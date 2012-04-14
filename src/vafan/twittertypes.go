/*
   Copyright 2012 Saul Howard. All rights reserved.

   Twitter base types.

   This file was taken from:

       https://github.com/araddon/httpstream/blob/master/twittertypes.go

   (With uneeded values removed). If extra values are needed, check
   there for help.
*/

package vafan

type twitterUser struct {
	Screen_name       string
	Url               string
	Profile_image_url string
	Id                int64
}

type tweet struct {
	Text       string
	RawBytes   []byte
	Id         int64
	Created_at string
	User       *twitterUser
}

type twitterSiteStreamMessage struct {
	For_user int64
	Message  tweet
}

type twitterEvent struct {
	Target     twitterUser
	Source     twitterUser
	Created_at string
	Event      string
}

type twitterFriendList struct {
	Friends []int64
}
