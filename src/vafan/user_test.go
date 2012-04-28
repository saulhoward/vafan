// Copyright 2012 Saul Howard. All rights reserved.

// Tests for User functions.

package vafan

import (
	"reflect"
	"testing"
)

// Tests.

func TestNewUser(t *testing.T) {
	u := NewUser()
	if reflect.TypeOf(u).String() != "*vafan.user" {
		t.Error("New User is wrong type.")
	} else if u.ID == "" {
		t.Error("New User has no ID.")
	} else if u.Role != defaultUserRole {
		t.Error("New User role is not default.")
	} else {
		t.Log("NewUser test passed.")
	}
}

func TestCreateSalt(t *testing.T) {
	s := createSalt()
	if s == "" {
		t.Error("Failed to create salt.")
	} else if len(s) != 16 {
		t.Error("Salt created not 16 characters long.")
	} else {
		t.Log("createSalt test passed.")
	}
}

func TestHashPassword(t *testing.T) {
	const pass = "password123"
	const salt = "1234567890123456"
	hp := hashPassword(pass, salt)
	if hp == "" {
		t.Error("Failed to create hashed password.")
	} else if pass == hp {
		t.Error("Failed to hash password.")
	} else {
		t.Log("hashPassword test passed.")
	}
}

func TestGetUser(t *testing.T) {
	const id = "1234567890"
	u := GetUser(id)
	if reflect.TypeOf(u).String() != "*vafan.user" {
		t.Error("Get User is wrong type.")
	} else if u.ID != id {
		t.Error("Get User has wrong ID.")
	} else if u.Role != defaultUserRole {
		t.Error("Get User role is not default.")
	} else {
		t.Log("GetUser test passed.")
	}
}

func TestGetUserForUserInfo(t *testing.T) {
	userInfo := map[string]string{
		"Id":           "1234567890",
		"Username":     "bob",
		"EmailAddress": "bob@example.com",
		"Role":         "admin",
	}
	u, err := getUserForUserInfo(userInfo)

	if err != nil {
		t.Error("Get user for user info returned error.")
	} else if reflect.TypeOf(u).String() != "*vafan.user" {
		t.Error("Get user for user info is wrong type.")
	} else if u.ID != userInfo["Id"] {
		t.Error("Get user for user info has wrong ID.")
	} else if u.Username != userInfo["Username"] {
		t.Error("Get user for user info has wrong Username.")
	} else if u.Role != userInfo["Role"] {
		t.Error("Get user for user info role is wrong.")
	} else if u.EmailAddress != userInfo["EmailAddress"] {
		t.Error("Get user for user info emailAddress is wrong.")
	} else {
		t.Log("getUserForUserInfo test with correct input passed.")
	}

	// Now, test with incorrect input (no ID).

	badUserInfo := map[string]string{
		"Username":     "bob",
		"EmailAddress": "bob@example.com",
		"Role":         "admin",
	}
	_, err = getUserForUserInfo(badUserInfo)
	if err == nil {
		t.Error("Get user for user info returned no error for incorrect input.")
	} else {
		t.Log("getUserForUserInfo test with incorrect input passed.")
	}
}

func TestNewUUID(t *testing.T) {
	u := newUUID()
	if u == "" {
		t.Error("Failed to create UUID.")
	} else if len(u) != 36 {
		t.Error("UUID created not 36 characters long.")
	} else {
		t.Log("newUUID test passed.")
	}
}
