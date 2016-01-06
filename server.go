package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type LoggingResponseWriter struct {
	logResponse bool
	http.ResponseWriter
}

func (w *LoggingResponseWriter) Write(data []byte) (int, error) {
	if w.logResponse {
		log.Printf("> %s\n\n", data)
	}
	return w.ResponseWriter.Write(data)
}

func logWrapper(handler http.HandlerFunc) http.Handler {
	doLogRequest := *doLogPtr || *doLogRestPtr
	doLogResponse := *doLogRestPtr

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &LoggingResponseWriter{doLogResponse, w}

		if doLogRequest {
			method := r.Method
			if method == "" {
				method = "GET"
			}
			log.Printf("< %s %s\n", method, r.RequestURI)
		}

		handler(writer, r)
	})
}

func idFromURL(reqURL *url.URL, prefix string) string {
	return reqURL.Path[len(prefix)+1:]
}

func findItemByID(requestedID string, menu *Menu) *MenuItem {
	index := -0
	for index < len(menu.Items) {
		if menu.Items[index].ID == requestedID {
			return &menu.Items[index]
		}
		index += 1
	}
	return nil
}

func GetAllItemsServer(menu *Menu) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(menu.Items) // Backbone wants the only the array
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
}

func GetItemByIDServer(menu *Menu, prefix string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedID := idFromURL(r.URL, prefix)
		probe := findItemByID(requestedID, menu)
		if probe == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			b, _ := json.Marshal(*probe) // Backbone wants the only the array
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		}
	})
}
