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
}

func (ind index) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(ind, req, s, nil)
}

func (ind index) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {

	logger.Info("index")

	s, _ := getSite(r)
	res := Resource{
		title:       s.Tagline,
		description: "Home page",
	}
	res.content = make(resourceContent)

	var err error

	ind.videos, err = getFeaturedVideos(s)
	if err == nil {
		for i, v := range ind.videos {
			ind.videos[i].URL = v.GetURL(r, nil).String()
		}
		res.content["videos"] = ind.videos
	}

	ind.dvds, err = getFeaturedDVDs(s)
	if err == nil {
		for i, d := range ind.dvds {
			ind.dvds[i].URL = d.GetURL(r, nil).String()
		}
		res.content["dvds"] = ind.dvds
	}

	ind.tweets, err = getFeaturedTweets()
	if err == nil {
		res.content["tweets"] = ind.tweets
	}

	res.write(w, r, ind, reqU)
	return
}
