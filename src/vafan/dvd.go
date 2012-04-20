// Copyright 2012 Saul Howard. All rights reserved.

// A DVD.

package vafan

import (
	"code.google.com/p/gorilla/mux"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var errDVDNotFound = errors.New("DVD: doesn't exist")

// A dvd 
type dvd struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"` // names are unique
	Title            string            `json:"title"`
	Prices           map[string]string `json:"prices"`
	Date             time.Time         `json:"date"`
	ShortDescription string            `json:"shortDescription"`
	Description      Markdown          `json:"description"`
	URL              string            `json:"url"`
	Thumbnail        Image             `json:"thumbnail"`
	sites            []*site           // the sites that display this dvd
}

// Video constructor.
func newDVD() (d *dvd) {
	d = new(dvd)
	return
}

// Methods to implement Resource interface

func (d dvd) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(d, req, s, []string{"name", d.Name})
}

func (d dvd) GetContent(req *http.Request, s *site) (c resourceContent) {
	c.title = "DVD"
	c.description = "DVD page"
	c.content = map[string]interface{}{"dvd": d}
	return
}

func (d dvd) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		var err error
		if vars["name"] == "brighton-wok-pal" {
			d = *getBrightonWokDVD()
		} else {
			err = errDVDNotFound
		}
		if err != nil {
			if err == ErrVideoNotFound {
				notFound{}.ServeHTTP(w, r, reqU)
				return
			}
			_ = logger.Err(fmt.Sprintf("Failed to get dvd by name: %v", err))
			notFound{}.ServeHTTP(w, r, reqU)
			return
		}
		writeResource(w, r, &d, reqU)
		return
	}
}

func getBrightonWokDVD() *dvd {
	createdOn, _ := time.Parse("2006-01-02", "2008-01-01")
	desc := `## DVD
dvddvdvd`
	descMarkdown := Markdown(desc)
	thumb := Image{URL: "/img/brighton-wok/dvd.png", Width: "640", Height: "360"}
	allSites := []*site{&sites[0], &sites[1]}
	prices := map[string]string{"GBP": "9.99", "USD": "19.99"}
	return &dvd{
		ID:               "001",
		Name:             "brighton-wok-pal",
		Title:            "Brighton Wok DVD",
		Date:             createdOn,
		ShortDescription: "Brighton Wok DVD",
		Description:      descMarkdown,
		Thumbnail:        thumb,
		Prices:           prices,
		sites:            allSites,
	}
}

func getFeaturedDVDs(s *site) (dvds map[string]*dvd, err error) {
	bwok := getBrightonWokDVD()
	dvds = map[string]*dvd{"brightonWok": bwok}
	return
}
