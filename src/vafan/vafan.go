// Vafan - a web server for a film studio
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"os"
	//"fmt"
	"regexp"
	"log"
	"path/filepath"
	"net/url"
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

// this is silly - we should just match any host
var hostRe = `{host:[a-z0-9\.\:\-]*}`

// set up the router
var router = new(mux.Router)  // gorilla's router

// Start the server up
func StartServer() {
	registerHandlers()
	http.Handle("/", router)
	http.ListenAndServe(":8888", router)
}

// for when we haven't yet got a resource...
// restores host and scheme to the request's url
func getCurrentUrl(r *http.Request) *url.URL {
   url := r.URL
   if url.Host == "" {
       site, env := getSite(r)
       host := env + "." + sites[site] + ":8888"
       url.Host = host
   }
   if url.Scheme == "" {
       if r.TLS != nil {
           url.Scheme = "https"
       } else {
           url.Scheme = "http"
       }
   }
   return url
}

// Handler wrapper - wraps resource requests
func makeHandler(res Resource) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // is there a login cookie?

        // if not, use the normal user cookie
        u := userCookie(w, r)
        res.ServeHTTP(w, r, u)
    }
}

// Set mux handlers
func registerHandlers() {
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
	router.PathPrefix("/fonts").Handler(
		http.StripPrefix("/fonts", http.FileServer(http.Dir(
			filepath.Join(baseDir, "static", "fonts")))))

	// Dynamic handlers - the resources
	formatRe := `{format:(\.{1}[a-z]+)?}`
    for _, r := range resources {
        router.Host(hostRe).
            Path(r.urlSchema() + formatRe).
            Name(r.name()).
            Handler(makeHandler(r))
    }

    /*

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
        HandlerFunc(usersRegistrarHandler)

    // Video resources
	router.Host(hostRe).
        Path("/videos/{video}" + formatRe).
        Name("videos").
        HandlerFunc(videoHandler)
        */
}

func getCanonicalSite(r Resource) (site string) {
    if resourceCanonicalSites[r.name()] != ""  {
        site = resourceCanonicalSites[r.name()]
    }
    return
}


func getUrlForSite(res Resource, site string, req *http.Request) *url.URL {
    curSite, env := getSite(req)
    canonicalSite := getCanonicalSite(res)
    if canonicalSite != "" && canonicalSite != site {
        site = canonicalSite
    }
    if site == "" {
        site = curSite
    }

	format := getFormat(req)
    format = "." + format
    if format == ".html" {
        format = ""
    }
    host := env + "." + sites[site] + ":8888"
    url, err := router.GetRoute(res.name()).Host(hostRe).URL("format", format, "host", host)
    checkError(err)
    return url
}

func getUrl(res Resource, req *http.Request) *url.URL {
    return getUrlForSite(res, "", req)
}

func writeResource(w http.ResponseWriter, req *http.Request, res Resource, u *User) {
    // get the site and env requested
	site, env := getSite(req)
	format := getFormat(req)
    // should we redirect to a canonical host for this resource?
    canonicalSite := getCanonicalSite(res)
    if canonicalSite != "" && canonicalSite != site {
       rUrl := getUrl(res, req)
        print("\nRedirecting to canonical url... " + rUrl.String())
        http.Redirect(w, req, rUrl.String(), http.StatusMovedPermanently)
        return
    }

    // Add defaults to the content, that are in every format
    content := res.content()
    links := make(map[string]interface{})
    links["site"] = getSiteLinks(req)
    content["links"] = links
    content["user"] = u

    // write the resource in requested format
    if format == "html" {
        content["site"] = site
        content["environment"] = env
        content["url"] = getUrl(res, req)
        content["resource"] = res.name()
        if flashes := getFlashContent(w, req); len(flashes) > 0 {
            content["flashes"] = flashes
        } else {
            content["flashes"] = ""
        }

		w.Header().Add("Content-Type", "text/html")
		t := getPageTemplate(format, res, site)
        err := t.Execute(w, content)
		checkError(err)
	} else if format == "json" {
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(content)
		checkError(err)
	} else {
		// error checking here pls
		os.Exit(1)
	}
    return
}

// a map of site urls included in every response
func getSiteLinks(req *http.Request) map[string]string {
    l := make(map[string]string)
    l["convictFilms"] = getUrlForSite(resources["index"], "convict-films", req).String()
    l["brightonWok"] = getUrlForSite(resources["index"], "brighton-wok", req).String()
    l["index"] = getUrl(resources["index"], req).String()
    l["usersAuth"] = getUrl(resources["usersAuth"], req).String()
    l["usersRegistrar"] = getUrl(resources["usersRegistrar"], req).String()
    return l
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

