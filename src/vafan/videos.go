// Copyright 2012 Saul Howard. All rights reserved.

// Collection of videos. 

package vafan

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Represents a collection of videos.
type videos struct {
	video  *video   // a new video, as it is being added
	videos []*video // the collection of videos
}

func (vids videos) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(vids, req, s, nil)
}

// View the videos, and, if user has permission, get a form to post a
// new video.
func (vids videos) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res := Resource{
		title:       "Video Library",
		description: "Video collection",
	}
	res.content = make(resourceContent)

	switch r.Method {
	case "GET":
		var err error
		site, _ := getSite(r)
		vids.videos, err = getAllVideos(site)
		if err == nil {
			for i, v := range vids.videos {
				vids.videos[i].URL = v.GetURL(r, nil).String()
			}
			res.content["videos"] = vids.videos
		}
		res.write(w, r, vids, reqU)
		return

	case "POST":
		if reqU.Role == "superadmin" {
			// This is a post to create a new video
			vids.video = newVideo()
			r.ParseForm()
			decoder.Decode(vids.video, r.Form)

			// as markdown is not a string, ParseForm() won't decode it
			vids.video.Description = Markdown(r.Form.Get("Description"))

			// Date - two possible formats.
			var dateErr error
			vids.video.Date, dateErr = time.Parse("2006-01-02", r.Form.Get("Date"))
			if dateErr != nil {
				vids.video.Date, dateErr = time.Parse("2006-01-02 15:04:05 +0000 UTC", r.Form.Get("Date"))
			}

			// All other video data
			// TODO check youtube, vimeo ids
			if vids.video.Sites == nil || vids.video.isNameLegal() == false ||
				vids.video.ShortDescription == "" || dateErr != nil ||
				vids.video.Title == "" || vids.video.Description == "" {
				// found errors in post
				errors := map[string]interface{}{}
				if vids.video.isNameLegal() == false {
					errors["Name"] = "Must contain only alphanumericals and dashes.."
				}
				if dateErr != nil {
					errors["Date"] = "Date is unreadable. It must look like 2012-04-01."
				}
				if vids.video.Title == "" {
					errors["Title"] = "Must have title."
				}
				if vids.video.ShortDescription == "" {
					errors["ShortDescription"] = "Must have short description."
				}
				if vids.video.Description == "" {
					errors["Description"] = "Must have description."
				}
				if vids.video.Sites == nil {
					errors["Sites"] = "Must have at least one site."
				}
				res.content["errors"] = errors

				res.content["video"] = vids.video
				res.write(w, r, vids, reqU)
				return
			}

			// legal viddya, try to save
			err := vids.video.save()

			// Fetch extra metadata from external sources (concurrently)
			go vids.video.UpdateExternalData()

			var url *url.URL
			if err != nil {
				logger.Err(fmt.Sprintf("Failed to save new video: %v", err))
				url = videos{}.GetURL(r, nil)
				addFlash(w, r, "Failed to save new video", "error")
			} else {
				logger.Info(fmt.Sprintf("Saved new video: %v", vids.video.ID))
				url = vids.video.GetURL(r, nil)
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
