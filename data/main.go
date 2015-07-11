/* Copyright (c) 2014-2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var (
	listen = flag.String("l", ":8080", "Host and port to listen to")
)

type handler struct {
	tmpl  *template.Template
	clubs []Club
	get   endpoints
}

type endpoints map[string]func(http.ResponseWriter, *http.Request)

func newHttpHandler() *handler {
	h := &handler{
		tmpl:  template.Must(template.ParseFiles("index.html")),
		clubs: getClubs(),
	}
	h.get = endpoints{
		"/":              h.index,
		"/places":        h.places,
		"/places-pretty": h.placespretty,
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.ToLower(r.URL.Path)
	switch r.Method {
	case "GET":
		f, e := h.get[p]
		if !e {
			http.Error(w, "unknown endpoint", http.StatusBadRequest)
			return
		}
		f(w, r)
	default:
		http.Error(w, "unsupported action", http.StatusBadRequest)
	}
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "index.html", struct{
		Get endpoints
	}{
		Get: h.get,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) places(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	e.Encode(&h.clubs)
}

func (h *handler) placespretty(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(&h.clubs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
}

func main() {
	flag.Parse()
	handler := newHttpHandler()
	log.Printf("listen = %s", *listen)
	http.Handle("/", handler)
	log.Println("Up and running!")
	log.Fatal(http.ListenAndServe(*listen, nil))
}
