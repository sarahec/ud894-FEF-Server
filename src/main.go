package main

import (
	"encoding/json"
	"net/http"
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

func handler(w http.ResponseWriter, r *http.Request) {
	salad := MenuItem{"chicken-pomegranate-salad", "Chicken Pomegranate Salad", "", 430, 4.1, "Fud. Is good", "", ""}
	b, _ := json.Marshal(salad)
	w.Write(b)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
