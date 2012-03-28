// Vafan - a web server for a film studio
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kless/goconfig/config"
	"log/syslog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

// Setup logging.
var logger = getLogger()

func getLogger() (l *syslog.Writer) {
	l, err := syslog.New(syslog.LOG_INFO, "vaf-serv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed initializing syslog: %v", err)
		panic(err)
	}
	return
}

// Get the config.
var conf = getConfig()

func getConfig() (c *config.Config) {
	//var conf, _ = config.ReadDefault("/srv/vafan/config/config.ini")
	c, err := config.ReadDefault("/home/saul/vafan-config.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading configuration: %v", err)
		panic(err)
	}
	return
}

var author string = "Saul Howard - saul@convictfilms.com"

// Represents a site served by this server
type site struct {
	Name      string
	Host      string
	Title     string
	FullTitle string
	Tagline   string
}

// Should be in config
var sites = [...]site{
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
var router = new(mux.Router) // gorilla's router

// Start the server up
func StartServer() {
	_ = logger.Info("Starting the Vafan server, port 8888.")
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
				err = logger.Err(fmt.Sprintf("Failed getting user: %v", err))
				u = NewUser()
			}
		}
		res.ServeHTTP(w, r, u)
		return
	}
}

// Set mux handlers
func registerHandlers() {
	_ = logger.Info("Setting Handlers.")

	var baseDir, err = conf.String("default", "base-dir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading 'base-dir' from configuration: %v", err)
		os.Exit(1)
	}

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
		Name("index").Handler(callHandler(index{}))

	// user resources
	router.Host(hostRe).Path(`/users/auth` + formatRe).
		Name("userAuth").Handler(callHandler(userAuth{}))
	router.Host(hostRe).Path(`/users/sync` + formatRe).
		Name("userSync").Handler(callHandler(userSync{}))
	router.Host(hostRe).Path(`/users/registrar` + formatRe).
		Name("userRegistrar").Handler(callHandler(userRegistrar{}))
	router.Host(hostRe).Path(`/users/` + uuidRe + formatRe).
		Name("user").Handler(callHandler(user{}))

	// media resources
	router.Host(hostRe).Path(`/videos` + formatRe).
		Name("videos").Handler(callHandler(videos{}))
	router.Host(hostRe).Path(`/videos/` + nameRe + formatRe).
		Name("video").Handler(callHandler(video{}))

	/* router.Host(hostRe).Path(`/movies/` + nameRe + formatRe). */
	/* Name("moviesResource").Handler(callHandler(moviesResource{})) */

	// http status codes
	router.Host(hostRe).Path(`/403` + formatRe).
		Name("forbidden").Handler(callHandler(forbidden{}))
	router.Host(hostRe).Path(`/404` + formatRe).
		Name("notFound").Handler(callHandler(notFound{}))

	// 404
	router.NotFoundHandler = callHandler(notFound{})
}

func getCanonicalSite(r Resource) (s *site, err error) {
	err = nil
	if resourceCanonicalSites[resourceName(r)] != nil {
		s = resourceCanonicalSites[resourceName(r)]
		return
	}
	err = errors.New("No canonical site set.")
	s = defaultSite
	return
}

func writeResource(w http.ResponseWriter, req *http.Request, res Resource, u *user) {
	_ = logger.Info(fmt.Sprintf("Requested url: '%v' writing resource '%v'", req.URL.String(), resourceName(res)))
	// get the site and env requested
	s, env := getSite(req)
	format := getFormat(req)
	// should we redirect to a canonical host for this resource?
	canonicalSite, err := getCanonicalSite(res)
	if err == nil && canonicalSite.Name != s.Name {
		rUrl := res.URL(req, nil)
		_ = logger.Info(fmt.Sprintf("Redirecting to canonical url: " + rUrl.String()))
		http.Redirect(w, req, rUrl.String(), http.StatusMovedPermanently)
		return
	}

	// Add defaults to the content, that are in every format
	resContent := res.Content(req, s)
	content := resContent.content
	content["links"] = getLinks(req)
	u.Location = u.URL(req, nil).String()
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
		t, err := getPageTemplate(format, res, s)
		if err != nil {
			_ = logger.Err(fmt.Sprintf("Failed to get template: %v", err))
			return
		}
		err = t.Execute(w, content)
		if err != nil {
			_ = logger.Err(fmt.Sprintf("Failed executing template: %v", err))
		}
	} else if format == "json" {
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(content)
		if err != nil {
			err = logger.Err(fmt.Sprintf("Failed encoding JSON: %v", err))
		}
	} else {
		err = logger.Err(fmt.Sprintf("Format unknown: %v", format))
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

	cfLinks["index"] = index{}.URL(req, convictFilms).String()
	cfLinks["videos"] = videos{}.URL(req, convictFilms).String()
	cfLinks["userAuth"] = userAuth{}.URL(req, nil).String()
	cfLinks["userRegistrar"] = userRegistrar{}.URL(req, nil).String()

	bwLinks["index"] = index{}.URL(req, brightonWok).String()
	bwLinks["videos"] = videos{}.URL(req, brightonWok).String()

	siteLinks["index"] = index{}.URL(req, nil).String()
	siteLinks["videos"] = videos{}.URL(req, nil).String()
	siteLinks["userAuth"] = userAuth{}.URL(req, nil).String()
	siteLinks["userRegistrar"] = userRegistrar{}.URL(req, nil).String()

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
	return "html"
}
