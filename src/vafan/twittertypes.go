/*
   Copyright 2012 Saul Howard. All rights reserved.

   Twitter base types.

   This file was taken from:

       https://github.com/araddon/httpstream/blob/master/twittertypes.go

   (With uneeded values removed and json names added). If extra values
   are needed, check there for help.
*/

package vafan

type twitterUser struct {
	Screen_name       string `json:"screenName"`
	Url               string `json:"url"`
	Profile_image_url string `json:"profileImageURL"`
	Id                int64  `json:"id"`
}

type tweet struct {
	Text       string       `json:"text"`
	Id         int64        `json:"id"`
	Created_at string       `json:"createdAt"`
	User       *twitterUser `json:"user"`
}

type twitterSiteStreamMessage struct {
	For_user int64 `json:"forUser"`
	Message  tweet `json:"message"`
}

type twitterEvent struct {
	Target     twitterUser `json:"target"`
	Source     twitterUser `json:"source"`
	Created_at string      `json:"createdAt"`
	Event      string      `json:"event"`
}

type twitterFriendList struct {
	Friends []int64 `json:"friends"`
}
