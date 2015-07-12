package main

import (
	"encoding/json"
	"log"
	"os"
)

type User struct {
	Unique
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func getUsers() []User {
	f, err := os.Open("const/users.json")
	if err != nil {
		log.Fatal(err)
	}
	var users []User
	json.NewDecoder(f).Decode(&users)
	return users
}
