/* Copyright (c) 2014-2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	switch strings.ToLower(r.URL.Path) {
	case "/":
		fmt.Fprintf(w, "Endpoints:\n\n")
		fmt.Fprintf(w, "  GET /places\n")
		fmt.Fprintf(w, "  GET /places-pretty\n")
	case "/places":
		e := json.NewEncoder(w)
		e.Encode(&h.Clubs)
	case "/places-pretty":
		b, err := json.MarshalIndent(&h.Clubs, "", "  ")
		if err != nil {
			log.Println(err)
		}
		w.Write(b)
	default:
		http.Error(w, "unknown endpoint", http.StatusBadRequest)
	}
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
