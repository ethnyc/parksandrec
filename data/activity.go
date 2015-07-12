package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Activity struct {
	Unique
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Point string `json:"point"`
	Loc   string `json:"loc"`
	Start string `json:"start"`
	End   string `json:"end"`
	Parts []int  `json:"participants"`
	Cap   int    `json:"cap"`
}

func getActivities() []Activity {
	f, err := os.Open("var/activities.json")
	if err != nil {
		log.Fatal(err)
	}
	var activities []Activity
	json.NewDecoder(f).Decode(&activities)
	return activities
}

func (a Activity) Matches(s string) bool {
	s = strings.ToLower(s)
	fields := []string{
		strings.ToLower(a.Name),
		strings.ToLower(a.Desc),
		strings.ToLower(a.Loc),
		strings.ToLower(a.Start),
		strings.ToLower(a.End),
	}
	for _, sub := range strings.Split(s, " ") {
		if !anyContains(fields, sub) {
			return false
		}
	}
	return true
}
