package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

const filepath = "_data/menu.json"
const resetpath = "_data/starter_menu.json"

func main() {
	port := flag.Int("port", 8000, "server port (on localhost, default 8000")
	wwwPath := flag.String("www", "_www", "path for serving web files")
	restPath := flag.String("api", "/api/items", "prefix for REST path")
	doReset := flag.Bool("reset", false, "reset model from starter file before using")
	doLog := flag.Bool("log", false, "log REST requests and responses")
	doSummarize := flag.Bool("verbose", false, "show summary of settings")
	flag.Parse()

	if *doSummarize {
		if *doReset {
			fmt.Printf("  Re-loading model from %v\n", resetpath)
		}
		fmt.Printf("  Serving on http://localhost:%v/\n", *port)
		fmt.Printf("  Serving files from %v on /\n", *wwwPath)
		fmt.Printf("  Serving REST requests on %v\n", *restPath)
		if *doLog {
			fmt.Printf("  Logging requests and responses for %v\n", *restPath)
		}
	}

	menu := &Menu{}
	if *doReset {
		menu.Load(resetpath)
	} else {
		menu.Load(filepath)
	}
	http.Handle("/", http.FileServer(http.Dir(*wwwPath)))
	http.Handle(*restPath, logWrapper(*doLog, NewRouter(GetAllItemsServer(menu), nil)))
	http.Handle(*restPath+"/", logWrapper(*doLog, NewRouter(
		GetItemByIDServer(menu, *restPath),
		PutItemServer(menu, *restPath, filepath))))
	// TODO catch and report errors
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
