package main

import (
	"testing"
	"net/http"
	"bytes"
	"net/http/httptest"
	"fmt"
	"sync"
)

func loadDB() {
	for i := 0; i < 9999; i++ {
		code := fmt.Sprintf("YRT6-72AS-K736-%04d", i)
		produce := Produce{code, "apple", 12.12}
		db.data = append(db.data, &produce)
		db.cache[produce.Code] = true
	}
}

func clearDB() {
	db.data = []*Produce{}
	db.cache = map[string]bool{}
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

func TestAddHandler(t *testing.T) {
	sampleRequest := []byte(`{"name":"apple","code":"YRT6-72AS-K736-L4ee", "price":"12.12"}`)
	badCode := []byte(`{"name":"apple","code":"YRT6-72AS-K736-L4eee", "price":"12.12"}`)
	badName := []byte(`{"name":"apple--","code":"YRT6-72AS-K736-L4ee", "price":"12.12"}`)
	badJson := []byte(`{"name":"apple--","code":"YRT6-72AS-K736-L4ee", "price":"12.12"`)

	// valid request
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

	// try adding duplicate
	req, err = http.NewRequest("POST", "/add", bytes.NewReader(sampleRequest))
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusConflict{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	// bad produce code
	req, err = http.NewRequest("POST", "/add", bytes.NewReader(badCode))
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnprocessableEntity{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	// bad produce name
	req, err = http.NewRequest("POST", "/add", bytes.NewReader(badName))
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnprocessableEntity{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	// bad method
	req, err = http.NewRequest("GET", "/add", bytes.NewReader(badName))
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	// bad json
	req, err = http.NewRequest("POST", "/add", bytes.NewReader(badJson))
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnprocessableEntity{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	// test db synchronizity... add 9999 entries asynchronously
	payloadBase := `{"name":"apple","code":"YRT6-72AS-K736-%04d", "price":"12.12"}`
	var wg sync.WaitGroup

	for i := 0; i < 9999; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := []byte(fmt.Sprintf(payloadBase, i))
			req, err := http.NewRequest("POST", "/add", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(addHandler)
			handler.ServeHTTP(recorder, req)
			if status := recorder.Code; status != http.StatusCreated {
				fmt.Println(recorder.Body)
				t.Errorf("%d handler returned wrong status code: got %v want %v",
					i, status, http.StatusCreated)
			}
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

	var wg sync.WaitGroup
	for i := 0; i < 9999; i++ {
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			url := fmt.Sprintf("/delete?code=YRT6-72AS-K736-%04d", i)
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(deleteHandler)
			handler.ServeHTTP(recorder, req)
			if status := recorder.Code; status != http.StatusNoContent{
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusNotFound)
			}
		}(i)
	}
	wg.Wait()

	if len(db.data) != 0 {
		t.Errorf("database filled after delete: got %d entries want 0", len(db.data))
	}

	// bad method
	loadDB()
	req, err := http.NewRequest("GET", "/delete?code=YRT6-72AS-K736-1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusMethodNotAllowed{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	// bad key
	req, err = http.NewRequest("DELETE", "/delete?code=YRT6-72AS-K736-10000", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler = http.HandlerFunc(deleteHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnprocessableEntity{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	// entity not found
	req, err = http.NewRequest("DELETE", "/delete?code=YRT6-72AS-K738-1000", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler = http.HandlerFunc(deleteHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNotFound{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestFetchHandler(t *testing.T) {
	clearDB()

	// get empty
	req, err := http.NewRequest("GET", "/fetch", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(fetchHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// get full
	loadDB()
	req, err = http.NewRequest("GET", "/fetch", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	handler = http.HandlerFunc(fetchHandler)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
