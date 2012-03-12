// Vafan - a web server for Convict Films
//
// Resource handlers
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"fmt"
	//"log"
	"net/http"
	"code.google.com/p/gorilla/mux"
    "code.google.com/p/gorilla/schema"
)

// decodes form values
var decoder = schema.NewDecoder()

type resource struct {
	name    string
	url    string
	canonicalSite    string
	content map[string]interface{}
}

type user struct {
    Username string
    EmailAddress string
    Password string
}

// Index resource
func indexHandler(w http.ResponseWriter, r *http.Request) {
    index := new(resource)
	index.name = "index"
	index.url = "/"
	index.content = map[string]interface{}{"title": "Go"}
	writeResource(w, r, index)
}

// Registrar resource
func userRegistrarHandler(w http.ResponseWriter, r *http.Request) {
    res := new(resource)
    res.name = "usersRegistrar"
    res.url = "/users/registrar"
    res.canonicalSite = "convict-films"
    switch r.Method {
    case "POST":
        // post to create a new user
        r.ParseForm()
        user := new(user)
        decoder.Decode(user, r.Form)
        print(user.Username)
        print(r.Form.Get("repeatPassword"))
        return
    case "GET":
        res.content = map[string]interface{}{"title": "Go"}
        writeResource(w, r, res)
        return
    }
}

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

