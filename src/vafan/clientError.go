// Copyright 2012 Saul Howard. All rights reserved.

// Provides content for 4XX Client Errors. Eg, 403, 404.

package vafan

import (
	"net/http"
	"net/url"
)

// 404 - Not Found.

type notFound struct {
}

func (nf notFound) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(nf, req, s, nil)
}

func (nf notFound) ServeHTTP(w http.ResponseWriter, r *http.Request, u *user) {
	res := Resource{
		title:       "404 - Not Found",
		description: "The requested resource does not exist.",
	}
	res.content = make(resourceContent)
	res.content["message"] = "404 Not found"
	res.content["body"] = "Sorry, this resource could not be found."

	w.WriteHeader(http.StatusNotFound)
	res.write(w, r, nf, u)
	return
}

// 403 - Forbidden 

type forbidden struct {
}

func (fb forbidden) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(fb, req, s, nil)
}

func (fb forbidden) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res := Resource{
		title:       "403 - Forbidden",
		description: "You are forbidden to access this resource",
	}
	res.content = make(resourceContent)
	res.content["message"] = "403 Forbidden"
	res.content["body"] = "Sorry, you may not access this resource."

	w.WriteHeader(http.StatusForbidden)
	res.write(w, r, fb, reqU)
	return
}
