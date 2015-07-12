package main

import (
	"regexp"
	"strings"
)

type Unique struct {
	Id int `json:"id"`
}

type Typed struct {
	Type string `json:"type"`
}

type Place struct {
	Unique
	Typed
	Name  string `json:"name"`
	Point string `json:"point"`
}

var typeMatch = regexp.MustCompile(`^type:(.*)$`)

func (p Place) Matches(s string) bool {
	s = strings.ToLower(s)
	n := strings.ToLower(p.Name)
	t := strings.ToLower(p.Type)
	for _, sub := range strings.Split(s, " ") {
		if m := typeMatch.FindStringSubmatch(sub); m != nil {
			if t != m[1] {
				return false
			}
		} else {
			if !strings.Contains(n, sub) {
				return false
			}
		}
	}
	return true
}
