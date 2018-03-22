package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHome(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}
	rec := httptest.NewRecorder()

	home(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	if string(body) != Help {
		t.Fatalf("incorrect response: expected %s, got %s", Help, body)
	}
}

func TestUploadRawString(t *testing.T) {
	rawString := "testing raw string"
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}
	req.Form = url.Values{"f1": {rawString}}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	home(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
	bodyString := string(bodyBytes)

	if bodyString != rawString {
		t.Fatalf("incorrect response: expected %s, got %s", rawString, bodyString)
	}
}

func TestRawGet(t *testing.T) {
	t.Errorf("Not implemented")
}

func TestConsoleHighlight(t *testing.T) {
	t.Errorf("Not implemented")
}
