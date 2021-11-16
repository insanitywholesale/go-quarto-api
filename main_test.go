package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//TODO: use functions for things that are reused in tests

// Test creating a user
func TestCreateUser(t *testing.T) {
	// run test http server
	testServer := httptest.NewServer(setupRouter())
	// define URL
	testURL := testServer.URL + "/user/register"
	// create some data in the form of an io.Reader from a string of json
	jsonData := []byte(`{"username": "myself", "password": "mypasswd"}`)
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
	// close test http server
	testServer.Close()
}

func TestCreateGame(t *testing.T) {
	// run test http server
	testServer := httptest.NewServer(setupRouter())
	// define URL
	testURL := testServer.URL + "/user/register"

	// basically TestCreateUser without err checks
	jsonUserData := []byte(`{"username": "myself", "password": "mypasswd"}`)
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
	// close test http server
	testServer.Close()
}

func TestInviteAndJoin(t *testing.T) {
	// run test http server
	testServer := httptest.NewServer(setupRouter())

	// redo things from previous tests but without err checks
	// define URL
	testURL := testServer.URL + "/user/register"
	// basically TestCreateUser without err checks
	jsonUserData := []byte(`{"username": "myself", "password": "mypasswd"}`)
	r, _ := http.Post(testURL, "application/json", bytes.NewBuffer(jsonUserData))
	defer r.Body.Close()
	b, _ := io.ReadAll(r.Body)
	u := &UserId{}
	_ = json.Unmarshal(b, u)
	// basically TestCreateUser without err checks for user to be invited
	jsonUserData = []byte(`{"username": "notme", "password": "notmypasswd"}`)
	r, _ = http.Post(testURL, "application/json", bytes.NewBuffer(jsonUserData))
	defer r.Body.Close()
	b, _ = io.ReadAll(r.Body)
	u2 := &UserId{}
	_ = json.Unmarshal(b, u2)
	// basically TestCreateGame without err checks
	testURL = testServer.URL + "/game"
	jsonGameData := []byte(`{"username": "` + u.UserName + `", "user_id": "` + u.UserId + `"}`)
	r, _ = http.Post(testURL, "application/json", bytes.NewBuffer(jsonGameData))
	defer r.Body.Close()
	b, _ = io.ReadAll(r.Body)
	g := &Game{}
	_ = json.Unmarshal(b, g)

	// test invite
	// change URL
	testURL = testServer.URL + "/game/" + g.GameId + "/invite/" + u2.UserName
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
	// set success message
	successMessage := `{"message": "success"}`
	// check if body has success message
	if string(body) != successMessage {
		t.Error("inviting user did not yield success message")
	}

	// test game creator join
	// change URL
	testURL = testServer.URL + "/game/" + g.GameId + "/join"
	// create some data in the form of an io.Reader from a string of json
	jsonData = []byte(`{"username": "` + u.UserName + `", "user_id": "` + u.UserId + `"}`)
	// do a simple Post request with the above data
	res, err = http.Post(testURL, "application/json", bytes.NewBuffer(jsonData))
	// check for request errors
	if err != nil {
		t.Error("POST error:", err)
	}
	// be responsible and close the response some time
	defer res.Body.Close()
	// log response
	t.Log(res)

	// save response body to check later
	body, err = io.ReadAll(res.Body)
	// check for response body read errors
	if err != nil {
		t.Error("resp.Body error:", err)
	}
	// log response body
	t.Log(string(body))

	// test invited player join
	// create some data in the form of an io.Reader from a string of json
	jsonData = []byte(`{"username": "` + u2.UserName + `", "user_id": "` + u2.UserId + `"}`)
	// do a simple Post request with the above data
	res, err = http.Post(testURL, "application/json", bytes.NewBuffer(jsonData))
	// check for request errors
	if err != nil {
		t.Error("POST error:", err)
	}
	// be responsible and close the response some time
	defer res.Body.Close()
	// log response
	t.Log(res)

	// save response body to check later
	body, err = io.ReadAll(res.Body)
	// check for response body read errors
	if err != nil {
		t.Error("resp.Body error:", err)
	}
	// log response body
	t.Log(string(body))

	// close test http server
	//testServer.Close()
}

//func TestPlayInGame(t *testing.T) {
//	t.Error("test not written")
//}
