package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var wwwPathPtr, restPathPtr *string
var doLogPtr, doLogRestPtr *bool
var portPtr *int

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
