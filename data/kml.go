package main

import (
	"encoding/xml"
	"log"
	"os"
	"regexp"
)

type clubsKML struct {
	Document struct {
		Placemarks []placemark `xml:"Placemark"`
	}
}

type placemark struct {
	Name  string `xml:"name"`
	Desc  desc   `xml:"description"`
	Point point  `xml:"Point"`
}

type desc struct {
	Lines []struct {
		Colour struct {
			Bold string `xml:"b"`
		} `xml:"font"`
	} `xml:"p"`
}

type point struct {
	Coord string `xml:"coordinates"`
}

type Typed struct {
	Type string `json:"type"`
}

type Place struct {
	Name  string `json:"name"`
	Point string `json:"point"`
	Typed
}

var (
	ccNameMatch = regexp.MustCompile(`(.* CC)($| CC .*)`)
	descMatch   = regexp.MustCompile(`Description - (.*)`)
)

func (k *clubsKML) Places() []Place {
	var places []Place
	for _, p := range k.Document.Placemarks {
		n := p.Name
		if m := ccNameMatch.FindStringSubmatch(n); m != nil {
			n = m[1]
		}
		place := Place{
			Name:  n,
			Point: p.Point.Coord,
		}
		place.Type = "club"
		places = append(places, place)
	}
	return places
}

func getClubs() []Place {
	f, err := os.Open("clubs.kml")
	if err != nil {
		log.Fatal(err)
	}
	k := clubsKML{}
	xml.NewDecoder(f).Decode(&k)
	return k.Places()
}
