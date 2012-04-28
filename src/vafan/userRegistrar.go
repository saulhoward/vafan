// Copyright 2012 Saul Howard. All rights reserved.

// User Registrar. Registers new users.

package vafan

import (
	"fmt"
	"net/http"
	"net/url"
)

type userRegistrar struct {
}

func (reg userRegistrar) GetURL(req *http.Request, s *site) *url.URL {
	// limit registration to default site
	return makeURL(reg, req, defaultSite, nil)
}

func (reg userRegistrar) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res := Resource{
		title:       "Register",
		description: "Register here to access Convict Films",
	}
	res.content = make(resourceContent)

	switch r.Method {
	case "POST":
		// This is a post to register the requesting user
		r.ParseForm()
		u := &user{ID: reqU.ID}
		decoder.Decode(u, r.Form)

		// check for errors in post
		if !u.isLegal(r.Form.Get("Password")) || r.Form.Get("Password") != r.Form.Get("RepeatPassword") || !u.isNew() {
			// found errors in post
			errors := map[string]interface{}{}
			if !u.isUsernameLegal() {
				errors["Username"] = "Must contain only letters and numbers, with no spaces."
			} else if !u.isUsernameNew() {
				errors["Username"] = "Username already taken, sorry."
			}
			if !u.isEmailAddressLegal() {
				errors["EmailAddress"] = "Must be a valid email address."
			} else if !u.isEmailAddressNew() {
				errors["EmailAddress"] = "This email address is already associated with another user."
			}
			if !u.isPasswordLegal(r.Form.Get("Password")) {
				errors["Password"] = "Password must be more than 6 characters."
			} else if r.Form.Get("Password") != r.Form.Get("RepeatPassword") {
				errors["Password"] = "Password must match repeat password."
			}
			res.content["errors"] = errors
			res.write(w, r, reg, u)
			return
		}

		// legal user, try to save
		err := u.save(r.Form.Get("Password"))
		var url *url.URL
		if err != nil {
			logger.Err(fmt.Sprintf("Failed to save new user: %v", err))
			url = reg.GetURL(r, nil)
			addFlash(w, r, "Failed to save new user", "error")
		} else {
			url = userAuth{}.GetURL(r, nil)
			addFlash(w, r, "Registered a new user, please log in.", "success")
		}
		http.Redirect(w, r, url.String(), http.StatusSeeOther)
		return

	case "GET":
		if reqU.isNew() {
			res.write(w, r, reg, reqU)
		} else {
			url := userAuth{}.GetURL(r, nil)
			addFlash(w, r, "Your user ID already has an account, please log in.", "warning")
			http.Redirect(w, r, url.String(), http.StatusSeeOther)
		}
		return
	}
}
