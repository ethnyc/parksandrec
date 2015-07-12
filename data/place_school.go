package main

import (
	"encoding/json"
	"log"
	"os"
)

func getSchools() []Place {
	f, err := os.Open("schools.json")
	if err != nil {
		log.Fatal(err)
	}
	var places []Place
	json.NewDecoder(f).Decode(&places)
	for i := range places {
		places[i].Type = "school"
	}
	return places
}
