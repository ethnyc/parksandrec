package main

type User struct {
	Unique
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
