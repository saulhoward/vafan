// Copyright 2012 Saul Howard. All rights reserved.

// Youtube videos. 

// Takes a youtube ID and fills in properties from
// the youtube API.

// https://developers.google.com/youtube/2.0/developers_guide_protocol_video_entries

package vafan

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const singleYoutubeVideoURL = "https://gdata.youtube.com/feeds/api/videos/{id}?&v=2&key={key}"

var ErrYoutubeNotFound = errors.New("youtube: id doesn't exist")

type youtubeVideo struct {
	Id       string
	Location string
	Data     youtubeXML
}

type youtubeXML struct {
	XMLName       xml.Name           `xml:"entry"`
	Title         string             `xml:"title"`
	Statistics    youtubeStats       `xml:"statistics"`
	PublishedDate string             `xml:"published"`
	Thumbnails    []youtubeThumbnail `xml:"group>thumbnail"`
	Links         []youtubeLink      `xml:"link"`
}
type youtubeStats struct {
	Views      int `xml:"viewCount,attr"`
	Favourites int `xml:"favoriteCount,attr"`
}
type youtubeThumbnail struct {
	URL    string `xml:"url,attr"`
	Height string `xml:"height,attr"`
	Width  string `xml:"width,attr"`
	Time   string `xml:"time,attr"`
	Name   string `xml:"name,attr"`
}
type youtubeLink struct {
	Rel      string `xml:"rel,attr"`
	Location string `xml:"href,attr"`
}

func (y *youtubeVideo) FetchData() (err error) {
	youtubeDevKey, err := conf.String("default", "youtube-dev-key")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to fetch youtube dev key from config: %v", err))
		return
	}
	r := strings.NewReplacer("{id}", y.Id, "{key}", youtubeDevKey)
	res, err := http.Get(r.Replace(singleYoutubeVideoURL))
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to GET youtube URL: %v", err))
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed reading youtube response body: %v", err))
		return
	}
	err = xml.Unmarshal([]byte(data), &y.Data)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed unmarshalling youtube XML: %v", err))
		return
	}
	y.Location, err = y.getLocation()
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed setting youtube video URL: %v", err))
	}
	return
}

func (y *youtubeVideo) getDefaultThumbnail() (t youtubeThumbnail, err error) {
	maxWidth := 0
	for _, thumb := range y.Data.Thumbnails {
		width, err := strconv.Atoi(thumb.Width)
		if err == nil {
			if width > maxWidth {
				t = thumb
			}
			maxWidth = width
		}
	}
	if t.URL == "" {
		err = errors.New("youtube: default thumbnail not found")
		_ = logger.Err(fmt.Sprintf("Failed getting default youtube thumbnail: %v", err))
		return
	}
	return
}

func (y *youtubeVideo) getLocation() (l string, err error) {
	for _, link := range y.Data.Links {
		if link.Rel == "alternate" {
			l = link.Location
            return
		}
	}
	err = errors.New("youtube: video location not found")
	_ = logger.Err(fmt.Sprintf("Failed getting youtube video URL: %v", err))
	return
}

