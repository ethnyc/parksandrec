package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type kml struct {
	Document struct {
		Placemarks []placemark `xml:"Placemark"`
	}
}

type placemark struct {
	Name  string `xml:"name"`
	Desc  string `xml:"description"`
	Point point  `xml:"Point"`
}

type point struct {
	Coord coord `xml:"coordinates"`
}

type coord string

const (
	communityClubsKML = "communityclubs.kml"
)

func (k *kml) debugPrint() {
	for _, p := range k.Document.Placemarks {
		fmt.Println(p.Name)
		fmt.Println(p.Point)
		fmt.Println(p.Desc)
	}
}

func main() {
	f, err := os.Open(communityClubsKML)
	if err != nil {
		log.Fatal(err)
	}
	k := kml{}
	xml.NewDecoder(f).Decode(&k)
	k.debugPrint()
}
