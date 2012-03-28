// Copyright 2012 Saul Howard. All rights reserved.

// Index - the home page for a vafan site.

package vafan

import (
	"net/http"
	"net/url"
)

type index struct {
}

func (res index) URL(req *http.Request, s *site) *url.URL {
	return getUrl(res, req, s, nil)
}

func (res index) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = s.Tagline
	c.description = "Home page"
	c.content = emptyContent
	return
}

func (res index) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	writeResource(w, r, res, reqU)
	return
}
