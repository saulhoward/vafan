// Copyright 2012 Saul Howard. All rights reserved.

// User Sync. Provides syncing of user IDs across domains.
// Used by session.go

package vafan

import (
	"fmt"
	"net/http"
	"net/url"
)

type userSync struct {
}

func (res userSync) GetURL(req *http.Request, s *site) *url.URL {
	// limit sync to default site
	return makeURL(res, req, defaultSite, nil)
}

func (res userSync) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "User Sync"
	c.description = "Performs a user sync redirect"
	c.content = emptyContent
	return
}

// send people back to the redirect-url param, with a canonical user id
func (res userSync) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	ruStr := r.URL.Query().Get("redirect-url")
	if ruStr == "" {
		ruStr = "/"
	}
	ru, err := url.Parse(ruStr)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to parse redirect URL: %v", err))
		ru, _ = url.Parse("/")
	}
	q := ru.Query()
	ru.RawQuery = "" // remove the query string
	q.Set("canonical-user-id", url.QueryEscape(reqU.ID))
	fullURL := ru.String() + "?" + q.Encode() // and add it back
	logger.Info(fmt.Sprintf("Returning to URL: %v", fullURL))
	http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
	return
}
