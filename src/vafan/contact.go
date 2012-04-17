// Copyright 2012 Saul Howard. All rights reserved.

// Contact resource. Returns contact information for a site.

package vafan

import (
	"net/http"
	"net/url"
)

type contact struct {
}

func (con contact) URL(req *http.Request, s *site) *url.URL {
	return getUrl(con, req, s, nil)
}

func (con contact) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = "Contact Us"
	c.description = "Contact details for " + s.Name
	c.content = emptyContent
	return
}

func (con contact) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	writeResource(w, r, con, reqU)
	return
}
