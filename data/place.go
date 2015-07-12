package main

import "strings"

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

func (p Place) Matches(s string) bool {
	s = strings.ToLower(s)
	return strings.Contains(strings.ToLower(p.Name), s)
}
