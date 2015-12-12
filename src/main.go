package main

import (
	"fmt"
	"os"
	"encoding/json"
)

type MenuItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Calories     uint    `json:"calories"`
	Rating       float32 `json:"rating"`
	Description  string  `json:"description"`
	Source       string  `json:"source"`
	Photographer string  `json:"photographer"`
}

func main() {
	salad := MenuItem{"chicken-pomegranate-salad", "Chicken Pomegranate Salad", "", 430, 4.1, "Fud. Is good", "", ""}
	b, err := json.Marshal(salad)
	if err != nil {
		fmt.Println(err)
	}
	os.Stdout.Write(b)
}
