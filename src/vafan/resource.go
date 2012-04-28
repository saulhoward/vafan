// Copyright 2012 Saul Howard. All rights reserved.

// Resource interface and helper functions.

package vafan

import (
	"code.google.com/p/gorilla/schema"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// Decodes form values.
var decoder = schema.NewDecoder()

// ResourceServer types serve Resource data over HTTP.
type ResourceServer interface {
	GetURL(req *http.Request, s *site) *url.URL
	//GetContent(req *http.Request, s *site) resourceContent
	ServeHTTP(w http.ResponseWriter, r *http.Request, u *user)
}

// Generic data map for Resource content.
type resourceContent map[string]interface{}

// Resource is the data structure commonly served by ResourceServer
// types, that can be served over HTTP.
type Resource struct {
	title       string
	description string
	content     resourceContent
}

// somewhat crufty helper map DELETE
//var emptyContent = map[string]interface{}{}

// crufty...
var resourceCanonicalSites = map[string]*site{
	"usersRegistrarResource": defaultSite,
	"usersAuthResource":      defaultSite,
	"usersSyncResource":      defaultSite,
}

// Makes a URL for a resource server.
// Used as a helper by struct.URL(r, s) methods.
func makeURL(res ResourceServer, req *http.Request, s *site, urlData []string) *url.URL {
	curSite, env := getSite(req)
	canonicalSite, err := getCanonicalSite(res)
	if s == nil {
		s = curSite
	}
	if err == nil && canonicalSite.Name != s.Name {
		s = canonicalSite
	}
	format := getFormat(req)
	format = "." + format
	if format == ".html" {
		format = ""
	}
	var host string
	if env == "production" {
		host = s.Host
	} else {
		host = env + "." + s.Host
	}
	/*
		if vafanConfig.port != "" {
			host = host + ":" + vafanConfig.port
		}
	*/
	urlPairs := []string{"format", format, "host", host}
	if urlData != nil {
		for _, p := range urlData {
			urlPairs = append(urlPairs, p)
		}
	}
	url, err := router.GetRoute(resourceName(res)).Host(hostRe).URL(urlPairs...)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to get URL for resource: %v", err))
		url, _ = url.Parse("/")
	}
	return url
}

// get the resource's type name by reflection
func resourceName(r ResourceServer) string {
	n := reflect.TypeOf(r).String()
	re := regexp.MustCompile(`\.([a-zA-Z]+)$`)
	m := re.FindStringSubmatch(n)
	return m[1]
}

// Helper function for resource ServeHTTP methods.
func (res Resource) write(w http.ResponseWriter, req *http.Request, serv ResourceServer, u *user) {
	// get the site and env requested
	s, env := getSite(req)
	format := getFormat(req)
	logger.Info(fmt.Sprintf("Requested url: '%v' writing resource '%v' as format '%v'.", req.URL.String(), resourceName(serv), format))
	// should we redirect to a canonical host for this resource?
	canonicalSite, err := getCanonicalSite(serv)
	if err == nil && canonicalSite.Name != s.Name {
		rUrl := serv.GetURL(req, nil)
		logger.Info(fmt.Sprintf("Redirecting to canonical url: " + rUrl.String()))
		http.Redirect(w, req, rUrl.String(), http.StatusMovedPermanently)
		return
	}

	// Add defaults to the content, that are in every format
	/* resContent := *new(resourceContent) */
	/* resContent = res.GetContent(req, s) */
	//logger.Info(fmt.Sprintf("%v", resContent.content["javascriptLibraryHTML"]))

	/* var content map[string]interface{} */
	/* content = resContent.content */

	res.content["links"] = getLinks(req)
	u.URL = u.GetURL(req, nil).String()
	res.content["requestingUser"] = u

	// write the resource in requested format
	if format == "html" {
		res.content["author"] = author
		res.content["title"] = res.title
		res.content["description"] = res.description
		res.content["site"] = s
		res.content["environment"] = env
		res.content["url"] = serv.GetURL(req, nil)
		res.content["resource"] = resourceDirName(serv)
		if flashes := getFlashContent(w, req); len(flashes) > 0 {
			res.content["flashes"] = flashes
		} else {
			res.content["flashes"] = ""
		}

		// add in CSS & JS (may be minified)
		res.content["javascriptLibraryHTML"] = getJavascriptLibraryHTML(s, env)
		res.content["cssHTML"] = getCSSHTML(s, env)

		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		t, err := getPageTemplate(format, serv, s, env)
		if err != nil {
			logger.Err(fmt.Sprintf("Failed to get template: %v", err))
			return
		}
		err = t.Execute(w, res.content)
		if err != nil {
			logger.Err(fmt.Sprintf("Failed executing template: %v", err))
		}
	} else if format == "json" {
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(res.content)
		if err != nil {
			err = logger.Err(fmt.Sprintf("Failed encoding JSON: %v", err))
		}
	} else {
		err = logger.Err(fmt.Sprintf("Format unknown: %v", format))
	}
	return
}
