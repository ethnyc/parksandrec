package main

import (
	"time"
)

type Activity struct {
	Unique
	Name  string    `json:"name"`
	Desc  string    `json:"desc"`
	Point string    `json:"point"`
	Loc   string    `json:"loc"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Users []int     `json:"users"`
	Cap   int       `json:"cap"`
}
