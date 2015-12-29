package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

type Menu struct {
	Items    []MenuItem `json:"menu"`
	testpath string     // file path used only in unit testing
}

const filepath = "_data/menu.json"

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

var wwwPathPtr, restPathPtr *string
var doLogPtr, doLogRestPtr *bool

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
func (menu *Menu) reset() {
	menu.Items = make([]MenuItem, 0)
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

func (menu *Menu) Save() error {
	var path = filepath
	if menu.testpath != "" {
		path = menu.testpath
	}
	body, err := json.Marshal(menu)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, body, 0644)
	return err
}

func (menu *Menu) Load() error {
	var path = filepath
	if menu.testpath != "" {
		path = menu.testpath
	}
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &menu)
	return err
}

var menu *Menu

func logRequest(r *http.Request) {
	method := r.Method
	if method == "" {
		method = "GET"
	}
	log.Printf("< %s %s\n", method, r.RequestURI)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	if *doLogPtr || *doLogRestPtr {
		logRequest(r)
	}
	b, _ := json.Marshal(menu)
	if *doLogRestPtr {
		log.Printf("> %s\n\n", b)
	}
	w.Write(b)
}

func main() {
	wwwPathPtr = flag.String("www", "_www", "path for serving web files")
	restPathPtr = flag.String("rest", "/api/items", "prefix for REST path")
	doLogPtr = flag.Bool("log", false, "log incoming requests")
	doLogRestPtr = flag.Bool("logrest", false, "log REST transactions (requests and responses)")
	flag.Parse()

	menu = &Menu{}
	menu.Load()
	http.Handle("/", http.FileServer(http.Dir(*wwwPathPtr)))
	http.HandleFunc(*restPathPtr, menuHandler)
	http.ListenAndServe(":8080", nil)
}
