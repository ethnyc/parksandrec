package main

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
