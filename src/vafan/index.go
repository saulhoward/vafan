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
	videos []*video     // featured videos
	tweets tweets       // recent tweets
	data   resourceData // assembled data for response
}

func (res index) URL(req *http.Request, s *site) *url.URL {
	return getUrl(res, req, s, nil)
}

func (res index) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = s.Tagline
	c.description = "Home page"

	var err error

	res.videos, err = getFeaturedVideos(s)
	if err == nil {
		for i, v := range res.videos {
			res.videos[i].Location = v.URL(req, nil).String()
		}
		res.data["videos"] = res.videos
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
