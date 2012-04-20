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
	"sort"
	"strconv"
	"time"
)

// A global auth type, used by oauth package.
var twAuth *oauth.OAuth

// When tweets are fetched from cache, their timestamp is checked and
// may return these errors.
var ErrStaleTweets = errors.New("Twitter: Tweets are more than 10 min old.")
var ErrNoTweetsFound = errors.New("Twitter: No tweets found.")

// Twitter user IDs - should be in config!
const (
	twConvictFilmsID = 25679402
	twSaulHowardID   = 7273252
)

// The list of tweets. A tweet is defined in twittertypes.go
type tweets []*tweet

// Methods implementing Sort interface

func (tws tweets) Len() int {
	return len(tws)
}

func (tws tweets) Swap(i, j int) {
	tws[i], tws[j] = tws[j], tws[i]
}

// Provides a specific timestamp sort.
type reverseCreatedAtTweets struct {
	tweets
}

// Sort on 'created_at' timestamp, latest first.
func (tws reverseCreatedAtTweets) Less(i, j int) bool {
	iTime, _ := time.Parse("Mon Jan 02 15:04:05 +0000 2006", tws.tweets[i].Created_at)
	jTime, _ := time.Parse("Mon Jan 02 15:04:05 +0000 2006", tws.tweets[j].Created_at)
	return iTime.After(jTime)
}

// Methods implementing Resource interface.

func (tw tweets) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(tw, req, s, nil)
}

func (tw tweets) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "Crew Tweets"
	c.description = "Tweets about Convict Films"
	c.content = emptyContent
	if tw != nil {
		c.content["tweets"] = tw
	}
	return
}

func (tw tweets) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	tw, err := getFeaturedTweets()
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Error when getting tweets: %v", err))
	}
	writeResource(w, r, tw, reqU)
	return
}

// Functions to fetch tweet collections from REST API, save them to
// Redis cache, and return them.

func getFeaturedTweets() (tws tweets, err error) {
	tws, err = getTweets("featured")
	if err != nil {
		if err == ErrStaleTweets || err == ErrNoTweetsFound {
			go storeFeaturedTweets()
			err = nil
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting featured tweets (from cache): %v", err))
		}
	}
	return
}

// Featured tweets are compiled from the other types of tweets.
func storeFeaturedTweets() (err error) {
	timeline, err := getTweets("user_timeline")
	if err != nil {
		if err == ErrStaleTweets || err == ErrNoTweetsFound {
			go storeUserTimelineTweets()
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting user timeline tweets (from cache): %v", err))
		}
	}
	mentions, err := getTweets("mentions")
	if err != nil {
		if err == ErrStaleTweets || err == ErrNoTweetsFound {
			go storeMentionTweets()
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting mention tweets (from cache): %v", err))
		}
	}
	favorites, err := getTweets("favorites")
	if err != nil {
		if err == ErrStaleTweets || err == ErrNoTweetsFound {
			go storeFavoriteTweets()
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting favorite tweets (from cache): %v", err))
		}
	}

	featured := tweets{}
	if len(favorites) > 0 {
		featured = append(featured, favorites...)
	}
	if len(mentions) > 0 {
		featured = append(featured, mentions...)
	}
	if len(timeline) > 0 {
		featured = append(featured, timeline...)
	}

	if len(featured) > 0 {
		sort.Sort(reverseCreatedAtTweets{featured})
		err = saveTweets("featured", featured[:8])
	}
	return
}

// Fetch and cache recent tweets posted by the user.
func storeUserTimelineTweets() (err error) {
	params := make(map[string]string)
	params["user_id"] = strconv.Itoa(twConvictFilmsID)
	params["count"] = "20"
	params["include_rts"] = "1"
	tws, err := fetchTweets("https://api.twitter.com/1/statuses/user_timeline.json", params)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed fetching user timeline tweets (ntp?): %v", err))
		return
	}
	err = saveTweets("user_timeline", tws)
	return
}

// Fetch and cache recent tweets which mention the user.
func storeMentionTweets() (err error) {
	params := make(map[string]string)
	params["user_id"] = strconv.Itoa(twConvictFilmsID)
	params["count"] = "20"
	params["include_rts"] = "1"
	tws, err := fetchTweets("https://api.twitter.com/1/statuses/mentions.json", params)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed fetching mention tweets (ntp?): %v", err))
		return
	}
	err = saveTweets("mentions", tws)
	return
}

// Fetch and cache recent tweets which the user has 'favored'.
func storeFavoriteTweets() (err error) {
	params := make(map[string]string)
	params["user_id"] = strconv.Itoa(twConvictFilmsID)
	params["count"] = "20"
	tws, err := fetchTweets("https://api.twitter.com/1/favorites.json", params)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed fetching favorite tweets (ntp?): %v", err))
		return
	}
	err = saveTweets("favorites", tws)
	return
}

