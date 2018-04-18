package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// Produce represents a single produce entry in the database
type Produce struct {
	Code  string
	Name  string
	Price float32 `json:",string"`
}

// produceDB is the database which consists of a slice of produce pointers,
// a reader/writer mutex, and a cache to prevent unnecessary database walks.
type produceDB struct {
	data  []*Produce
	cache map[string]bool
	lock  sync.RWMutex
}

// Globals include the database, cache, and regex. Need to be
// accessed by goroutines.
var (
	db         produceDB
	nameRegexp *regexp.Regexp
	codeRegexp *regexp.Regexp
)

func init() {
	nameRegexp = regexp.MustCompile("[0-9A-Za-z]$")
	codeRegexp = regexp.MustCompile("[0-9A-Za-z]{4}-[0-9A-Za-z]{4}-[0-9A-Za-z]{4}-[0-9A-Za-z]{4}$")
	db = produceDB{}
	db.cache = map[string]bool{}
}

func main() {
	log.Println("Starting gannet-market-api")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/fetch", fetchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// AddHandler is responsible for adding a produce entry to the database.
// This function accepts POST request and expects a json of following this criteria:
// 	 - name: alphanumeric and case insensitive
// 	 - produce codes: alphanumeric and case insensitive and are sixteen
// 	   characters long, with dashes separating each four character group
// 	 - price: number with up to 2 decimal places
// Sample add request:
// 	 $ curl -X POST -d '{"name":"apple","code":"YRT6-72AS-K736-L4AR", "price": "12.12"}' localhost:8080/add
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "/add requires POST", http.StatusMethodNotAllowed)
		return
	}
	var p Produce
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "unable to process request", http.StatusUnprocessableEntity)
		return
	}
	if !(nameRegexp.Match([]byte(p.Name))) {
		http.Error(w, "invalid name", http.StatusUnprocessableEntity)
		return
	}
	if !(codeRegexp.Match([]byte(p.Code))) {
		http.Error(w, "invalid code", http.StatusUnprocessableEntity)
		return
	}

	// handle case insensitivity
	// TODO test this
	p.Name = strings.ToLower(p.Name)
	p.Code = strings.ToLower(p.Code)

	err = db.add(&p)
	if err != nil {
		http.Error(w, "entry already exists", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// DeleteHandler is responsible for removing a produce entry from the database.
// This function accepts DELETE requests and expects a query param `code`
// Sample delete request:
// 	 $  curl -X "DELETE" localhost:8080/delete?code=YRT6-72AS-K736-L4ee
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "/delete requires DELETE", http.StatusMethodNotAllowed)
		return
	}
	code := r.URL.Query().Get("code")
	if !(codeRegexp.Match([]byte(code))) {
		http.Error(w, "invalid code", http.StatusUnprocessableEntity)
		return
	}

	// handle case insensitivity
	code = strings.ToLower(code)

	err := db.delete(code)
	if err != nil {
		http.Error(w, "entry does not exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// FetchHandler is responsible for reporting all the entries in the database.
// This function accepts GET requests.
// Sample fetch request:
// 	$  curl -X GET 0.0.0.0:8080/fetch
func fetchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "/fetch requires GET", http.StatusMethodNotAllowed)
		return
	}
	resp, err := json.Marshal(db.data)
	if err != nil {
		http.Error(w, "failed to create entry", http.StatusForbidden)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Add is a helper function called from addHandler and is responsible for adding a produce
// entry to the produce database. This function first grabs a reader lock and
// checks the db cache to see if the produce entry exists. If it does, unlock the reader and return an error.
// If the entry does not exist, unlock the reader, grab a write lock and add the produce to the database and update
// the cache.
func (db *produceDB) add(produce *Produce) error {
	// ensure readers while checking cache
	db.lock.RLock()
	if exists, _ := db.cache[produce.Code]; exists {
		db.lock.RUnlock()
		return fmt.Errorf("code already exists")

	}
	db.lock.RUnlock()

	// produce does not exist, need to grab a lock for the write
	db.lock.Lock()
	defer db.lock.Unlock()
	// update database and cache
	db.data = append(db.data, produce)
	db.cache[produce.Code] = true
	return nil
}

// Delete is responsible for removing a produce entry from the produce database.
// This function errs on the side of caution and grabs a write lock right away
// to avoid any data races taking the form of bad indexing when removing the entry
// from our 'database' slice. If the entry exists in the database we remove it and update
// the cache to reflect the change. If the entry does not exist, return an error.
func (db *produceDB) delete(code string) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	for i, produce := range db.data {
		if produce.Code == code {
			// remove from db, update cache
			copy(db.data[i:], db.data[i+1:])
			db.data[len(db.data)-1] = nil
			db.data = db.data[:len(db.data)-1]
			db.cache[code] = false
			return nil
		}
	}
	return fmt.Errorf("entry does not exist")
}
