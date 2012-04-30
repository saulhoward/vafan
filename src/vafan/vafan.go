// Vafan - a web server for a film studio
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"code.google.com/p/go.net/websocket"
	"code.google.com/p/gorilla/mux"
	"errors"
	"fmt"
	gzipFileServer "github.com/saulhoward/go-gzip-file-server"
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

// Config! Srsly...
var author string = "Saul Howard - saul@convictfilms.com"

// Represents a site served by this server
type site struct {
	Name              string
	Host              string
	Title             string
	FullTitle         string
	Tagline           string
	GoogleAnalyticsID string
}

// Should be in config
var sites = [...]site{
	site{"convict-films", "convictfilms.com", "Convict Films", "Convict Films", "We make movies", "UA-349594-6"},
	site{"brighton-wok", "brighton-wok.com", "Brighton Wok", "Brighton Wok: The Legend of Ganja Boxing", "Ninjas! Ganja! Kung Fu!", "UA-349594-1"},
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
	fullHost := vafanConf.host
	if vafanConf.port != "" {
		fullHost = fullHost + ":" + vafanConf.port
	}

	logger.Info(fmt.Sprintf("Starting the Vafan server, host '%v', port '%v.'", vafanConf.host, vafanConf.port))
	registerHandlers()
	http.Handle("/", router)
	http.ListenAndServe(fullHost, router)
}

// for when we haven't yet got a resource...
// restores host and scheme to the request's url
func getCurrentUrl(r *http.Request) *url.URL {
	url := r.URL
	if url.Host == "" {
		site, env := getSite(r)
		var host string
		if env == "production" {
			host = site.Host
		} else {
			host = env + "." + site.Host
		}
		/*
			if vafanConf.port != "" {
				host = host + ":" + vafanConf.port
			}
		*/
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
func callHandler(res ResourceServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := userCookie(w, r)

		logger.Info(fmt.Sprintf(
			"User: %v", u.ID))

		if err != nil {
			if err == ErrResourceRedirected {
				return
			} else {
				logger.Err(fmt.Sprintf(
					"Failed getting user: %v", err))
				u = NewUser()
			}
		}
		res.ServeHTTP(w, r, u)
		return
	}
}

// Set mux handlers
func registerHandlers() {
	logger.Info("Setting Handlers.")

	// Static directories
	router.PathPrefix("/css").Handler(
		http.StripPrefix("/css", gzipFileServer.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static", "css")))))
	router.PathPrefix("/js").Handler(
		http.StripPrefix("/js", gzipFileServer.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static", "js")))))
	router.PathPrefix("/img").Handler(
		http.StripPrefix("/img", http.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static", "img")))))
	router.PathPrefix("/fonts").Handler(
		http.StripPrefix("/fonts", http.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static", "fonts")))))

	// web standard static files
	router.Path("/favicon.ico").Handler(
		http.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static"))))
	router.Path(`/apple-touch-icon{appleIcon:[a-z0-9\-]*}.png`).Handler(
		http.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static"))))
	router.Path(`/robots.txt`).Handler(
		http.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static"))))
	router.Path(`/humans.txt`).Handler(
		http.FileServer(http.Dir(
			filepath.Join(vafanConf.baseDir, "static"))))

	// Regex strings used in url schemas
	formatRe := `{format:(\.{1}[a-z]+)?}`         // matches '.json' etc.
	uuidRe := `{id:[a-f0-9\-]+}`                  // matches UUIDs
	alphaNumericRe := `{%v:[\p{L}\p{M}\p{N}\-]+}` // matches unicode alphanumerics
	nameRe := fmt.Sprintf(alphaNumericRe, "name")

	// -- Resources

	router.Host(hostRe).Path(`/` + formatRe).
		Name("index").Handler(callHandler(index{}))

	router.Host(hostRe).Path(`/contact` + formatRe).
		Name("contact").Handler(callHandler(contact{}))

	// user resources
	router.Host(hostRe).Path(`/users/auth` + formatRe).
		Name("userAuth").Handler(callHandler(userAuth{}))
	router.Host(hostRe).Path(`/users/sync` + formatRe).
		Name("userSync").Handler(callHandler(userSync{}))
	router.Host(hostRe).Path(`/users/registrar` + formatRe).
		Name("userRegistrar").Handler(callHandler(userRegistrar{}))
	router.Host(hostRe).Path(`/users/` + uuidRe + formatRe).
		Name("user").Handler(callHandler(user{}))

	// Media resources (videos)
	router.Host(hostRe).Path(`/videos` + formatRe).
		Name("videos").Handler(callHandler(videos{}))
	router.Host(hostRe).Path(`/videos/` + nameRe + formatRe).
		Name("video").Handler(callHandler(video{}))

	// DVD resources
	router.Host(hostRe).Path(`/dvds/` + nameRe + formatRe).
		Name("dvd").Handler(callHandler(dvd{}))
	router.Host(hostRe).Path(`/dvds/` + nameRe + `/stockists/` + fmt.Sprintf(alphaNumericRe, "dvdStockist") + formatRe).
		Name("dvdStockist").Handler(callHandler(dvdStockist{}))
	router.Host(hostRe).Path(`/dvds/` + nameRe + `/stockists` + formatRe).
		Name("dvdStockists").Handler(callHandler(dvdStockists{}))

	// Twitter resource, inc. websockets resource.
	router.Host(hostRe).Path(`/tweets` + formatRe).
		Name("tweets").Handler(callHandler(tweets{}))
	router.Host(hostRe).Path(`/tweets/stream` + formatRe).
		Name("tweetStream").Handler(websocket.Handler(streamTweets))

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

func getCanonicalSite(r ResourceServer) (s *site, err error) {
	err = nil
	if resourceCanonicalSites[resourceName(r)] != nil {
		s = resourceCanonicalSites[resourceName(r)]
		return
	}
	err = errors.New("No canonical site set.")
	s = defaultSite
	return
}

// a map of site urls included in every response
// config?
func getLinks(req *http.Request) map[string]interface{} {
	l := make(map[string]interface{})

	bwLinks := make(map[string]string)
	cfLinks := make(map[string]string)
	siteLinks := make(map[string]string)

	cfLinks["index"] = index{}.GetURL(req, convictFilms).String()
	cfLinks["videos"] = videos{}.GetURL(req, convictFilms).String()
	cfLinks["contact"] = contact{}.GetURL(req, convictFilms).String()
	cfLinks["userAuth"] = userAuth{}.GetURL(req, nil).String()
	cfLinks["userRegistrar"] = userRegistrar{}.GetURL(req, nil).String()

	bwLinks["index"] = index{}.GetURL(req, brightonWok).String()
	bwLinks["videos"] = videos{}.GetURL(req, brightonWok).String()
	bwLinks["dvd"] = getBrightonWokDVD().GetURL(req, brightonWok).String()
	bwLinks["dvdStockists"] = getBrightonWokDVDStockists().GetURL(req, brightonWok).String()

	siteLinks["index"] = index{}.GetURL(req, nil).String()
	siteLinks["videos"] = videos{}.GetURL(req, nil).String()
	siteLinks["userAuth"] = userAuth{}.GetURL(req, nil).String()
	siteLinks["userRegistrar"] = userRegistrar{}.GetURL(req, nil).String()

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
	/*
		    vars := mux.Vars(r)
				if vars["host"] != "" {
					host = vars["host"]
				}
	*/
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
