package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

var wwwPathPtr, restPathPtr *string
var doLogPtr, doLogRestPtr *bool
var portPtr *int

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

func (menu *Menu) Save(path string) error {
	body, err := json.Marshal(menu)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, body, 0644)
	return err
}

func (menu *Menu) Load(path string) error {
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
	b, _ := json.Marshal(menu.Items) // Backbone wants the only the array
	if *doLogRestPtr {
		log.Printf("> %s\n\n", b)
	}
	w.Write(b)
}

func main() {
	portPtr = flag.Int("port", 8000, "server port (on localhost, default 8000")
	wwwPathPtr = flag.String("www", "_www", "path for serving web files")
	restPathPtr = flag.String("rest", "/api/items", "prefix for REST path")
	doLogPtr = flag.Bool("log", false, "log incoming requests")
	doLogRestPtr = flag.Bool("logrest", false, "log REST transactions (requests and responses)")
	flag.Parse()

	fmt.Println("Front-End Frameworks server")
	fmt.Println("---------------------------")
	fmt.Printf("  Serving on http://localhost:%v/\n", *portPtr)
	fmt.Printf("  Serving files from %v on /\n", *wwwPathPtr)
	fmt.Printf("  Serving REST requests on %v\n", *restPathPtr)
	if *doLogPtr || *doLogRestPtr {
		fmt.Println("  Logging incoming requests")
	}
	if *doLogRestPtr {
		fmt.Println("  Logging outgoing responses")
	}

	const filepath = "_data/menu.json"

	menu = &Menu{}
	menu.Load(filepath)
	http.Handle("/", http.FileServer(http.Dir(*wwwPathPtr)))
	http.HandleFunc(*restPathPtr, menuHandler)
	http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)
}
