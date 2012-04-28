// Copyright 2012 Saul Howard. All rights reserved.

// Contact resource. Returns contact information for a site.

package vafan

import (
	"net/http"
	"net/url"
)

type contact struct {
}

func (con contact) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(con, req, s, nil)
}

func (con contact) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	s, _ := getSite(r)
	res := Resource{
		title:       "Contact Us",
		description: "Contact details for " + s.Name,
	}
	res.content = make(resourceContent)
	res.write(w, r, con, reqU)
	return
}
