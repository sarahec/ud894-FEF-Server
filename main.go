package main

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
)

type Menu struct {
	Items []MenuItem `json:"menu"`
}

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

var menu Menu

func loadMenu() Menu {
	  m := Menu{}
    body, err := ioutil.ReadFile("_data/menu.json")
		if err == nil {
			json.Unmarshal(body, &m)
		}
    return m
}

func handler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(menu)
	w.Write(b)
}

func main() {
	  menu = loadMenu()
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
