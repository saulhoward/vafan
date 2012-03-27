// Vafan - a web server for a film studio
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"os"
	"errors"
	"regexp"
	"log"
	"path/filepath"
	"net/url"
	"net/http"
	"encoding/json"
	"code.google.com/p/gorilla/mux"
	"github.com/kless/goconfig/config"
)

// get the config
//var conf, _ = config.ReadDefault("/srv/vafan/config/config.ini")
var conf, _ = config.ReadDefault("/home/saul/vafan-config.ini")

//var baseDir, _ = conf.String("default", "base-dir")
var baseDir string = "/srv/vafan"

var author string = "Saul Howard - saul@convictfilms.com"

// Represents a site served by this server
type site struct {
    Name string
    Host string
    Title string
    FullTitle string
    Tagline string
}
// Should be in config
var sites = [...]site {
    site{"convict-films", "convictfilms.com", "Convict Films", "Convict Films", "We make movies"},
    site{"brighton-wok", "brighton-wok.com", "Brighton Wok", "Brighton Wok: The Legend of Ganja Boxing", "Ninjas! Ganja! Kung Fu!"},
}
var defaultSite *site = &sites[0]
var convictFilms *site = &sites[0]
var brightonWok *site = &sites[1]

// Can stay as sensible default?
var envs = [...]string{"dev", "testing", "production"}
var defaultEnv string = "production"

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
       host := env + "." + site.Host + ":8888"
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
func callHandler(res Resource) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        u, err := userCookie(w, r)
        if err != nil {
            if err == ErrResourceRedirected {
                return
            } else {
                checkError(err)
            }
        }
        res.ServeHTTP(w, r, u)
        return
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

    // web standard static files
	router.Path("/favicon.ico").Handler(
        http.FileServer(http.Dir(
            filepath.Join(baseDir, "static"))))

    // Regex strings used in url schemas
	formatRe := `{format:(\.{1}[a-z]+)?}`   // matches '.json' etc.
    uuidRe := `{id:[a-f0-9\-]+}`            // matches UUIDs
    nameRe := `{name:[\p{L}\p{M}\p{N}\-]+}` // matches unicode alphanumerics

    // -- Resources

    router.Host(hostRe).Path(`/` + formatRe).
        Name("indexResource").Handler(callHandler(indexResource{}))

    // user resources
    router.Host(hostRe).Path(`/users/auth` + formatRe).
        Name("usersAuthResource").Handler(callHandler(usersAuthResource{}))
    router.Host(hostRe).Path(`/users/sync` + formatRe).
        Name("usersSyncResource").Handler(callHandler(usersSyncResource{}))
    router.Host(hostRe).Path(`/users/registrar` + formatRe).
        Name("usersRegistrarResource").Handler(callHandler(usersRegistrarResource{}))
    router.Host(hostRe).Path(`/users/` + uuidRe + formatRe).
        Name("usersResource").Handler(callHandler(usersResource{}))

    // media resources
    router.Host(hostRe).Path(`/videos` + formatRe).
        Name("videosResource").Handler(callHandler(videosResource{}))
    router.Host(hostRe).Path(`/videos/` + nameRe + formatRe).
        Name("videoResource").Handler(callHandler(videoResource{}))

    /* router.Host(hostRe).Path(`/movies/` + nameRe + formatRe). */
        /* Name("moviesResource").Handler(callHandler(moviesResource{})) */

    // http status codes
    router.Host(hostRe).Path(`/403` + formatRe).
        Name("forbiddenResource").Handler(callHandler(forbiddenResource{}))
    router.Host(hostRe).Path(`/404` + formatRe).
        Name("notFoundResource").Handler(callHandler(notFoundResource{}))

    // 404
    router.NotFoundHandler = callHandler(notFoundResource{})
}

func getCanonicalSite(r Resource) (s *site, err error) {
    err = nil
    if resourceCanonicalSites[resourceName(r)] != nil  {
        s = resourceCanonicalSites[resourceName(r)]
        return
    }
    err = errors.New("No canonical site set.")
    s = defaultSite
    return
}

func writeResource(w http.ResponseWriter, req *http.Request, res Resource, u *User) {
    print("\nWriting resource " + resourceName(res) + " for request " + req.URL.String())
    // get the site and env requested
	s, env := getSite(req)
	format := getFormat(req)
    // should we redirect to a canonical host for this resource?
    canonicalSite, err := getCanonicalSite(res)
    if err == nil && canonicalSite.Name != s.Name {
        rUrl := res.URL(req, nil)
        print("\nRedirecting to canonical url... " + rUrl.String())
        http.Redirect(w, req, rUrl.String(), http.StatusMovedPermanently)
        return
    }

    // Add defaults to the content, that are in every format
    resContent := res.Content(req, s)
    content := resContent.content
    content["links"] = getLinks(req)
    u.URL = u.getURL(req)
    content["requestingUser"] = u

    // write the resource in requested format
    if format == "html" {
        content["author"] = author
        content["title"] = resContent.title
        content["description"] = resContent.description
        content["site"] = s
        content["environment"] = env
        content["url"] = res.URL(req, nil)
        content["resource"] = resourceDirName(res)
        if flashes := getFlashContent(w, req); len(flashes) > 0 {
            content["flashes"] = flashes
        } else {
            content["flashes"] = ""
        }

		w.Header().Add("Content-Type", "text/html")
		t := getPageTemplate(format, res, s)
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
// config?
func getLinks(req *http.Request) map[string]interface{} {
    l := make(map[string]interface{})

    bwLinks := make(map[string]string)
    cfLinks := make(map[string]string)
    siteLinks := make(map[string]string)

    cfLinks["index"] = indexResource{}.URL(req, convictFilms).String()
    cfLinks["videos"] = videosResource{}.URL(req, convictFilms).String()
    cfLinks["usersAuth"] = usersAuthResource{}.URL(req, nil).String()
    cfLinks["usersRegistrar"] = usersRegistrarResource{}.URL(req, nil).String()

    bwLinks["index"] = indexResource{}.URL(req, brightonWok).String()
    bwLinks["videos"] = videosResource{}.URL(req, brightonWok).String()

    siteLinks["index"] = indexResource{}.URL(req, nil).String()
    siteLinks["videos"] = videosResource{}.URL(req, nil).String()
    siteLinks["usersAuth"] = usersAuthResource{}.URL(req, nil).String()
    siteLinks["usersRegistrar"] = usersRegistrarResource{}.URL(req, nil).String()

    l["site"] = siteLinks
    l["brightonWok"] = bwLinks
    l["convictFilms"] = cfLinks

    return l
}

func getSite(r *http.Request) (s *site, env string) {
    // default values
    env = defaultEnv
    s = defaultSite
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
    for _, possSite := range sites {
        var hostRe = regexp.MustCompile(possSite.Host)
        if hostRe.MatchString(host) {
            s = &possSite
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

