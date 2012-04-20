// Copyright 2012 Saul Howard. All rights reserved.

// 4XX Client Error. Eg, 403, 404.

package vafan

import (
	"net/http"
	"net/url"
)

// 404 - Not Found.
type notFound struct {
}

func (res notFound) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(res, req, s, nil)
}

func (res notFound) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "404 - Not Found"
	c.description = "The requested resource does not exist."
	content := map[string]interface{}{}
	content["message"] = "404 Not found"
	content["body"] = "Sorry, this resource could not be found."
	c.content = content
	return
}

func (res notFound) ServeHTTP(w http.ResponseWriter, r *http.Request, u *user) {
	w.WriteHeader(http.StatusNotFound)
	writeResource(w, r, res, u)
	return
}

// 403 - Forbidden 
type forbidden struct {
}

func (res forbidden) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(res, req, s, nil)
}

func (res forbidden) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "403 - Forbidden"
	c.description = "You are forbidden to access this resource"
	content := map[string]interface{}{}
	content["message"] = "403 Forbidden"
	content["body"] = "Sorry, you may not access this resource."
	c.content = content
	return
}

func (res forbidden) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	w.WriteHeader(http.StatusForbidden)
	writeResource(w, r, res, reqU)
	return
}
