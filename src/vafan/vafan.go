// Vafan - a web server for Convict Films
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package main

import (
	"os"
	"fmt"
	"log"
	"path/filepath"
	"net/http"
	"html/template"
	"encoding/json"
	"code.google.com/p/gorilla/mux"
	"github.com/kless/goconfig/config"
)

// get the config
var conf, _ = config.ReadDefault("/home/saul/code/vafan/config/config.ini")
var baseDir, _ = conf.String("default", "base-dir")

// set up the router
var router = new(mux.Router)

// Start the server up
func main() {
	setHandlers()
	http.Handle("/", router)
	http.ListenAndServe(":8888", router)
}

// Set mux handlers
func setHandlers() {
	// Static directories
	router.PathPrefix("/css").Handler(
		http.StripPrefix("/css", http.FileServer(http.Dir(
			filepath.Join(baseDir, "static", "css")))))
	router.PathPrefix("/js").Handler(
		http.StripPrefix("/js", http.FileServer(http.Dir(
			filepath.Join(baseDir, "static", "js")))))
	router.PathPrefix("/img").Handler(
		http.StripPrefix("/img", http.FileServer(http.Dir(
			filepath.Join(baseDir, "static", "img")))))

	// Dynamic funcs
	formatRe := `{format:(\.{1}[a-z]+)?}`

	router.Path("/").HandlerFunc(indexHandler)
	router.Path("/home" + formatRe).HandlerFunc(indexHandler)
	router.Path("/index" + formatRe).HandlerFunc(indexHandler)

	router.Path("/videos/{video}" + formatRe).HandlerFunc(videoHandler)
}

// Index resource
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var index resource
	index.name = "index"
	index.content = map[string]interface{}{"title": "Go"}
	writeResource(w, r, index)
}

// Video resource
func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var video resource
	video.name = "video"
	video.content = map[string]interface{}{"video": vars["video"]}
	writeResource(w, r, video)
}

// 404 resource
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Does not compute")
}

//-- resource

type resource struct {
	name    string
	content map[string]interface{}
}

func getSite(r *http.Request) string {
	return "brighton-wok"
}

func getFormat(r *http.Request) string {
	vars := mux.Vars(r)
	if vars["format"] == "" || vars["format"] == ".html" {
		return "html"
	} else if vars["format"] == ".json" {
		return "json"
	}
    // srsly?
    os.Exit(1)
    return "error"
}

func writeResource(w http.ResponseWriter, req *http.Request, res resource) {
	site := getSite(req)
	format := getFormat(req)
	if format == "html" {
		w.Header().Add("Content-Type", "text/html")
		t := getTemplate(format, res, site)
		err := t.Execute(w, res.content)
		checkError(err)
	} else if format == "json" {
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(res.content)
		checkError(err)
	} else {
		// error checking here pls
		os.Exit(1)
	}
}

//-- crappy helper things

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//-- Template functions

func getTemplate(format string, res resource, site string) *template.Template {
	//Check for the most specific template first
    for i:= 0; templateExists(format, res, site) == false; i++ {
		if i == 0 {
			site = "_anySite"
		} else if i == 1 {
			res.name = "_anyResource"
		} else if i == 2 {
			format = "_anyFormat"
		} else if i > 2 {
			// error checking here pls
			os.Exit(1)
		}
	}
	path := filepath.Join(baseDir, "templates", format, res.name, site, "main.html")
	t, err := template.New("main.html").ParseFiles(path)
	checkError(err)
	return t
}

func templateExists(format string, res resource, site string) bool {
	path := filepath.Join(baseDir, "templates", format, res.name, site, "main.html")
	_, err := os.Stat(path)
	if err != nil {
		print(err.Error() + "\n")
		return false
	}
	return true
}
