package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

var wwwPathPtr, restPathPtr *string
var doLogPtr, doResetPtr, doSummarizePtr *bool
var portPtr *int

const filepath = "_data/menu.json"
const resetpath = "_data/starter_menu.json"

func main() {
	portPtr = flag.Int("port", 8000, "server port (on localhost, default 8000")
	wwwPathPtr = flag.String("www", "_www", "path for serving web files")
	restPathPtr = flag.String("api", "/api/items", "prefix for REST path")
	doResetPtr = flag.Bool("reset", false, "reset model from starter file before using")
	doLogPtr = flag.Bool("log", false, "log REST requests and responses")
	doSummarizePtr = flag.Bool("verbose", false, "show summary of settings")
	flag.Parse()

	if *doSummarizePtr {
		if *doResetPtr {
			fmt.Printf("  Re-loading model from %v\n", resetpath)
		}
		fmt.Printf("  Serving on http://localhost:%v/\n", *portPtr)
		fmt.Printf("  Serving files from %v on /\n", *wwwPathPtr)
		fmt.Printf("  Serving REST requests on %v\n", *restPathPtr)
		if *doLogPtr {
			fmt.Printf("  Logging requests and responses for %v\n", *restPathPtr)
		}
	}

	menu := &Menu{}
	if *doResetPtr {
		menu.Load(resetpath)
	} else {
		menu.Load(filepath)
	}
	http.Handle("/", http.FileServer(http.Dir(*wwwPathPtr)))
	http.Handle(*restPathPtr, logWrapper(NewRouter(GetAllItemsServer(menu), nil)))
	http.Handle(*restPathPtr+"/", logWrapper(NewRouter(
		GetItemByIDServer(menu, *restPathPtr),
		PutItemServer(menu, *restPathPtr, filepath))))
	// TODO catch and report errors
	http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)
}
