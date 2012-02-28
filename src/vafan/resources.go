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
)

type resource struct {
	name    string
	content map[string]interface{}
}

// Index resource
func indexHandler(w http.ResponseWriter, r *http.Request) {
    index := new(resource)
	index.name = "index"
	index.content = map[string]interface{}{"title": "Go"}
	writeResource(w, r, index)
}

// Video resource
func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    video := new(resource)
	video.name = "video"
	video.content = map[string]interface{}{"video": vars["video"]}
	writeResource(w, r, video)
}

// 404 resource
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Does not compute")
}

