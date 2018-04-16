package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"regexp"
	"sync"
	_"io/ioutil"
	"encoding/json"
)

type Produce struct {
	Code  string
	Name string
	Price float32 `json:",string"`
}

type produceDB struct {
	data []*Produce
	lock sync.Mutex
}

var (
	db produceDB
	dbCache map[string]int
	nameRegexp *regexp.Regexp
	codeRegexp *regexp.Regexp
)

func init() {
	nameRegexp = regexp.MustCompile("[0-9A-Za-z]")
	codeRegexp = regexp.MustCompile("[0-9A-Za-z]{4}-[0-9A-Za-z]{4}-[0-9A-Za-z]{4}-[0-9A-Za-z]{4}$")
	db = produceDB{}
	dbCache = map[string]int{}

}

func main() {
	args := os.Args[1:]
	if args != nil {
		// load db
	}
	log.Println("Starting gannet-market-api service")

	http.HandleFunc("/", invalidHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/fetch", fetchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func invalidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "404 not found")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, "/add accepts POST requests")
		return
	}
	var p Produce
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprint(w, "malformed request")
		return
	}
	if !(nameRegexp.Match([]byte(p.Name))) {
		fmt.Fprint(w, "invalid name")
		return
	}
	if !(codeRegexp.Match([]byte(p.Code))) {
		fmt.Fprint(w, "invalid code")
		return
	}
	db.add(&p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Delete -- URL param `Produce Code`
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO -- GET all produce
}

func (db *produceDB) add(produce *Produce){
	// check if produce/price currently exists
	// THIS ASSUMES PRODUCE CODE PERSISTS. aka produce codes do not change between sessions
	var (
		idx int
		updatePrice bool
		exists bool
	)

	if idx, exists = dbCache[produce.Name]; exists {
		if db.data[idx].Price != produce.Price {
			// check if price is the same
			updatePrice = true
		} else {
			// produce already exists at the same price
			return
		}
	}

	if updatePrice {
		// update price
		db.lock.Lock()
		defer db.lock.Unlock()
		db.data[idx].Price = produce.Price
	} else {
		// add produce for the first time
		db.lock.Lock()
		defer db.lock.Unlock()
		db.data = append(db.data, produce)
		idx = len(db.data) - 1
		dbCache[produce.Name] = idx
	}
}

// helper scripts
// curl -H "Content-Type: application/json" -X POST -d '{"name":"apple","code":"YRT6-72AS-K736-L4AR", "price": "12.12"}' localhost:8080
