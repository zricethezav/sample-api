package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"regexp"
)

type produce struct {
	code  string
	label string
	price float32
}

// will need some globals as all the http actions will need channels
// Requirements state that the produce "DB" need only be an in memory
// array of data
// ASK ABOUT THIS!!!
var (
	db []produce // needs to be protected by a lock
	nameRegexp *regexp.Regexp
	codeRegexp *regexp.Regexp
)

func init() {
	nameRegexp = regexp.MustCompile("[0-9A-Za-z]")
	codeRegexp = regexp.MustCompile("[0-9A-Za-z]{4}-[0-9A-Za-z]{4}-[0-9A-Za-z]{4}-[0-9A-Za-z]{4}")
	db = []produce{}
}

func main() {
	args := os.Args[1:]
	if args != nil {
		// load db
	}
	fmt.Println("Running")
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
	errHandler(w, "404 not found")
}

func errHandler(w http.ResponseWriter, msg string) {
	fmt.Fprint(w, msg)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errHandler(w, "/add accepts POST requests")
	}
	r.ParseForm()

	// TODO these checks could be sent to a go routine
	name := r.Form.Get("name")
	if !(nameRegexp.Match([]byte(name))) {
		errHandler(w, "invalid name")
	}
	price := r.Form.Get("price")

	code := r.Form.Get("code")
	if !(codeRegexp.Match([]byte(code))) {
		errHandler(w, "invalid code")
	}

	// TODO go routine
	fmt.Println(name, price, code)
	fmt.Println(db)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Delete -- URL param `Produce Code`
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO -- GET all produce
}
