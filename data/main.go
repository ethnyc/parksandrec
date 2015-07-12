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
	"sync"

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
	m      sync.RWMutex
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
	h.m.RLock()
	defer h.m.RUnlock()
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

func (h *handler) postactivity(w http.ResponseWriter, r *http.Request) {
	h.m.Lock()
	defer h.m.Unlock()
	decoder := json.NewDecoder(r.Body)
	var a Activity
	if err := decoder.Decode(&a); err != nil {
		log.Print(err)
	}
	a.Id = len(h.activs) + 1
	a.Parts = []int{}
	if a.Cap < 1 {
		a.Cap = 1
	}
	if a.Cap > 50 {
		a.Cap = 50
	}
	h.activs = append(h.activs, a)
}

func (h *handler) getactivity(w http.ResponseWriter, r *http.Request) {
	h.m.RLock()
	defer h.m.RUnlock()
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

func (h *handler) getidimg(kind, name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		doFile(w, r, filepath.Join(kind, "img", name, id+".jpg"))
	}
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
	r.HandleFunc("/activity", h.postactivity).Methods("POST")
	r.HandleFunc("/activity/{id}", h.getactivity).Methods("GET")
	r.HandleFunc("/activities", h.searchactivs).Methods("GET")
	r.HandleFunc("/activities/", h.searchactivs).Methods("GET")
	r.HandleFunc("/activities/{search}", h.searchactivs).Methods("GET")
	r.HandleFunc("/img/place/{id}", h.getidimg("const", "place")).Methods("GET")
	r.HandleFunc("/img/user/{id}", h.getidimg("const", "user")).Methods("GET")
	r.HandleFunc("/img/activity/{id}", h.getidimg("var", "activity")).Methods("GET")
	log.Printf("listen = %s", *listen)
	http.Handle("/", r)
	log.Println("Up and running!")
	log.Fatal(http.ListenAndServe(*listen, nil))
}
