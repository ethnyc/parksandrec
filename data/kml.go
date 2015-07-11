package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
)

type kml struct {
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

const (
	clubsKML = "clubs.kml"
)

func (k *kml) debugPrint() {
	for _, p := range k.Document.Placemarks {
		fmt.Println(p.Name)
		fmt.Println(p.Point.Coord)
	}
}

type Typed struct {
	Type string `json:"type"`
}

type Club struct {
	Name  string `json:"name"`
	Point string `json:"point"`
	Typed
}

var (
	descMatch = regexp.MustCompile(`Description - (.*)`)
)

func NewClubs(k *kml) []Club {
	var clubs []Club
	for _, p := range k.Document.Placemarks {
		c := Club{
			Name:  p.Name,
			Point: p.Point.Coord,
		}
		c.Type = "club"
		clubs = append(clubs, c)
	}
	return clubs
}

func getClubs() []Club {
	f, err := os.Open(clubsKML)
	if err != nil {
		log.Fatal(err)
	}
	k := kml{}
	xml.NewDecoder(f).Decode(&k)
	return NewClubs(&k)
}
