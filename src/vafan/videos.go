// Copyright 2012 Saul Howard. All rights reserved.

// Collection of videos. 

package vafan

import (
	"fmt"
	"net/http"
	"net/url"
)

// Represents a collection of videos.
type videos struct {
	video *video
	data  resourceData
}

func (res videos) URL(req *http.Request, s *site) *url.URL {
	return getUrl(res, req, s, nil)
}

func (res videos) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = "Video"
	c.description = "Video page"
	if res.data != nil {
		c.content = res.data
	} else {
		c.content = emptyContent
	}
	if res.video != nil {
		res.data["video"] = res.video
	}
	return
}

func (res videos) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res.data = emptyContent
	switch r.Method {
	case "GET":

		writeResource(w, r, res, reqU)
		return
	case "POST":
		if reqU.Role == "superadmin" {
			// This is a post to create a new user
			r.ParseForm()
			res.video = new(video)
			decoder.Decode(res.video, r.Form)

			// check for errors in post
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
