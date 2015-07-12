/* Copyright (c) 2014-2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	listen = flag.String("l", ":8080", "Host and port to listen to")
)

type handler struct {
	tmpl   *template.Template
	places []Place
}

func newHttpHandler() *handler {
	h := &handler{}
	h.addPlaces(getClubs())
	h.addPlaces(getSchools())
	return h
}

func (h *handler) addPlaces(places []Place) {
	for _, p := range places {
		p.Id = len(h.places) + 1
		h.places = append(h.places, p)
	}
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "README.md")
}

func filterPlaces(places []Place, search string) []Place {
	if search == "" {
		return places
	}
	found := []Place{}
	for _, p := range places {
		if !p.Matches(search) {
			continue
		}
		found = append(found, p)
	}
	return found
}

func (h *handler) searchplaces(w http.ResponseWriter, r *http.Request) {
	search := mux.Vars(r)["search"]
	places := filterPlaces(h.places, search)
	marshal(w, &places)
}

func marshal(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//json.NewEncoder(w).Encode(v)
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func main() {
	flag.Parse()
	h := newHttpHandler()
	r := mux.NewRouter()
	r.HandleFunc("/", h.index).Methods("GET")
	r.HandleFunc("/places", h.searchplaces).Methods("GET")
	r.HandleFunc("/places/", h.searchplaces).Methods("GET")
	r.HandleFunc("/places/{search}", h.searchplaces).Methods("GET")
	log.Printf("listen = %s", *listen)
	http.Handle("/", r)
	log.Println("Up and running!")
	log.Fatal(http.ListenAndServe(*listen, nil))
}
