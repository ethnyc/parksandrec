/* Copyright (c) 2014-2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	listen = flag.String("l", ":8080", "Host and port to listen to")
)

type handler struct {
	tmpl   *template.Template
	places []Place
	users  []User
	activs []Activity
}

func newHttpHandler() *handler {
	h := &handler{
		users:  getUsers(),
		activs: getActivities(),
	}
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
	doMarshal(w, r, &places)
}

func filterActivs(activs []Activity, search string) []Activity {
	if search == "" {
		return activs
	}
	found := []Activity{}
	for _, p := range activs {
		if !p.Matches(search) {
			continue
		}
		found = append(found, p)
	}
	return found
}

func (h *handler) searchactivs(w http.ResponseWriter, r *http.Request) {
	search := mux.Vars(r)["search"]
	activs := filterActivs(h.activs, search)
	doMarshal(w, r, &activs)
}

func (h *handler) getplace(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	n, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n < 1 || n > len(h.places) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	doMarshal(w, r, &h.places[n-1])
}

func (h *handler) getuser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	n, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n < 1 || n > len(h.users) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	doMarshal(w, r, &h.users[n-1])
}

func (h *handler) getactivity(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	n, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n < 1 || n > len(h.activs) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	doMarshal(w, r, &h.activs[n-1])
}

func (h *handler) getusers(w http.ResponseWriter, r *http.Request) {
	doMarshal(w, r, &h.users)
}

func (h *handler) getavatar(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	doFile(w, r, filepath.Join("const", "avatars", id+".jpg"))
}

func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func doFile(w http.ResponseWriter, r *http.Request, path string) {
	addHeaders(w, r)
	http.ServeFile(w, r, path)
}

func doMarshal(w http.ResponseWriter, r *http.Request, v interface{}) {
	addHeaders(w, r)
	//json.NewEncoder(w).Encode(v)
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func main() {
	flag.Parse()
	h := newHttpHandler()
	r := mux.NewRouter()
	r.HandleFunc("/", h.index).Methods("GET")
	r.HandleFunc("/place/{id}", h.getplace).Methods("GET")
	r.HandleFunc("/places", h.searchplaces).Methods("GET")
	r.HandleFunc("/places/", h.searchplaces).Methods("GET")
	r.HandleFunc("/places/{search}", h.searchplaces).Methods("GET")
	r.HandleFunc("/user/{id}", h.getuser).Methods("GET")
	r.HandleFunc("/users", h.getusers).Methods("GET")
	r.HandleFunc("/users/", h.getusers).Methods("GET")
	r.HandleFunc("/avatar/{id}", h.getavatar).Methods("GET")
	r.HandleFunc("/activity/{id}", h.getactivity).Methods("GET")
	r.HandleFunc("/activities", h.searchactivs).Methods("GET")
	r.HandleFunc("/activities/", h.searchactivs).Methods("GET")
	r.HandleFunc("/activities/{search}", h.searchactivs).Methods("GET")
	log.Printf("listen = %s", *listen)
	http.Handle("/", r)
	log.Println("Up and running!")
	log.Fatal(http.ListenAndServe(*listen, nil))
}
