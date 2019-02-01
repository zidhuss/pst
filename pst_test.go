package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/zidhuss/pst/db"
)

func TestHome(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}
	rec := httptest.NewRecorder()

	handler(&app{}).ServeHTTP(rec, req)

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

func TestNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/abc", nil)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}
	rec := httptest.NewRecorder()

	handler(&app{db.CreatePasteDatabase("pst.db")}).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	if len(body) > 0 {
		t.Fatalf("expected to receive empty body, instead received \"%s\"", body)
	}
}

func TestUploadAndGetRawString(t *testing.T) {
	rawString := "testing raw string"
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("could not make upload request: %v", err)
	}
	req.Form = url.Values{"f:1": {rawString}}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler(&app{db.CreatePasteDatabase("pst.db")}).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read request path response: %v", err)
	}

	// Regex to match url with random ending
	r := regexp.MustCompile(`/[A-Za-z1-9]+$`)
	if !r.Match(bodyBytes) {
		t.Fatalf("didn't received valid URL, got \"%s\"", string(bodyBytes))
	}

	req, err = http.NewRequest("GET", string(bodyBytes), nil)
	if err != nil {
		t.Fatalf("could not make file retrieval request: %v", err)
	}

	handler(&app{}).ServeHTTP(rec, req)

	res = rec.Result()
	defer res.Body.Close()
	bodyBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read paste response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	if string(bodyBytes) != rawString {
		t.Fatalf("incorrect response: expected \"%s\", got \"%s\"", rawString, string(bodyBytes))
	}
}
