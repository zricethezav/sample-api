package main

import (
	"testing"
	"net/http"
	"bytes"
	"net/http/httptest"
)

func TestRegex(t *testing.T) {
	if !nameRegexp.Match([]byte("apple")) {
		t.Error("nameRegexp failed matching apple")
	}
	if nameRegexp.Match([]byte("---")) {
		t.Error("nameRegexp incorrectly matched '---'")
	}
	if nameRegexp.Match([]byte("")) {
		t.Error("nameRegexp incorrectly matched the empty string")
	}


	if codeRegexp.Match([]byte("")) {
		t.Error("codeRegexp incorrectly matched the empty string")
	}
	if !codeRegexp.Match([]byte("YRT6-72AS-K736-L4AR")) {
		t.Error("codeRegexp failed to match 'YRT6-72AS-K736-L4AR'")
	}
	if codeRegexp.Match([]byte("YRT6-72AS*-K736-L4AR")) {
		t.Error("codeRegexp incorrectly matched 'YRT6-72AS*-K736-L4AR'")
	}
	if codeRegexp.Match([]byte("YRT6-72AS-K736-L4ARa")) {
		t.Error("codeRegexp incorrectly matched 'YRT6-72AS*-K736-L4ARa'")
	}
	if codeRegexp.Match([]byte("72AS-K736-L4AR")) {
		t.Error("codeRegexp incorrectly matched '72AS-K736-L4AR'")
	}
}

func TestAddHandler(t *testing.T) {
	sampleRequest := []byte(`{"name":"apple","code":"YRT6-72AS-K736-L4ee", "price":"12.12"}`)
	req, err := http.NewRequest("POST", "/add", bytes.NewReader(sampleRequest))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(addHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusCreated{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestDeleteHandler(t *testing.T) {

}

func TestFetchHandler(t *testing.T) {

}

func TestLoad(t *testing.T) {
	// add, get, delete cycle for some range
}
