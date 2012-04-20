// Copyright 2012 Saul Howard. All rights reserved.

// Index - the home page for a vafan site.
// Will contain a bit of everything, some videos, some photos etc.
// As the main landing page, content should be A/B tested, rotated etc.

package vafan

import (
	"net/http"
	"net/url"
)

type index struct {
	videos []*video        // featured videos
	dvds   map[string]*dvd // featured dvds
	tweets tweets          // recent tweets
	data   resourceData    // assembled data for response
}

func (res index) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(res, req, s, nil)
}

func (res index) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = s.Tagline
	c.description = "Home page"

	var err error

	res.videos, err = getFeaturedVideos(s)
	if err == nil {
		for i, v := range res.videos {
			res.videos[i].URL = v.GetURL(req, nil).String()
		}
		res.data["videos"] = res.videos
	}

	res.dvds, err = getFeaturedDVDs(s)
	if err == nil {
		for i, d := range res.dvds {
			res.dvds[i].URL = d.GetURL(req, nil).String()
		}
		res.data["dvds"] = res.dvds
	}

	res.tweets, err = getFeaturedTweets()
	if err == nil {
		res.data["tweets"] = res.tweets
	}

	if res.data != nil {
		c.content = res.data
	} else {
		c.content = emptyContent
	}
	return
}

func (res index) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res.data = emptyContent
	writeResource(w, r, res, reqU)
	return
}
