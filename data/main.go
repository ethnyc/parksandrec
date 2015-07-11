/* Copyright (c) 2014-2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listen = flag.String("l", ":8080", "Host and port to listen to")
)

type httpHandler struct {
	Clubs []CommunityClub
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		http.Error(w, "unsupported action", http.StatusBadRequest)
	}
}

func (h *httpHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET")
}

func (h *httpHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST")
}

func main() {
	flag.Parse()
	handler := httpHandler{
		Clubs: getCommunityClubs(),
	}
	log.Printf("listen = %s", *listen)
	http.Handle("/", handler)
	log.Println("Up and running!")
	log.Fatal(http.ListenAndServe(*listen, nil))
}
