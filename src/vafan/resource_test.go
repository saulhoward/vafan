// Copyright 2012 Saul Howard. All rights reserved.

// Tests for the Resource interface and helper functions.

package vafan

import (
	"net/http"
	"net/url"
	"testing"
)

// A Mock ResourceServer, for testing.

type dummyResource struct {
}

func (d dummyResource) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(d, req, s, nil)
}

func (d dummyResource) ServeHTTP(w http.ResponseWriter, r *http.Request, u *user) {
	return
}

// Tests.

func TestResourceName(t *testing.T) {
	if resourceName(dummyResource{}) != "dummyResource" {
		t.Error("resourceName did not work as expected.")
	} else {
		t.Log("resourceName test passed.")
	}
}
