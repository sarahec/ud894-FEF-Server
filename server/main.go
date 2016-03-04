package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"log"
)

const noLogging = 0

func addLogging(logLevel int, handler http.Handler) http.Handler {
	if logLevel == noLogging {
		return handler
	}
	return logWrapper(logLevel, handler)
}

func main() {
	const dataPath = "../_data"

	port := flag.Int("port", 8000, "server port (on localhost, default 8000")
	wwwPath := flag.String("www", "../_www", "path for serving web files")
	restPath := flag.String("api", "/api/items", "prefix for REST path")
	reset := flag.Bool("reset", false, "reset model from starter file before using")
	logging := flag.Bool("log", false, "log REST requests (method and URL) and JSON responses")
	debugRequests := flag.Bool("debug", false, "show details of REST requests and responses (more verbose than -log)")
	summarize := flag.Bool("verbose", false, "show summary of settings")
	flag.Parse()

	if *summarize {
		if *reset {
			fmt.Println("  Re-loading model")
		}
		fmt.Printf("  Serving on http://localhost:%v/\n", *port)
		fmt.Printf("  Serving files from %v on /\n", *wwwPath)
		fmt.Printf("  Serving REST requests on %v\n", *restPath)
		if *logging {
			fmt.Printf("  Logging requests and responses for %v\n", *restPath)
		}
	}

	filePath, err := BuildStorageDir(*reset, dataPath)
	if err != nil {
		log.Fatal("Error creating backing store: %v", err)
	}

	menu := &Menu{}
	menu.Load(filePath)

	logLevel := noLogging
	switch {
	case *logging:
		logLevel = LogMethodAndURL
	case *debugRequests:
		logLevel = LogWholeRequest
	}

	http.Handle("/", http.FileServer(http.Dir(*wwwPath)))
	server := &Server{menu, filePath}
	handler := addLogging(logLevel, http.StripPrefix(*restPath, server))
	http.Handle(*restPath, handler)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
