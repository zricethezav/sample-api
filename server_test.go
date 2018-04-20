package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func loadDB() {
	for i := 0; i < 9999; i++ {
		code := strings.ToUpper(fmt.Sprintf("YRT6-72AS-K736-%04d", i))
		produce := Produce{code, "apple", 12.12}
		db.data = append(db.data, &produce)
		db.cache[produce.Code] = true
	}
}

func clearDB() {
	db.data = []*Produce{}
	db.cache = map[string]bool{}
}

func TestValidPrice(t *testing.T) {
	if validPrice(12.111) {
		t.Errorf("12.111 is not a valid price")
	}
	if validPrice(0) {
		t.Errorf("0 is not a valid price")
	}
	if validPrice(-12.111) {
		t.Errorf("-12.111 is not a valid price")
	}
	if validPrice(-12.11) {
		t.Errorf("-12.11 is not a valid price")
	}
	if !validPrice(1) {
		t.Errorf("1 is a valid price")
	}
	if !validPrice(0.1) {
		t.Errorf("0.1 is a valid price")
	}
}

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
	if nameRegexp.Match([]byte("apple---")) {
		t.Error("nameRegexp incorrectly matched 'apple---'")
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

func requestHelper(t *testing.T, handler http.HandlerFunc, method string, url string, payload io.Reader,
	expectedStatus int) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	h := http.HandlerFunc(handler)
	h.ServeHTTP(recorder, req)
	if status := recorder.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}
}

func TestAddHandler(t *testing.T) {
	sampleRequest := []byte(`{"name":"apple","code":"YRT6-72AS-K736-L4ee", "price":"12.12"}`)
	badPrice := []byte(`{"name":"apple","code":"YRT6-72AS-K736-L4ee", "price":"12.123"}`)
	badCode := []byte(`{"name":"apple","code":"YRT6-72AS-K736-L4eee", "price":"12.12"}`)
	badName := []byte(`{"name":"apple--","code":"YRT6-72AS-K736-L4ee", "price":"12.12"}`)
	badJSON := []byte(`{"name":"apple--","code":"YRT6-72AS-K736-L4ee", "price":"12.12"`)

	requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(sampleRequest), http.StatusCreated)
	requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(sampleRequest), http.StatusConflict)
	requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(badPrice), http.StatusUnprocessableEntity)
	requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(badCode), http.StatusUnprocessableEntity)
	requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(badName), http.StatusUnprocessableEntity)
	requestHelper(t, addHandler, "GET", "/add", bytes.NewReader(sampleRequest), http.StatusMethodNotAllowed)
	requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(badJSON), http.StatusUnprocessableEntity)

	// test db synchronizity... add 9999 entries asynchronously
	payloadBase := `{"name":"apple","code":"YRT6-72AS-K736-%04d", "price":"12.12"}`
	var wg sync.WaitGroup

	for i := 0; i < 9999; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := []byte(fmt.Sprintf(payloadBase, i))
			requestHelper(t, addHandler, "POST", "/add", bytes.NewReader(payload), http.StatusCreated)
		}(i)
	}
	wg.Wait()

	if len(db.data) != 10000 {
		t.Errorf("database not filled: got %d entries want 10000", len(db.data))
	}

}

func TestDeleteHandler(t *testing.T) {
	// load up database with some values
	clearDB()
	loadDB()

	// clear all entries
	var wg sync.WaitGroup
	for i := 0; i < 9999; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			url := fmt.Sprintf("/delete?code=YRT6-72AS-K736-%04d", i)
			requestHelper(t, deleteHandler, "DELETE", url, nil, http.StatusNoContent)
		}(i)
	}
	wg.Wait()

	if len(db.data) != 0 {
		t.Errorf("database filled after delete: got %d entries want 0", len(db.data))
	}

	// bad method
	loadDB()
	requestHelper(t, deleteHandler, "GET", "/delete?code=YRT6-72AS-K736-1000",
		nil, http.StatusMethodNotAllowed)
	// bad code
	requestHelper(t, deleteHandler, "DELETE", "/delete?code=YRT6-72AS-K736-10000",
		nil, http.StatusUnprocessableEntity)

	// entity not found
	requestHelper(t, deleteHandler, "DELETE", "/delete?code=YRT6-72AS-K736-1000",
		nil, http.StatusNoContent)
	requestHelper(t, deleteHandler, "DELETE", "/delete?code=YRT6-72AS-K736-1000",
		nil, http.StatusNotFound)
}

func TestFetchHandler(t *testing.T) {
	clearDB()
	// get empty
	requestHelper(t, fetchHandler, "GET", "/fetch",
		nil, http.StatusOK)

	// get full
	loadDB()
	req, err := http.NewRequest("GET", "/fetch", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(fetchHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	dbResp := []Produce{}
	json.Unmarshal(recorder.Body.Bytes(), &dbResp)
	if len(dbResp) != 9999 {
		t.Errorf("database not filled: got %d entries want 9999", len(db.data))
	}

	// verify response body. Order is determined by loadDB
	if dbResp[0].Name != "apple" {
		t.Errorf("expecting name: apple, got %s", dbResp[0].Name)
	}
	if dbResp[0].Code != "YRT6-72AS-K736-0000" {
		t.Errorf("code expecting YRT6-72AS-K736-0000, got %s", dbResp[0].Code)
	}

	// bad method
	requestHelper(t, fetchHandler, "POST", "/fetch",
		nil, http.StatusMethodNotAllowed)
}
