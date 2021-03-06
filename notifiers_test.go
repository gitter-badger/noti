package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPushbulletNotify(t *testing.T) {
	orig := struct {
		url string
		tok string
	}{
		url: pushbulletAPI,
		tok: os.Getenv(pushbulletTokEnv),
	}
	defer func() {
		pushbulletAPI = orig.url
		os.Setenv(pushbulletTokEnv, orig.tok)
	}()
	n := notification{"title", "mesg", false}

	os.Unsetenv(pushbulletTokEnv)
	if err := pushbulletNotify(n); err == nil {
		t.Error("Missing access token.")
	}

	os.Setenv(pushbulletTokEnv, "fu")
	var hitServer bool
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if r.Header.Get("Access-Token") == "" {
			t.Error("Missing access token.")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Content type should be application/json.")
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("Missing payload.")
		}
	}))
	defer ts.Close()

	pushbulletAPI = ts.URL
	if err := pushbulletNotify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("Didn't reach server.")
	}
}

func TestSlackNotify(t *testing.T) {
	orig := struct {
		url  string
		tok  string
		dest string
	}{
		url:  slackAPI,
		tok:  os.Getenv(slackTokEnv),
		dest: os.Getenv(slackDestEnv),
	}
	defer func() {
		slackAPI = orig.url
		os.Setenv(slackTokEnv, orig.tok)
		os.Setenv(slackDestEnv, orig.dest)
	}()
	n := notification{"title", "mesg", false}

	os.Unsetenv(slackTokEnv)
	if err := slackNotify(n); err == nil {
		t.Error("Missing access token.")
	}

	os.Unsetenv(slackDestEnv)
	if err := slackNotify(n); err == nil {
		t.Error("Missing destination.")
	}

	os.Setenv(slackTokEnv, "fu")
	os.Setenv(slackDestEnv, "fa")

	var hitServer bool
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if r.FormValue("token") == "" {
			t.Error("Missing access token.")
		}
		if r.FormValue("channel") == "" {
			t.Error("Missing destination.")
		}
	}))
	defer ts.Close()

	slackAPI = ts.URL
	if err := slackNotify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("Didn't reach server.")
	}
}

func TestHipChatNotify(t *testing.T) {
	orig := struct {
		url  string
		tok  string
		dest string
	}{
		url:  hipChatAPI,
		tok:  os.Getenv(hipChatTokEnv),
		dest: os.Getenv(hipChatDestEnv),
	}
	defer func() {
		hipChatAPI = orig.url
		os.Setenv(hipChatTokEnv, orig.tok)
		os.Setenv(hipChatDestEnv, orig.dest)
	}()
	n := notification{"title", "mesg", false}

	os.Unsetenv(hipChatTokEnv)
	if err := hipChatNotify(n); err == nil {
		t.Error("Missing access token.")
	}

	os.Unsetenv(hipChatDestEnv)
	if err := hipChatNotify(n); err == nil {
		t.Error("Missing destination.")
	}

	os.Setenv(hipChatTokEnv, "fu")

	var hitServer bool
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if r.Header.Get("Authorization") == "" {
			t.Error("Missing access token.")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Content type should be application/json.")
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(b) == "" {
			t.Error("Missing payload.")
		}
	}))
	defer ts.Close()

	// In real life, to calculate hipChatAPI, we need to Sprintf the
	// destination env var into the URL. This just pretends that HipChat Room
	// is at ts.URL.
	os.Setenv(hipChatDestEnv, ts.URL)
	hipChatAPI = "%s"
	if err := hipChatNotify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("Didn't reach server.")
	}
}

func TestPushoverNotify(t *testing.T) {
	orig := struct {
		url  string
		tok  string
		dest string
	}{
		url:  pushoverAPI,
		tok:  os.Getenv(pushoverTokEnv),
		dest: os.Getenv(pushoverDestEnv),
	}
	defer func() {
		slackAPI = orig.url
		os.Setenv(pushoverTokEnv, orig.tok)
		os.Setenv(pushoverDestEnv, orig.dest)
	}()
	n := notification{"title", "mesg", false}

	os.Unsetenv(pushoverTokEnv)
	if err := pushoverNotify(n); err == nil {
		t.Error("Missing access token.")
	}

	os.Unsetenv(pushoverDestEnv)
	if err := pushoverNotify(n); err == nil {
		t.Error("Missing destination.")
	}

	os.Setenv(pushoverTokEnv, "fu")
	os.Setenv(pushoverDestEnv, "fa")

	var hitServer bool
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST.")
		}

		if r.FormValue("token") == "" {
			t.Error("Missing access token.")
		}
		if r.FormValue("user") == "" {
			t.Error("Missing destination.")
		}
	}))
	defer ts.Close()

	pushoverAPI = ts.URL
	if err := pushoverNotify(n); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("Didn't reach server.")
	}
}
