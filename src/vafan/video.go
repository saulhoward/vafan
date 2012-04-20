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
	"time"
)

// ErrVideoNotFound is returned by video when the named video does not
// exist in the database.
var ErrVideoNotFound = errors.New("video: doesn't exist")

// A video represents data describing a video, hosted on an external
// site such as Youtube or Vimeo
type video struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"` // names are unique
	Title            string         `json:"title"`
	Date             time.Time      `json:"date"`
	ShortDescription string         `json:"shortDescription"`
	Description      Markdown       `json:"description"`
	URL              string         `json:"url"`
	Thumbnail        Image          `json:"thumbnail"`
	Sites            []*site        `json:"sites"` // the sites that display this vid
	ExternalVideos   externalVideos `json:"externalVideos"`
}

// Image type.
type Image struct {
	URL    string `json:"id"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

// These satisfy the externalVideoProvider interface, below.
// I am explicit about their type so the mongo lib can write data
// straight into them.
type externalVideos struct {
	Youtube *youtubeVideo `json:"youtube"`
	Vimeo   *vimeoVideo   `json:"vimeo"`
}

// External video interface, eg, youtube, vimeo.
type externalVideoProvider interface {
	FetchData() (err error)
	getDefaultThumbnail() (i Image, err error)
}

// Video constructor.
func newVideo() (v *video) {
	v = new(video)
	return v
}

// Video url uses the video name, eg, `/videos/brighton-wok`
func (v video) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(v, req, s, []string{"name", v.Name})
}

func (v video) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "Video"
	c.description = "Video page"
	c.content = map[string]interface{}{"video": v}
	relatedVideos, err := getRelatedVideos(&v, s)
	if err == nil && len(relatedVideos) > 0 {
		for i, v := range relatedVideos {
			relatedVideos[i].URL = v.GetURL(req, nil).String()
		}
		c.content["relatedVideos"] = relatedVideos
	}
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

// Save a video to the DB.
func (v *video) save() (err error) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")

	//err = c.Insert(v)
	id, err := c.Upsert(bson.M{"name": v.Name}, v)

	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to insert or update video (Mongo): %v", err))
		return
	}
	_ = logger.Info(fmt.Sprintf("Inserted or updated video (Mongo): %v", id))

	return
}

// Fetches external data from the External Video Providers and saves it
// to the video.
func (v *video) UpdateExternalData() (err error) {
	_ = logger.Info("Updating video external data.")
	externalVideos := map[string]externalVideoProvider{}
	if v.ExternalVideos.Youtube.ID != "" {
		externalVideos["youtube"] = v.ExternalVideos.Youtube
	}
	if v.ExternalVideos.Vimeo.ID != "" {
		externalVideos["vimeo"] = v.ExternalVideos.Vimeo
	}
	var thumbs []Image
	for provider, extVid := range externalVideos {
		err = extVid.FetchData()
		if err != nil {
			_ = logger.Err(fmt.Sprintf("Failed to fetch external video details (%v): %v", provider, err))
			continue
		}
		err = v.save()
		if err != nil {
			_ = logger.Err(fmt.Sprintf("Failed to save video (%v): %v", provider, err))
			continue
		}

		// Set a default thumbnail for this video
		thumb, err := extVid.getDefaultThumbnail()
		if err == nil {
			thumbs = append(thumbs, thumb)
		}
		_ = logger.Info(fmt.Sprintf("Fetched video data for %v.", provider))
	}
	// Choose a thumbnail to be the default.
	v.Thumbnail = thumbs[0]
	v.save()
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

// -- Get Videos from DB.

// Get one video.
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

// Generic get videos - takes a selector.
func getVideos(selector bson.M) (vids []*video, err error) {
	vids = []*video{}
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")
	err = c.Find(selector).All(&vids)
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

// All videos for a particular site.
func getAllVideos(s *site) (v []*video, err error) {
	v = []*video{}
	v, err = getVideos(bson.M{"sites.name": s.Name})
	return
}

// Featured videos - used for the index videos.
func getFeaturedVideos(s *site) (v []*video, err error) {
	v = []*video{}
	v, err = getVideos(bson.M{"sites.name": s.Name})
	return
}

// Returns videos, except the video passed in, and will have tag
// matching.
func getRelatedVideos(v *video, s *site) (relVids []*video, err error) {
	relVids = []*video{}
	selector := bson.M{"sites.name": s.Name, "name": bson.M{"$ne": v.Name}}
	relVids, err = getVideos(selector)
	return
}
