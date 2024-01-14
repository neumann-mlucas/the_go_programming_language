package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	//!+main
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/get", db.get)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)
	//!-main
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]int

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%q: $%d\n", item, price)
	}
}

func (db database) get(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "$ get %q : %d\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "item not found: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	p := req.URL.Query().Get("price")
	price, err := strconv.Atoi(p)
	if _, ok := db[item]; !ok && err == nil {
		db[item] = price
		fmt.Fprintf(w, "created %q : $%d\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "item already exists: %q\n", item)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	p := req.URL.Query().Get("price")
	price, err := strconv.Atoi(p)
	if _, ok := db[item]; ok && err == nil {
		db[item] = price
		fmt.Fprintf(w, "updated %s : $%d\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "item not found: %q\n", item)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	fmt.Println(item)
	if _, ok := db[item]; ok {
		delete(db, item)
		fmt.Fprintf(w, "deleted %q\n", item)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "item not found: %q\n", item)
	}
}
