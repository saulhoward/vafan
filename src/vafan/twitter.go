// Copyright 2012 Saul Howard. All rights reserved.

// Twitter.

package vafan

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/araddon/httpstream"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	twConvictFilmsID = 25679402
	twSaulHowardID   = 7273252
)

type twitter struct {
	tweets []*httpstream.Tweet
	data   resourceData // assembled data for response
}

func (tw tweets) URL(req *http.Request, s *site) *url.URL {
	return getUrl(tw, req, s, nil)
}

func (tw tweets) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = "Tweets"
	c.description = "Tweets"

	tw.tweets, err = getTweets()
	if err != nil {
		tw.tweets = nil
	}
	tw.data["tweets"] = tw.tweets

	if tw.data != nil {
		c.content = tw.data
	} else {
		c.content = emptyContent
	}
	return
}

func (tw tweets) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	writeResource(w, r, tw, reqU)
	return
}

func getTweets() (tw []*httpstream.Tweet) {

	// get latest tweets from redis

	// return them (may be empty)

	//if they're empty, or if they're old, get them from the twitter
	// api concurrently and put them in redis for next time.

}

// Websocket handler.
func streamTweets(ws *websocket.Conn) {
	writeTweetStream(ws)
}

// Connect to twitter streaming api and send lines to be written.
func writeTweetStream(w io.Writer) {
	var twUser, twPwd string
	var err error
	twUser, err = conf.String("twitter", "user")
	twPwd, err = conf.String("twitter", "password")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed reading twitter user/password from configuration: %v", err))
		return
	}

	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), "info")
	stream := make(chan []byte, 1000)
	done := make(chan bool)
	client := httpstream.NewBasicAuthClient(twUser, twPwd, func(line []byte) {
		stream <- line
	})
	//err = client.Filter([]int64{twSaulHowardID, twConvictFilmsID}, []string{"brightonwok", "brighton wok", "convictfilms"}, false, done)
	err = client.Filter([]int64{twSaulHowardID, twConvictFilmsID}, []string{"tsunami", "brighton wok", "convictfilms"}, false, done)
	// this opens a go routine that is effectively thread 1
	// err := client.Sample(done)
	if err != nil {
		println(err.Error())
	}
	// 2nd thread
	go func() {
		for {
			line := <-stream
			writeTweet(w, line)
		}
	}()
	// 3rd thread
	for {
		line := <-stream
		writeTweet(w, line)
	}
}

// Write one line from the streaming api.
func writeTweet(w io.Writer, line []byte) {
	switch {
	case bytes.HasPrefix(line, []byte(`{"event":`)):
		var event httpstream.Event
		json.Unmarshal(line, &event)
	case bytes.HasPrefix(line, []byte(`{"friends":`)):
		var friends httpstream.FriendList
		json.Unmarshal(line, &friends)
	default:
		tweet := httpstream.Tweet{}
		json.Unmarshal(line, &tweet)
		if tweet.User != nil {
			//println(th, " ", tweet.User.Screen_name, ": ", tweet.Text)
			w.Write([]byte(tweet.Text))
		}
	}
}
