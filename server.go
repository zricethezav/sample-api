package main

import (
	"os"
	"net/http"
	"log"
)

// Questions about the acceptance criteria:
// 1. "The produce database is only a single, in memory array of data and"
//	- excludes things like redis-lite, boltdb, etc?
//	- should this be persistent (lets go with yes)

// produce code should be a hash of the produce string... that way
// we dont have to walk the entire "database" to see if the code exists or not
type produce struct {
	code  string
	label string
	// price maybe should be two ints. dollar amount and change??
	price float32
}

// will need some globals as all the http actions will need channels
var db []produce

func main() {
	args := os.Args[1:]
	if args != nil {
		// load db
	} else {
		// init empty
		db := []produce{}
	}
	run(db)
}

func run(produceDB []produce) {
	// need an update channel to update/invalidate the db persistence
	// check out gob
	http.HandleFunc("/", invalidHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/fetch", fetchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func invalidHandler(w http.ResponseWriter, r *http.Request) {
	// TODO - default page 404
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	// TODO POST request to add produce to db
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Delete -- URL param `Produce Code`
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO -- GET all produce
}
