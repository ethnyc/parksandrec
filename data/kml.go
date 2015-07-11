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
		for _, l := range p.Desc.Lines {
			fmt.Println(l.Colour.Bold)
		}
	}
}

type CommunityClub struct {
	Name  string `json:"name"`
	Point string `json:"point"`
	Desc  string `json:"desc"`
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
		for _, l := range p.Desc.Lines {
			s := l.Colour.Bold
			if m := descMatch.FindStringSubmatch(s); m != nil {
				fmt.Println(m)
				c.Desc = m[1]
			}
		}
		clubs = append(clubs, c)
	}
	return clubs
}

func printCommunityClubs(clubs []CommunityClub) {
	for _, c := range clubs {
		fmt.Println(c.Name)
		fmt.Println(c.Desc)
		fmt.Println(c.Point)
	}
}

func main() {
	f, err := os.Open(communityClubsKML)
	if err != nil {
		log.Fatal(err)
	}
	k := kml{}
	xml.NewDecoder(f).Decode(&k)
	clubs := NewCommunityClubs(&k)
	printCommunityClubs(clubs)
}
