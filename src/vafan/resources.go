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
)

// decodes form values
var decoder = schema.NewDecoder()

// data map returned by a resource
type resourceData map[string]interface{}

type Resource interface {
	name() string
	urlSchema() string
	title(s *site) string
	description() string
	content() resourceData
    ServeHTTP(w http.ResponseWriter, r *http.Request, u *User)
}

// list of resource instances
var resources = map[string]Resource{
    "index":          new(index),
    "usersRegistrar": new(usersRegistrar),
    "usersAuth":      new(usersAuth),
    "usersSync":      new(usersSync),
    "notFound":       new(notFound),
}

// from config, eventually
var resourceCanonicalSites = map[string]*site{
    "usersRegistrar": defaultSite,
    "usersAuth":      defaultSite,
    "usersSync":      defaultSite,
}

// unnecessary helper cruft
var emptyContent = map[string]interface{}{}

// -- Index resource

type index struct {}

func (res *index) name() string {
    return "index"
}
func (res *index) title(s *site) string {
    return s.Tagline
}
func (res *index) description() string {
    return "Home page"
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

func (res *usersRegistrar) title(s *site) string {
    return "Register"
}

func (res *usersRegistrar) description() string {
    return "Register here to access Convict Films"
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
        if !u.isLegal(r.Form.Get("Password")) || r.Form.Get("Password") != r.Form.Get("RepeatPassword") || !u.isNew() {
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
			if !u.isPasswordLegal(r.Form.Get("Password")) {
				errors["Password"] = "Password must be more than 6 characters."
            } else if r.Form.Get("Password") != r.Form.Get("RepeatPassword") {
				errors["Password"] = "Password must match repeat password."
			}
			res.data["errors"] = errors
			writeResource(w, r, res, u)
			return
		}

        // legal user, try to save
        err := u.save(r.Form.Get("Password"))
        var url *url.URL
        if err != nil {
            url = getUrl(resources["usersRegistrar"], r)
            addFlash(w, r, "Failed to save new user", "error")
        } else {
            url = getUrl(resources["usersAuth"], r)
            addFlash(w, r, "Registered a new user, please log in.", "success")
        }

        http.Redirect(w, r, url.String(), http.StatusSeeOther)
		return
	case "GET":
        if u.isNew() {
            writeResource(w, r, res, u)
        } else {
            url := getUrl(resources["usersAuth"], r)
            addFlash(w, r, "Your user ID already has an account, please log in.", "warning")
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

func (res *usersAuth) title(s *site) string {
    return "Login"
}

func (res *usersAuth) description() string {
    return "Login here to access Convict Films"
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
        // try to login
        if r.Form.Get("login") != "" {
            // TODO: THE CURRENT USER MUST BE LOGGED OUT

            var url *url.URL
            // login user
            loginUser, err := login(r.Form.Get("UsernameOrEmailAddress"), r.Form.Get("Password"))
            if err != nil {
                url = getUrl(resources["usersAuth"], r)
                addFlash(w, r, "Failed to login", "error")
            } else {
                // set the login session
                _, err := newLoginSession(w, r, loginUser)
                if err != nil {
                    checkError(err)
                }
                url = getUrl(resources["index"], r)
                addFlash(w, r, "Login!", "success")
            }
            http.Redirect(w, r, url.String(), http.StatusSeeOther)
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

func (res *usersSync) title(s *site) string {
    return "User Sync"
}

func (res *usersSync) description() string {
    return "Performs a user sync redirect"
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
*/

// -- 404 resource
type notFound struct {
}

func (res *notFound) name() string {
    return "notFound"
}

func (res *notFound) title(s *site) string {
    return "404 - Not Found"
}

func (res *notFound) description() string {
    return "The requested resource does not exist."
}

func (res *notFound) urlSchema() string {
    return "/404"
}

func (res *notFound) content() resourceData {
    var content = map[string]interface{}{}
    content["message"] = "404 Not Found"
    content["body"] = "Sorry, this resource could not be found"
    return content
}

func (res *notFound) ServeHTTP(w http.ResponseWriter, r *http.Request, u *User) {
    w.WriteHeader(http.StatusNotFound)
    writeResource(w, r, res, u)
    return
}

