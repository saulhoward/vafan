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
	"code.google.com/p/gorilla/schema"
	"net/http"
)

// decodes form values
var decoder = schema.NewDecoder()

// data map returned by a resource
type resourceData map[string]interface{}

type Resource interface {
	name() string
	url() string
	content() resourceData
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// list of resource instances
var resources = map[string]Resource{
    "index": new(index),
    "usersRegistrar": new(usersRegistrar),
}

var emptyContent = map[string]interface{}{"title": "Go"}

// -- Index resource

type index struct {}

func (i *index) name() string {
    return "index"
}

func (i *index) url() string {
    return "/"
}

func (i *index) content() resourceData {
    return emptyContent
}

func (i *index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeResource(w, r, i)
}

// -- Registrar resource

type usersRegistrar struct {
    data resourceData
	canonicalSite string
}

func (u *usersRegistrar) name() string {
    return "usersRegistrar"
}

func (u *usersRegistrar) url() string {
    return "/users/registrar"
}

func (u *usersRegistrar) content() resourceData {
    return u.data
}

func (u *usersRegistrar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.data = map[string]interface{}{}
	switch r.Method {
	case "POST":
		// This is a post to create a new user
		r.ParseForm()
        user := NewUser()
		decoder.Decode(user, r.Form)

		// check for errors in post
        if !user.isLegal() || user.Password != r.Form.Get("RepeatPassword") {
            // found errors in post
			errors := map[string]interface{}{}
			if !user.isUsernameLegal() {
				errors["Username"] = "Must contain only letters and numbers, with no spaces."
            }
			if !user.isEmailAddressLegal() {
				errors["EmailAddress"] = "Must be a valid email address."
            }
			if !user.isPasswordLegal() {
				errors["Password"] = "Password must be more than 6 characters."
            } else if user.Password != r.Form.Get("RepeatPassword") {
				errors["Password"] = "Password must match repeat password."
			}
			u.data["user"] = user
			u.data["errors"] = errors
			writeResource(w, r, u)
			return
		}

        user.save()
        //rUrl := getUrl("usersAuth", r)
        //http.Redirect(w, r, rUrl.String(), http.StatusSeeOther)
		return
	case "GET":
        user := NewUser()
		u.data["user"] = user
		writeResource(w, r, u)
		return
	}
}

/*
// Auth resource
func userAuthHandler(w http.ResponseWriter, r *http.Request) {
	auth := new(resource)
	auth.name = "usersAuth"
	auth.content = map[string]interface{}{"title": "Go"}
	writeResource(w, r, auth)
}

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
