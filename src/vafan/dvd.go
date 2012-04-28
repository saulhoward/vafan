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
var errDVDStockistNotFound = errors.New("DVD Stockist: doesn't exist")

// A DVD 
type dvd struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"` // names are unique
	Title            string            `json:"title"`
	Prices           map[string]string `json:"prices"`
	Date             time.Time         `json:"date"`
	ShortDescription string            `json:"shortDescription"`
	Description      Markdown          `json:"description"`
	URL              string            `json:"url"`
	LargeImage       Image             `json:"largeImage"`
	Thumbnail        Image             `json:"thumbnail"`
	sites            []*site           // the sites that display this dvd
	ThankUser        bool              `json:"thankUser"` // true if the user has already purchased
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

func (d dvd) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		var err error
		if vars["name"] == "brighton-wok-first-edition-pal" {
			d = *getBrightonWokDVD()
		} else {
			err = errDVDNotFound
		}
		if err != nil {
			if err == ErrVideoNotFound {
				notFound{}.ServeHTTP(w, r, reqU)
				return
			}
			logger.Err(fmt.Sprintf("Failed to get dvd by name: %v", err))
			notFound{}.ServeHTTP(w, r, reqU)
			return
		}

		// Thankyou query string, used after user has purchased
		thanks := r.URL.Query().Get("thankyou")
		if thanks != "" {
			d.ThankUser = true
		}
		d.URL = d.GetURL(r, nil).String()
		res := Resource{
			title:       d.Title,
			description: d.ShortDescription,
		}
		res.content = make(resourceContent)
		res.content["dvd"] = d
		res.write(w, r, &d, reqU)
		return
	}
}

func getBrightonWokDVD() *dvd {
	createdOn, _ := time.Parse("2006-01-02", "2008-01-01")
	desc := `## The Brighton Wok DVD
dvddvdvd`
	descMarkdown := Markdown(desc)
	lgImg := Image{URL: "/img/brighton-wok/dvd/box.png", Width: "346", Height: "476"}
	thumb := Image{URL: "/img/brighton-wok/dvd/box-174x240.png", Width: "174", Height: "240"}
	allSites := []*site{&sites[0], &sites[1]}
	prices := map[string]string{"GBP": "9.99", "USD": "19.99"}

	return &dvd{
		ID:               "001",
		Name:             "brighton-wok-first-edition-pal",
		Title:            "Brighton Wok: The DVD",
		Date:             createdOn,
		ShortDescription: "The first edition of Brighton Wok: The Legend of Ganja Boxing on DVD-5.",
		Description:      descMarkdown,
		LargeImage:       lgImg,
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

// DVD Stockists

// A DVD stockist
type dvdStockist struct {
	Name        string   `json:"name"` // names are unique
	Title       string   `json:"title"`
	Description Markdown `json:"description"`
	MapImage    Image    `json:"description"`
	Location    location `json:"location"`
	DVD         *dvd     `json:"dvd"`
	Website     website  `json:"website"`
	Address     address  `json:"address"`
	URL         string   `json:"url"`
}

type dvdStockists struct {
	DVD       *dvd           `json:"dvd"`
	Stockists []*dvdStockist `json:"stockists"`
}

// A location
type location struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

// A website
type website struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// A address
type address struct {
	ShortAddress string `json:"shortAddress"`
}

func getBrightonWokDVDStockists() dvdStockists {
	bwok := *getBrightonWokDVD()
	mktplDesc := `*The Marketplace* is a legendary Brighton headshop in the Lanes, established in 1967.`

	marketplace := dvdStockist{
		Name:        "the-marketplace-brighton",
		Title:       "The Marketplace",
		Description: Markdown(mktplDesc),
		Location:    location{Lat: "40", Long: "-40"},
		MapImage:    Image{URL: "/img/brighton-wok/maps/The-Marketplace-Map.png", Width: "600", Height: "420"},
		Website:     website{URL: "http://www.marketplace-brighton.co.uk/", Title: "MarketplaceBrighton.com"},
		Address:     address{ShortAddress: "Brighton, UK"},
	}
	marketplace.DVD = &bwok
	timeslip := dvdStockist{
		Name:     "timeslip-brighton",
		Title:    "Timeslip",
		Location: location{Lat: "40", Long: "-40"},
		MapImage: Image{URL: "/img/brighton-wok/maps/Timeslip-Map.png", Width: "600", Height: "442"},
		Address:  address{ShortAddress: "Brighton, UK"},
	}
	timeslip.DVD = &bwok
	stockists := []*dvdStockist{&marketplace, &timeslip}
	return dvdStockists{Stockists: stockists, DVD: &bwok}
}

// DVD Stockist resource

func (d dvdStockist) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(d, req, s, []string{"name", d.DVD.Name, "dvdStockist", d.Name})
}

func (d dvdStockist) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		var err error
		if vars["name"] == "brighton-wok-first-edition-pal" {
			for _, s := range getBrightonWokDVDStockists().Stockists {
				if vars["dvdStockist"] == s.Name {
					d = *s
				}
			}
			if d.Name == "" {
				err = errDVDStockistNotFound
				notFound{}.ServeHTTP(w, r, reqU)
				return
			}

		} else {
			err = errDVDNotFound
		}
		if err != nil {
			if err == ErrVideoNotFound {
				notFound{}.ServeHTTP(w, r, reqU)
				return
			}
			logger.Err(fmt.Sprintf("Failed to get dvd by name: %v", err))
			notFound{}.ServeHTTP(w, r, reqU)
			return
		}

		d.DVD.URL = d.DVD.GetURL(r, nil).String()
		res := Resource{
			title:       d.Title,
			description: d.Title,
		}
		res.content = make(resourceContent)
		res.content["stockist"] = d
		res.write(w, r, &d, reqU)
		return
	}
	return
}

// All DVD Stockists resource

func (d dvdStockists) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(d, req, s, []string{"name", d.DVD.Name})
}

func (d dvdStockists) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		var err error
		if vars["name"] == "brighton-wok-first-edition-pal" {
			d.DVD = getBrightonWokDVDStockists().DVD
			d.Stockists = getBrightonWokDVDStockists().Stockists
		} else {
			err = errDVDNotFound
		}
		if err != nil {
			if err == ErrVideoNotFound {
				notFound{}.ServeHTTP(w, r, reqU)
				return
			}
			logger.Err(fmt.Sprintf("Failed to get dvd by name: %v", err))
			notFound{}.ServeHTTP(w, r, reqU)
			return
		}

		d.DVD.URL = d.DVD.GetURL(r, nil).String()
		for _, s := range d.Stockists {
			s.URL = s.GetURL(r, nil).String()
		}

		res := Resource{
			title:       "Where to buy the " + d.DVD.Title,
			description: "Shops and businesses where you can buy " + d.DVD.Title,
		}
		res.content = make(resourceContent)
		res.content["stockists"] = d.Stockists
		res.content["dvd"] = d.DVD
		res.write(w, r, &d, reqU)
		return
	}
	return
}
