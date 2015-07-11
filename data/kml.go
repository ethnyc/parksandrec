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
	communityClubsKML = "communityclubs.kml"
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

type CommunityClub struct {
	Name  string `json:"name"`
	Point string `json:"point"`
	Typed
}

var (
	descMatch = regexp.MustCompile(`Description - (.*)`)
)

func NewCommunityClubs(k *kml) []CommunityClub {
	var clubs []CommunityClub
	for _, p := range k.Document.Placemarks {
		c := CommunityClub{
			Name:  p.Name,
			Point: p.Point.Coord,
		}
		c.Type = "communityclub"
		clubs = append(clubs, c)
	}
	return clubs
}

func getCommunityClubs() []CommunityClub {
	f, err := os.Open(communityClubsKML)
	if err != nil {
		log.Fatal(err)
	}
	k := kml{}
	xml.NewDecoder(f).Decode(&k)
	return NewCommunityClubs(&k)
}
