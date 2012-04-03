// Copyright 2012 Saul Howard. All rights reserved.

// Collection of videos. 

package vafan

import (
	"fmt"
	"launchpad.net/mgo"
	"net/http"
	"net/url"
)

// Represents a collection of videos.
type videos struct {
	video  *video       // a new video, being added
	videos []video      // the collection of videos
	data   resourceData // other data
}

func (res videos) URL(req *http.Request, s *site) *url.URL {
	return getUrl(res, req, s, nil)
}

func (res videos) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = "Video Library"
	c.description = "Video collection"
	if res.data != nil {
		c.content = res.data
	} else {
		c.content = emptyContent
	}
	if res.video != nil {
		res.data["video"] = res.video
	}
	if res.videos != nil {
		for i, v := range res.videos {
			res.videos[i].Location = v.URL(req, nil).String()
		}
		res.data["videos"] = res.videos
	}
	return
}

// View the videos, and, if user has permission, get a form to post a new video
func (res videos) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res.data = emptyContent
	var err error
	res.videos, err = getAllVideos()
	if err != nil {
		res.videos = nil
	}
	switch r.Method {
	case "GET":
		writeResource(w, r, res, reqU)
		return
	case "POST":
		if reqU.Role == "superadmin" {
			// This is a post to create a new video
			r.ParseForm()
			res.video = new(video)
			decoder.Decode(res.video, r.Form)
			// as markdown is not a string, ParseForm() won't decode it
			res.video.Description = Markdown(r.Form.Get("Description"))

			if res.video.Sites == nil || res.video.isNameLegal() == false ||
				res.video.Title == "" || res.video.Description == "" {
				// found errors in post
				errors := map[string]interface{}{}
				if res.video.isNameLegal() == false {
					errors["Name"] = "Must contain only alphanumericals and dashes.."
				}
				if res.video.Title == "" {
					errors["Title"] = "Must have title."
				}
				if res.video.Description == "" {
					errors["Description"] = "Must have description."
				}
				if res.video.Sites == nil {
					errors["Sites"] = "Must have at least one site."
				}
				res.data["errors"] = errors
				writeResource(w, r, res, reqU)
				return
			}

			// legal viddya, try to save
			err := res.video.save()

			// Fetch extra metadata from external sources (concurrently)
			go res.video.UpdateExternalData()

			var url *url.URL
			if err != nil {
				_ = logger.Err(fmt.Sprintf("Failed to save new video: %v", err))
				url = videos{}.URL(r, nil)
				addFlash(w, r, "Failed to save new video", "error")
			} else {
				_ = logger.Info(fmt.Sprintf("Saved new video: %v", res.video.Id))
				url = res.video.URL(r, nil)
				addFlash(w, r, "Added a video!", "success")
			}

			http.Redirect(w, r, url.String(), http.StatusSeeOther)
			return

		} else {
			forbidden{}.ServeHTTP(w, r, reqU)
			return
		}
	}
	return
}

func getAllVideos() (v []video, err error) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")
	err = c.Find(nil).All(&v)
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
