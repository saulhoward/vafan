// Copyright 2012 Saul Howard. All rights reserved.

// Vimeo videos. 

// Takes a Vimeo ID and fills in properties from
// their API.

// http://vimeo.com/api/docs/simple-api#video

package vafan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const singleVimeoVideoURLSchema = "http://vimeo.com/api/v2/video/{id}.json"

var ErrVimeoNotFound = errors.New("vimeo: id doesn't exist")

type vimeoVideo struct {
	ID       string
	Location string
	Data     vimeoJSON
}

type vimeoJSON struct {
	title                 string
	url                   string
	upload_date           string
	thumbnail_small       string
	thumbnail_medium      string
	thumbnail_large       string
	stats_number_of_plays int
	stats_number_of_likes int
	width                 int
	height                int
}

func (v *vimeoVideo) FetchData() (err error) {
	r := strings.NewReplacer("{id}", v.ID)
	res, err := http.Get(r.Replace(singleVimeoVideoURLSchema))
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to GET Vimeo URL: %v", err))
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed reading Vimeo response body: %v", err))
		return
	}
	err = json.Unmarshal([]byte(data), &v.Data)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed unmarshalling Vimeo JSON: %v", err))
		return
	}
	// Set default url
	v.Location = v.Data.url
	return
}

func (v *vimeoVideo) getDefaultThumbnail() (i Image, err error) {
	i = Image{URL: v.Data.url, Width: "640"}
	return
}

/*
Example API result for one vid:

[
{
    id: 35687624,
    title: "Martin Parr: Teddy Gray's Sweet Factory",
    description: "Magnum photographer, Martin Parr returns to using a film camera in this wonderfully engaging documentary about Teddy Gray’s sweet factory in Dudley in the West Midlands.   <br />
    <br />
    Established in 1826, Teddy Gray’s has always been a family owned and run business. Five generations have worked and contributed towards the business of keeping the traditional, hand-made methods of sweet making alive.<br />
    <br />
    The film is part of the Black Country Stories body of work commissioned by Multistory to document life in the Black Country by capturing and celebrating the unique mix of communities living in the area and of existing traditional Black Country life.   <br />
    <br />
    A Multistory Production, 2011<br />
    <br />
    Filmed and directed by Martin Parr<br />
    Sound by Andrew Yarme<br />
    Film editing by Darren Flaxstone<br />
    Music by Rob Dunstone<br />
    Produced by Multistory<br />
    <br />
    © Martin Parr, 2011.  All rights reserved.<br />
    <br />
    Thanks to Edward Gray, Betty Guest and all the workers at Teddy Gray’s.<br />
    <br />
    WARNING: all rights of this DVD are reserved and it is strictly prohibited to use this DVD other than for private viewing.  Duplication, unless authorised, is strictly prohibited.",
    url: "http://vimeo.com/35687624",
    upload_date: "2012-01-26 08:02:06",
    thumbnail_small: "http://b.vimeocdn.com/ts/244/194/244194997_100.jpg",
    thumbnail_medium: "http://b.vimeocdn.com/ts/244/194/244194997_200.jpg",
    thumbnail_large: "http://b.vimeocdn.com/ts/244/194/244194997_640.jpg",
    user_name: "Multistory",
    user_url: "http://vimeo.com/multistory",
    user_portrait_small: "http://b.vimeocdn.com/ps/214/573/2145739_30.jpg",
    user_portrait_medium: "http://b.vimeocdn.com/ps/214/573/2145739_75.jpg",
    user_portrait_large: "http://b.vimeocdn.com/ps/214/573/2145739_100.jpg",
    user_portrait_huge: "http://b.vimeocdn.com/ps/214/573/2145739_300.jpg",
    stats_number_of_likes: 9,
    stats_number_of_plays: 526,
    stats_number_of_comments: 0,
    duration: 1200,
    width: 640,
    height: 360,
    tags: "Martin Parr, Black Country Stories, Teddy Gray's Sweet Factory, Dudley, Multistory",
    embed_privacy: "anywhere"
}
]
*/
