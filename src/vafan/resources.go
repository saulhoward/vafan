// Vafan - a web server for Convict Films
//
// Resource handlers
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	//"fmt"
	//"log"
	//"code.google.com/p/gorilla/mux"
	"net/url"
	"net/http"
	"code.google.com/p/gorilla/schema"
    "code.google.com/p/gorilla/sessions"
)

var sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

func MyHandler(w http.ResponseWriter, r *http.Request) {

}

// decodes form values
var decoder = schema.NewDecoder()

// data map returned by a resource
type resourceData map[string]interface{}

type Resource interface {
	name() string
	urlSchema() string
	content() resourceData
    ServeHTTP(w http.ResponseWriter, r *http.Request, u *User)
}

// list of resource instances
var resources = map[string]Resource{
    "index": new(index),
    "usersRegistrar": new(usersRegistrar),
    "usersAuth": new(usersAuth),
    "usersSync": new(usersSync),
}

var resourceCanonicalSites = map[string]string{
    "usersRegistrar": "convict-films",
    "usersAuth": "convict-films",
    "usersSync": "convict-films",
}

// unnecessary helper cruft
var emptyContent = map[string]interface{}{}

// -- Index resource

type index struct {}

func (res *index) name() string {
    return "index"
}

func (res *index) urlSchema() string {
    return "/"
}

func (res *index) content() resourceData {
    return emptyContent
}

func (res *index) ServeHTTP(w http.ResponseWriter, r *http.Request, u *User) {
	writeResource(w, r, res, u)
    return
}

// -- Registrar resource

type usersRegistrar struct {
    data resourceData
}

func (res *usersRegistrar) name() string {
    return "usersRegistrar"
}

func (res *usersRegistrar) urlSchema() string {
    return "/users/registrar"
}

func (res *usersRegistrar) content() resourceData {
    return res.data
}

func (res *usersRegistrar) ServeHTTP(w http.ResponseWriter, r *http.Request, u *User) {
	res.data = map[string]interface{}{}
	switch r.Method {
	case "POST":
		// This is a post to create a new user
		r.ParseForm()
		decoder.Decode(u, r.Form)

		// check for errors in post
        if !u.isLegal() || u.Password != r.Form.Get("RepeatPassword") || !u.isNew() {
            // found errors in post
			errors := map[string]interface{}{}
			if !u.isUsernameLegal() {
				errors["Username"] = "Must contain only letters and numbers, with no spaces."
            } else if !u.isUsernameNew() {
				errors["Username"] = "Username already taken, sorry."
            }
			if !u.isEmailAddressLegal() {
				errors["EmailAddress"] = "Must be a valid email address."
            } else if !u.isEmailAddressNew() {
				errors["EmailAddress"] = "This email address is already associated with another user."
            }
			if !u.isPasswordLegal() {
				errors["Password"] = "Password must be more than 6 characters."
            } else if u.Password != r.Form.Get("RepeatPassword") {
				errors["Password"] = "Password must match repeat password."
			}
			res.data["errors"] = errors
			writeResource(w, r, res, u)
			return
		}

        // legal user, try to save
        err := u.save()
        var url *url.URL
        if err != nil {
            url = getUrl(resources["usersRegistrar"], r)
            // add flash to session
            session, _ := sessionStore.Get(r, "vafanFlashes")
            session.AddFlash("Failed to save new user", "error")
            session.Save(r, w)
        } else {
            url = getUrl(resources["usersAuth"], r)
            // add flash to session
            session, _ := sessionStore.Get(r, "vafanFlashes")
            session.AddFlash("Registered a new user, please log in.", "success")
            session.Save(r, w)
        }

        http.Redirect(w, r, url.String(), http.StatusSeeOther)
		return
	case "GET":
        if u.isNew() {
            writeResource(w, r, res, u)
        } else {
            url := getUrl(resources["usersAuth"], r)
            session, _ := sessionStore.Get(r, "vafanFlashes")
            session.AddFlash("Your user ID already has an account, please log in.", "warning")
            session.Save(r, w)
            http.Redirect(w, r, url.String(), http.StatusSeeOther)
        }
		return
	}
}

// -- Auth resource

type usersAuth struct {
    data resourceData
}

func (res *usersAuth) name() string {
    return "usersAuth"
}

func (res *usersAuth) urlSchema() string {
    return "/users/auth"
}

func (res *usersAuth) content() resourceData {
    return res.data
}

func (res *usersAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, u *User) {
	res.data = emptyContent

	switch r.Method {
	case "POST":
		// This is a post to login or logout
		r.ParseForm()
		decoder.Decode(u, r.Form)
        if r.Form.Get("login") != "" {
            // login user
            sessionId, err := u.Login()
            if err != nil {
                // set the session
                // Get a session. We're ignoring the error resulted from decoding an
                // existing session: Get() always returns a session, even if empty.
                session, _ := store.Get(r, "vafanLogin")
                // Set some session values.
                session.Values["id"] = sessionId
                // Save it.
                session.Save(r, w)
            }
        }

	case "GET":
        writeResource(w, r, res, u)
    }
    return
}

// -- Sync resource

type usersSync struct {
}

func (res *usersSync) name() string {
    return "usersSync"
}

func (res *usersSync) urlSchema() string {
    return "/users/sync"
}

func (res *usersSync) content() resourceData {
    return emptyContent
}

// send people back to the redirect-url param, with a canonical user id
func (res *usersSync) ServeHTTP(w http.ResponseWriter, r *http.Request, u *User) {
    ruStr := r.URL.Query().Get("redirect-url")
    if ruStr == "" {
        ruStr = "/"
    }
    ru, err := url.Parse(ruStr)
    checkError(err)
    q := ru.Query()
    ru.RawQuery = "" // remove the query string
    q.Set("canonical-user-id", url.QueryEscape(u.Id))
    fullUrl := ru.String() + "?" + q.Encode() // and add it back
    print("\nReturning to url... " + fullUrl)
    http.Redirect(w, r, fullUrl, http.StatusTemporaryRedirect)
    return
}

/*
// Video resource
func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	video := new(resource)
	video.name = "videos"
	video.content = map[string]interface{}{"video": vars["video"]}
	writeResource(w, r, video)
}

// 404 resource
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Does not compute")
}
*/
