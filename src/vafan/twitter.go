// Copyright 2012 Saul Howard. All rights reserved.

// Twitter.

package vafan

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/araddon/httpstream"
	"github.com/fzzbt/radix"
	oauth "github.com/reinaldons/goauth"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var twAuth *oauth.OAuth

const (
	twConvictFilmsID = 25679402
	twSaulHowardID   = 7273252
)

type tweets struct {
	tweets []*httpstream.Tweet
	data   resourceData // assembled data for response
}

func (tw tweets) URL(req *http.Request, s *site) *url.URL {
	return getUrl(tw, req, s, nil)
}

func (tw tweets) Content(req *http.Request, s *site) (c resourceContent) {
	c.title = "Tweets"
	c.description = "Tweets"
	var err error
	tw.tweets, err = getLatestTweets()
	if err != nil {
		tw.tweets = nil
	}
	if tw.tweets != nil {
		tw.data["tweets"] = tw.tweets
	}

	if tw.data != nil {
		c.content = tw.data
	} else {
		c.content = emptyContent
	}

	return
}

func (tw tweets) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	tw.data = emptyContent
	writeResource(w, r, tw, reqU)
	return
}

// Return a 'clever' selection of the latest tweets.
func getLatestTweets() (tws []*httpstream.Tweet, err error) {

	// get latest tweets from redis

	// return them (may be empty)

	//if they're empty, or if they're old, get them from the twitter
	// api concurrently and put them in redis for next time.

	// --

	tws, _ = getTweets("user_timeline")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed fetching user timeline tweets: %v", err))
	}

	// Only fetch new lot of tweets if necessary.
	go storeUserTimelineTweets()

	//tws = append(tws, userTimeline)
	return
}

// Fetch and cache recent tweets posted by the user.
func storeUserTimelineTweets() (err error) {
	_ = logger.Info("storing...")
	params := make(map[string]string)
	params["user_id"] = string(twConvictFilmsID)
	params["count"] = "20"
	params["include_rts"] = "1"
	tws, err := fetchTweets("https://api.twitter.com/1/statuses/user_timeline.json", params)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed fetching user timeline tweets: %v", err))
		//return
	}
	err = saveTweets("user_timeline", tws)
	return
}

// Save an array of tweets in Redis.
func saveTweets(key string, tws []*httpstream.Tweet) (err error) {
	_ = logger.Info("saving...")
	// Saving tweets as a json blob.
	twsJSON, err := json.Marshal(tws)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed marshalling tweets json: %v", err))
		return
	}
	db := radix.NewClient(redisConfiguration)
	defer db.Close()
	reply := db.Command("set", "tweets:"+key, twsJSON)
	if reply.Error() != nil {
		errText := fmt.Sprintf("Failed to set tweet data (Redis): %v", reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	return
}

// Get tweets from Redis.
func getTweets(key string) (tws []*httpstream.Tweet, err error) {
	db := radix.NewClient(redisConfiguration)
	defer db.Close()
	reply := db.Command("get", "tweets:"+key)
	if reply.Error() != nil {
		errText := fmt.Sprintf("Failed to get tweet data (Redis): %v", reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	twsJSON := reply.String()
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed getting tweets json (Redis): %v", err))
		return
	}
	tws = []*httpstream.Tweet{}
	err = json.Unmarshal([]byte(twsJSON), &tws)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed unmarshalling tweets json: %v", err))
		return
	}
	return
}

// Generic tweet fetcher. Requires REST URL.
func fetchTweets(url string, params map[string]string) (tws []*httpstream.Tweet, err error) {
	// Get OAuth, set as a Global
	if twAuth == nil {
		twAuth, err = getTwitterAuth()
		if err != nil {
			_ = logger.Err(fmt.Sprintf("Failed initialising twitter authentication: %v", err))
			return
		}
	}
	response, err := twAuth.Get(url, params)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed accessing twitter REST API: %v", err))
		return
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed reading twitter response body: %v", err))
		return
	}
	tws = []*httpstream.Tweet{}
	err = json.Unmarshal([]byte(data), &tws)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Error unmarshalling tweet JSON: %v", err))
		err = nil
	}
	return
}

// Creates oauth using details from config.
func getTwitterAuth() (o *oauth.OAuth, err error) {
	o = new(oauth.OAuth)
	var (
		consumerKey    string
		consumerSecret string
		accessToken    string
		accessSecret   string
	)
	consumerKey, err = conf.String("twitter", "consumer-key")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading twitter details from configuration: %v", err)
	}
	consumerSecret, err = conf.String("twitter", "consumer-secret")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading twitter details from configuration: %v", err)
	}
	accessToken, err = conf.String("twitter", "access-token")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading twitter details from configuration: %v", err)
	}
	accessSecret, err = conf.String("twitter", "access-secret")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading twitter details from configuration: %v", err)
	}
	o.ConsumerKey = consumerKey
	o.ConsumerSecret = consumerSecret
	o.AccessToken = accessToken
	o.AccessSecret = accessSecret
	o.SignatureMethod = "HMAC-SHA1"
	return
}

// -- Streaming tweets via websockets --

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
	err = client.Filter([]int64{twSaulHowardID, twConvictFilmsID}, []string{"brightonwok", "brighton wok", "convictfilms", "convict films"}, false, done)
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
