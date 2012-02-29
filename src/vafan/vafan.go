// Vafan - a web server for Convict Films
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"os"
	"regexp"
	"log"
	"path/filepath"
	"net/http"
	"encoding/json"
	"code.google.com/p/gorilla/mux"
	//"github.com/kless/goconfig/config"
)

// get the config
//var conf, _ = config.ReadDefault("/home/saul/code/vafan/config/config.ini")
//var baseDir, _ = conf.String("default", "base-dir")
var baseDir string = "/srv/vafan"

// set up the router
var router = new(mux.Router)

// Start the server up
func StartServer() {
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
    hostRe := `{host:[a-z]*}`
	formatRe := `{format:(\.{1}[a-z]+)?}`

	router.Host(hostRe).Path("/").HandlerFunc(indexHandler)
	router.Host(hostRe).Path("/home" + formatRe).HandlerFunc(indexHandler)
	router.Host(hostRe).Path("/index" + formatRe).HandlerFunc(indexHandler)

	router.Host(hostRe).Path("/videos/{video}" + formatRe).HandlerFunc(videoHandler)
}

func writeResource(w http.ResponseWriter, req *http.Request, res *resource) {
	site, env := getSite(req)
	format := getFormat(req)
	if format == "html" {
        res.content["environment"] = env;
		w.Header().Add("Content-Type", "text/html")
		t := getPageTemplate(format, res, site)
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

var sites = map[string]string {
    "brighton-wok": "brighton-wok.com",
    "convict-films": "convictfilms.com",
}

var envs = [...]string{"dev", "testing", "production"}

func getSite(r *http.Request) (site string, env string) {
    // default values
    env = "production"
    site = "convict-films"
    // get the host (from mux or in the Host: field)
    host := r.Host
	vars := mux.Vars(r)
    if vars["host"] != "" {
        host = vars["host"]
    }
    /* Should use one regex, perhaps like...
        var envRe string
        for i, env := range envs {
            if i != 0 {
                envRe += "|"
            }
            envRe += env
        }
    */
    for possSite, possHost := range sites {
        var hostRe = regexp.MustCompile(possHost)
        if hostRe.MatchString(host) {
            site = possSite
            break
        }
    }
    for _, possEnv := range envs {
        var envRe = regexp.MustCompile("^" + possEnv + ".")
        if envRe.MatchString(host) {
            env = possEnv
            break
        }
    }
    return
}

func getFormat(r *http.Request) string {
    // should also allow Content-Accept
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

//-- crappy helper things

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

