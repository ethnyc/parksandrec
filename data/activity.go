package main

import (
	"time"
)

type Activity struct {
	Unique
	Name  string
	Desc  string
	Point string
	Loc   string // desc of Point
	Start time.Time
	End   time.Time
	Cap   int
}
