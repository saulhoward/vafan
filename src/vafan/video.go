// Copyright 2012 Saul Howard. All rights reserved.

// A video.

package vafan

import (
	"code.google.com/p/gorilla/mux"
	"errors"
	"fmt"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"net/url"
	"regexp"
)

// ErrVideoNotFound is returned by video when the named video does not
// exist in the database.
var ErrVideoNotFound = errors.New("video: doesn't exist")

// A video represents data describing a video, hosted on an external
// site such as Youtube or Vimeo
type video struct {
	Id          string
	Name        string // names are unique
	Title       string
	Description markdown
	Location    string
	Sites       []*site // the sites that display this vid
	Youtube     youtubeVideo
	Vimeo       vimeoVideo
}

// External video types, eg, youtube, vimeo
type externalVideo interface {
	FetchDetails() (err error)
}

// Video url uses the video name, eg, `/videos/brighton-wok`
func (v video) URL(req *http.Request, s *site) *url.URL {
	return getUrl(v, req, s, []string{"name", v.Name})
}

func (v video) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = "Video"
	c.description = "Video page"
	c.content = map[string]interface{}{"video": v}
	return
}

// GET sets the video from the URL vars.
func (v video) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		var err error
		vp, err := GetVideoByName(vars["name"])
		v = *vp
		if err != nil {
			if err == ErrVideoNotFound {
				notFound{}.ServeHTTP(w, r, reqU)
				return
			}
			_ = logger.Err(fmt.Sprintf("Failed to get video by name: %v", err))
			notFound{}.ServeHTTP(w, r, reqU)
			return
		}
		writeResource(w, r, &v, reqU)
		return
	}
}
func GetVideoByName(name string) (v *video, err error) {
	v = new(video)
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")
	err = c.Find(bson.M{"name": name}).One(v)
	if err != nil {
		if err == mgo.NotFound {
			err = ErrVideoNotFound
			return
		}
		_ = logger.Err(fmt.Sprintf("Failed to get video (Mongo): %v", err))
		return
	}
	return
}

func (v *video) save() (err error) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")
	err = c.Insert(v)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to insert video (Mongo): %v", err))
		return
	}
	return
}

// name must be unicode alphanumericals and dashes only
func (v *video) isNameLegal() bool {
	var illegalCharsRe = regexp.MustCompile(`[^\-\p{L}\p{M}\p{N}]+`)
	if illegalCharsRe.MatchString(v.Name) {
		return false
	}
	return true
}
