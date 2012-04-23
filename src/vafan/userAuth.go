// Copyright 2012 Saul Howard. All rights reserved.

// User Auth. Logs users in and out.

package vafan

import (
	"fmt"
	"net/http"
	"net/url"
)

type userAuth struct {
	data resourceData
}

func (res userAuth) GetURL(req *http.Request, s *site) *url.URL {
	// limit authentication to default site
	return makeURL(res, req, defaultSite, nil)
}

func (res userAuth) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "Login"
	c.description = "Login here to access Convict Films"
	if res.data == nil {
		res.data = emptyContent
	}
	c.content = res.data
	return
}

func (res userAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	switch r.Method {
	case "POST":
		// This is a post to login or logout
		var url *url.URL
		r.ParseForm()
		switch {
		case r.Form.Get("login") != "":
			// try to login
			// TODO: THE CURRENT USER MUST BE LOGGED OUT

			// login user
			loginUser, err := login(r.Form.Get("UsernameOrEmailAddress"), r.Form.Get("Password"))
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to login user: %v", err))
				url = res.GetURL(r, nil)
				addFlash(w, r, "Failed to login", "error")
			} else {
				// set the login session
				_, err := newLoginSession(w, r, loginUser)
				if err == nil {
					addFlash(w, r, "Login!", "success")
				} else {
					logger.Err(fmt.Sprintf("Failed to set user session: %v", err))
				}
				url = index{}.GetURL(r, nil)
				addFlash(w, r, "Login!", "success")
			}
			http.Redirect(w, r, url.String(), http.StatusSeeOther)
		case r.Form.Get("logout") != "":
			// try to logout
			logout(w, r, reqU)
			addFlash(w, r, "Logged out.", "success")
			url = index{}.GetURL(r, nil)
			http.Redirect(w, r, url.String(), http.StatusSeeOther)
		}
	case "GET":
		writeResource(w, r, res, reqU)
	}
	return
}
