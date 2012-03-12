// Vafan - a web server for a film studio
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

// Should be in config
var sites = map[string]string {
    "brighton-wok": "brighton-wok.com",
    "convict-films": "convictfilms.com",
}
// Can stay as sensible default?
var envs = [...]string{"dev", "testing", "production"}

var hostRe = `{host:[a-z0-9\.\:]*}`

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
    print("\nSetting Handlers")
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

    // Home resource
	router.Host(hostRe).Path("/").
        Name("index").
        HandlerFunc(indexHandler)
	router.Host(hostRe).
        Path("/home" + formatRe).
        HandlerFunc(indexHandler)
	router.Host(hostRe).
        Path("/index" + formatRe).
        HandlerFunc(indexHandler)

    // User resources
	router.Host(hostRe).
        Path("/users/auth" + formatRe).
        Name("usersAuth").
        HandlerFunc(userAuthHandler)
	router.Host(hostRe).
        Path("/users/registrar" + formatRe).
        Name("usersRegistrar").
        HandlerFunc(userRegistrarHandler)

    // Video resources
	router.Host(hostRe).
        Path("/videos/{video}" + formatRe).
        Name("videos").
        HandlerFunc(videoHandler)
}

func writeResource(w http.ResponseWriter, req *http.Request, res *resource) {
	site, env := getSite(req)
	format := getFormat(req)
    // should we redirect to a canonical host for this resource?
    if res.canonicalSite != "" && res.canonicalSite != site {
        rHost := env + "." + sites[res.canonicalSite] + ":8888"
        rFormat := "." + format
        if rFormat == ".html" {
            rFormat = ""
        }
        rUrl, err := router.GetRoute(res.name).Host(hostRe).URL("format", rFormat, "host", rHost)
        checkError(err)
        print("\nRedirecting to canonical url... " + rUrl.String())
        w.Header().Set("Location", rUrl.String())
        http.Redirect(w, req, rUrl.String(), http.StatusMovedPermanently)
        return
    }
    // write the resource in requested format
    if format == "html" {
        res.content["environment"] = env;
        res.content["url"] = res.url;
        res.content["resource"] = res.name;
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
    return
}

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