// Save an array of tweets in Redis.
func saveTweets(key string, tws tweets) (err error) {
	_ = logger.Info(fmt.Sprintf("Saving '%v' tweets to Redis", key))
	// Saving tweets as a json blob.
	twsJSON, err := json.Marshal(tws)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed marshalling '%v' tweets json: %v", key, err))
		return
	}
	db := radix.NewClient(redisConfiguration)
	defer db.Close()

	timestamp := time.Now().Unix()
	twsMap := map[string]string{
		"tweets":    string(twsJSON),
		"timestamp": strconv.Itoa(int(timestamp)),
	}
	tweetKey := "tweets:" + key
	reply := db.Command("hmset", tweetKey, twsMap)
	if reply.Error() != nil {
		errText := fmt.Sprintf("Failed to set '%v' tweet data (Redis): %v", key, reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	return
}

// Get tweets from Redis cache. Returned error warns if empty or stale.
func getTweets(key string) (tws tweets, err error) {
	db := radix.NewClient(redisConfiguration)
	defer db.Close()
	reply := db.Command("hgetall", "tweets:"+key)
	if reply.Error() != nil {
		errText := fmt.Sprintf("Failed to get '%v' tweet data (Redis): %v", key, reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	twsMap, err := reply.StringMap()
	if err != nil {
		errText := fmt.Sprintf("Stringmap for '%v' tweets failed (Redis): %v", key, reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}

	tws = tweets{}
	err = json.Unmarshal([]byte(twsMap["tweets"]), &tws)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed unmarshalling '%v' tweets json: %v", key, err))
		err = ErrNoTweetsFound
		return
	}

	// Check the timestamp for freshness.
	then, err := strconv.Atoi(twsMap["timestamp"])
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed getting timestamp for '%v' tweets: %v", key, err))
		return
	}
	now := int(time.Now().Unix())
	diff := now - then
	_ = logger.Info(fmt.Sprintf("'%v' tweets are %v seconds old.", key, diff))
	if diff > 600 {
		_ = logger.Info(fmt.Sprintf("'%v' tweets are more than 10 min old (%v secs).", key, diff))
		err = ErrStaleTweets
	}

	// Did we get any tweets?
	if len(tws) < 1 {
		err = ErrNoTweetsFound
	}

	return
}

// Generic tweet fetcher. Requires REST URL.
func fetchTweets(url string, params map[string]string) (tws tweets, err error) {
	_ = logger.Info(fmt.Sprintf("Fetching tweets from: %v", url))
	// Get OAuth, set as a Global
	if twAuth == nil {
		twAuth, err = getTwitterAuth()
		if err != nil || twAuth.Authorized() == false {
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
	tws = tweets{}
	err = json.Unmarshal([]byte(data), &tws)
	if err != nil {
		/* Commented out err log, as it's spamming (it complains
		 * about null values) */
		//_ = logger.Err(fmt.Sprintf("Error unmarshalling tweet JSON: '%v' '%v'", err, string(data)))

		// Was anything unmarshalled?
		if len(tws) > 0 {
			err = nil
		}
	}
	return
}

// Creates oauth using details from config.
func getTwitterAuth() (o *oauth.OAuth, err error) {
	o = new(oauth.OAuth)
	o.ConsumerKey = vafanConf.twitter.consumerKey
	o.ConsumerSecret = vafanConf.twitter.consumerSecret
	o.AccessToken = vafanConf.twitter.accessToken
	o.AccessSecret = vafanConf.twitter.accessSecret
	o.SignatureMethod = "HMAC-SHA1"
	return
}

// Functions for streaming tweets via websockets

// Websocket handler.
func streamTweets(ws *websocket.Conn) {
	writeTweetStream(ws)
}

// Connect to twitter streaming api and send lines to be written.
func writeTweetStream(w io.Writer) {
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), "info")
	stream := make(chan []byte, 1000)
	done := make(chan bool)
	client := httpstream.NewBasicAuthClient(vafanConf.twitter.user, vafanConf.twitter.password, func(line []byte) {
		stream <- line
	})
	err := client.Filter([]int64{twSaulHowardID, twConvictFilmsID}, []string{"brightonwok", "brighton wok", "convictfilms", "convict films"}, false, done)
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
		tw := tweet{}
		json.Unmarshal(line, &tw)
		if tw.User != nil {
			//println(th, " ", tweet.User.Screen_name, ": ", tweet.Text)
			w.Write([]byte(tw.Text))
		}
	}
}
