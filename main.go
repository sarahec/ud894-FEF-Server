package main

import (
	"encoding/json"
	"io/ioutil"
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

// Searches for the specified id string in the menu, returning its index
// or -1 if not found
func (menu *Menu) indexOf(id string) int {
	for i, v := range menu.Items {
		if v.ID == id {
			return i
		}
	}
	return -1
}

// Put this menu item into the collection, overwriting the one with the same ID if it exists.
func (menu *Menu) Put(item *MenuItem) {
	probe := menu.indexOf(item.ID)
	if probe == -1 {
		menu.Items = append(menu.Items, *item)
	} else {
		menu.Items[probe] = *item
	}
}

// Get a menu item from the collection, matching in ID. Return nil if not found.
func (menu *Menu) Get(id string) *MenuItem {
	probe := menu.indexOf(id)
	if probe == -1 {
		return nil
	}
	found := menu.Items[probe]
	return &found
}

// Remove the menu item with the specified ID from the collection
func (menu *Menu) Remove(id string) {
	i := menu.indexOf(id)
	if i == -1 {
		return
	}
	// delete the ith item
	menu.Items = menu.Items[:i+copy(menu.Items[i:], menu.Items[i+1:])]
}


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
