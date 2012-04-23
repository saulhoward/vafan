// Copyright 2012 Saul Howard. All rights reserved.

// Resource interface and helper functions.

package vafan

import (
	"code.google.com/p/gorilla/schema"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// Decodes form values.
var decoder = schema.NewDecoder()

// A Resource represents data that can be served over HTTP.
type Resource interface {
	GetURL(req *http.Request, s *site) *url.URL
	GetContent(req *http.Request, s *site) resourceContent
	ServeHTTP(w http.ResponseWriter, r *http.Request, u *user)
}

// Generic data map for resource content.
type resourceData map[string]interface{}

type resourceContent struct {
	title       string
	description string
	content     resourceData
}

// somewhat crufty helper map
var emptyContent = map[string]interface{}{}

// crufty...
var resourceCanonicalSites = map[string]*site{
	"usersRegistrarResource": defaultSite,
	"usersAuthResource":      defaultSite,
	"usersSyncResource":      defaultSite,
}

// Makes a URL for a resource.
// Used as a helper by struct.URL(r, s) methods.
func makeURL(res Resource, req *http.Request, s *site, urlData []string) *url.URL {
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
func resourceName(r Resource) string {
	n := reflect.TypeOf(r).String()
	re := regexp.MustCompile(`\.([a-zA-Z]+)$`)
	m := re.FindStringSubmatch(n)
	return m[1]
}
