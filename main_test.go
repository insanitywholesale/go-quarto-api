package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test creating a user
func TestCreateUser(t *testing.T) {
	// run test http server
	testServer := httptest.NewServer(setupRouter())
	// define URL
	testURL := testServer.URL + "/user/register"
	// create some data in the form of an io.Reader from a string of json
	jsonData := []byte(`{"username": "myself", "password": "my password"}`)
	// do a simple Post request with the above data
	res, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonData))
	// check for request errors
	if err != nil {
		t.Error("POST error:", err)
	}
	// be responsible and close the response some time
	defer res.Body.Close()
	// log response
	t.Log(res)

	// save response body to check later
	body, err := io.ReadAll(res.Body)
	// check for response body read errors
	if err != nil {
		t.Error("resp.Body error:", err)
	}
	// log response body
	t.Log(string(body))

	// response should contain json that can maps to the UserId type
	// set up empty UserId
	u := &UserId{}
	// try to unmarshal
	err = json.Unmarshal(body, u)
	// check for unmarshaling errors
	if err != nil {
		t.Error("unmarshal error:", err)
	}
	// log UserId
	t.Log(u)
}

func TestCreateGame(t *testing.T) {
	// run test http server
	testServer := httptest.NewServer(setupRouter())
	// define URL
	testURL := testServer.URL + "/user/register"

	// basically TestCreateUser without err checks
	jsonUserData := []byte(`{"username": "myself", "password": "my password"}`)
	res, _ := http.Post(testURL, "application/json", bytes.NewBuffer(jsonUserData))
	defer res.Body.Close()
	bod, _ := io.ReadAll(res.Body)
	u := &UserId{}
	_ = json.Unmarshal(bod, u)

	// change URL
	testURL = testServer.URL + "/game"
	// create some data in the form of an io.Reader from a string of json
	jsonData := []byte(`{"username": "` + u.UserName + `", "user_id": "` + u.UserId + `"}`)
	// do a simple Post request with the above data
	res, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonData))
	// check for request errors
	if err != nil {
		t.Error("POST error:", err)
	}
	// be responsible and close the response some time
	defer res.Body.Close()
	// log response
	t.Log(res)

	// save response body to check later
	body, err := io.ReadAll(res.Body)
	// check for response body read errors
	if err != nil {
		t.Error("resp.Body error:", err)
	}
	// log response body
	t.Log(string(body))

	// response should contain json that can maps to the Game type
	// set up empty Game
	g := &Game{}
	// try to unmarshal
	err = json.Unmarshal(body, g)
	// check for unmarshaling errors
	if err != nil {
		t.Error("unmarshal error:", err)
	}
	// log GameId
	t.Log(g)
}