/*
Example Youtube API result for one vid:

<entry gd:etag="W/"DUMCQH47eCp7I2A9WhVRGUQ."">
  <id>tag:youtube.com,2008:video:hFSlQrB3iGY</id>
  <published>2010-03-07T01:05:03.000Z</published>
  <updated>2012-03-29T05:04:21.000Z</updated>
  <category scheme="http://schemas.google.com/g/2005#kind" term="http://gdata.youtube.com/schemas/2007#video"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/categories.cat" term="Music" label="Music"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="John"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="Lee"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="Hooker"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="serves"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="me"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="right"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="to"/>
  <category scheme="http://gdata.youtube.com/schemas/2007/keywords.cat" term="suffer"/>
  <title>John Lee Hooker - serves me right to suffer</title>
  <content type="application/x-shockwave-flash" src="https://www.youtube.com/v/hFSlQrB3iGY?version=3&f=videos&d=AQzwMvx3P_vnDfRX4Vuzv2wO88HsQjpE1a8d1GxQnGDm&app=youtube_gdata"/>
  <link rel="alternate" type="text/html" href="https://www.youtube.com/watch?v=hFSlQrB3iGY&feature=youtube_gdata"/>
  <link rel="http://gdata.youtube.com/schemas/2007#video.responses" type="application/atom+xml" href="https://gdata.youtube.com/feeds/api/videos/hFSlQrB3iGY/responses?v=2"/>
  <link rel="http://gdata.youtube.com/schemas/2007#video.ratings" type="application/atom+xml" href="https://gdata.youtube.com/feeds/api/videos/hFSlQrB3iGY/ratings?v=2"/>
  <link rel="http://gdata.youtube.com/schemas/2007#video.complaints" type="application/atom+xml" href="https://gdata.youtube.com/feeds/api/videos/hFSlQrB3iGY/complaints?v=2"/>
  <link rel="http://gdata.youtube.com/schemas/2007#video.related" type="application/atom+xml" href="https://gdata.youtube.com/feeds/api/videos/hFSlQrB3iGY/related?v=2"/>
  <link rel="http://gdata.youtube.com/schemas/2007#mobile" type="text/html" href="https://m.youtube.com/details?v=hFSlQrB3iGY"/>
  <link rel="self" type="application/atom+xml" href="https://gdata.youtube.com/feeds/api/videos/hFSlQrB3iGY?v=2"/>
  <author>
    <name>alessiobiancheri</name>
    <uri>https://gdata.youtube.com/feeds/api/users/alessiobiancheri</uri>
    <yt:userId>k_EJjrrEGrxcoO_CalvEwA</yt:userId>
  </author>
  <yt:accessControl action="comment" permission="allowed"/>
  <yt:accessControl action="commentVote" permission="allowed"/>
  <yt:accessControl action="videoRespond" permission="moderated"/>
  <yt:accessControl action="rate" permission="allowed"/>
  <yt:accessControl action="embed" permission="allowed"/>
  <yt:accessControl action="list" permission="allowed"/>
  <yt:accessControl action="autoPlay" permission="allowed"/>
  <yt:accessControl action="syndicate" permission="allowed"/>
  <gd:comments>
    <gd:feedLink rel="http://gdata.youtube.com/schemas/2007#comments" href="https://gdata.youtube.com/feeds/api/videos/hFSlQrB3iGY/comments?v=2" countHint="20"/>
  </gd:comments>
  <media:group>
    <media:category label="Music" scheme="http://gdata.youtube.com/schemas/2007/categories.cat">Music</media:category>
    <media:content url="https://www.youtube.com/v/hFSlQrB3iGY?version=3&f=videos&d=AQzwMvx3P_vnDfRX4Vuzv2wO88HsQjpE1a8d1GxQnGDm&app=youtube_gdata" type="application/x-shockwave-flash" medium="video" isDefault="true" expression="full" duration="388" yt:format="5"/>
    <media:content url="rtsp://v1.cache7.c.youtube.com/CkULENy73wIaPAlmiHewQqVUhBMYDSANFEgGUgZ2aWRlb3NyIQEM8DL8dz_75w30V-Fbs79sDvPB7EI6RNWvHdRsUJxg5gw=/0/0/0/video.3gp" type="video/3gpp" medium="video" expression="full" duration="388" yt:format="1"/>
    <media:content url="rtsp://v8.cache5.c.youtube.com/CkULENy73wIaPAlmiHewQqVUhBMYESARFEgGUgZ2aWRlb3NyIQEM8DL8dz_75w30V-Fbs79sDvPB7EI6RNWvHdRsUJxg5gw=/0/0/0/video.3gp" type="video/3gpp" medium="video" expression="full" duration="388" yt:format="6"/>
    <media:credit role="uploader" scheme="urn:youtube" yt:display="alessiobiancheri">alessiobiancheri</media:credit>
    <media:description type="plain"/>
    <media:keywords>John, Lee, Hooker, serves, me, right, to, suffer</media:keywords>
    <media:license type="text/html" href="http://www.youtube.com/t/terms">youtube</media:license>
    <media:player url="https://www.youtube.com/watch?v=hFSlQrB3iGY&feature=youtube_gdata_player"/>
    <media:thumbnail url="http://i.ytimg.com/vi/hFSlQrB3iGY/default.jpg" height="90" width="120" time="00:03:14" yt:name="default"/>
    <media:thumbnail url="http://i.ytimg.com/vi/hFSlQrB3iGY/hqdefault.jpg" height="360" width="480" yt:name="hqdefault"/>
    <media:thumbnail url="http://i.ytimg.com/vi/hFSlQrB3iGY/1.jpg" height="90" width="120" time="00:01:37" yt:name="start"/>
    <media:thumbnail url="http://i.ytimg.com/vi/hFSlQrB3iGY/2.jpg" height="90" width="120" time="00:03:14" yt:name="middle"/>
    <media:thumbnail url="http://i.ytimg.com/vi/hFSlQrB3iGY/3.jpg" height="90" width="120" time="00:04:51" yt:name="end"/>
    <media:title type="plain">John Lee Hooker - serves me right to suffer</media:title>
    <yt:duration seconds="388"/>
    <yt:uploaded>2010-03-07T01:05:03.000Z</yt:uploaded>
    <yt:videoid>hFSlQrB3iGY</yt:videoid>
  </media:group>
  <gd:rating average="5.0" max="5" min="1" numRaters="157" rel="http://schemas.google.com/g/2005#overall"/>
  <yt:statistics favoriteCount="173" viewCount="24597"/>
  <yt:rating numDislikes="0" numLikes="157"/>
</entry>
*/
