// Copyright 2012 Saul Howard. All rights reserved.

// Tests for Video functions.

package vafan

import (
	"reflect"
	"testing"
)

func TestNewVideo(t *testing.T) {
	v := newVideo()
	if reflect.TypeOf(v).String() != "*vafan.video" {
		t.Error("New video is wrong type.")
	} else {
		t.Log("NewVideo test passed.")
	}
}

func TestIsVideoNameLegal(t *testing.T) {
	good := []string{
		"mean-streets",
		"g00dfellaz",
		`องค์บาก`,
	}
	for _, n := range good {
		v := video{Name: n}
		if !v.isNameLegal() {
			t.Error("Good name declared illegal.")
		} else {
			t.Log("isNameLegal test passed for good names.")
		}
	}

	bad := []string{
		"The Lord of The Rings",
		"delete;",
	}
	for _, n := range bad {
		v := video{Name: n}
		if v.isNameLegal() {
			t.Error("Bad name declared legal.")
		} else {
			t.Log("isNameLegal test passed for bad names.")
		}
	}
}
