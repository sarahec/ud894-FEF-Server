package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const CONTENT_TYPE = "Content-Type"
const TYPE_JSON = "application/json"

// === Logging support

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

func logWrapper(doLog bool, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &LoggingResponseWriter{doLog, w}

		if doLog {
			method := r.Method
			if method == "" {
				method = "GET"
			}
			log.Printf("< %s %s\n", method, r.RequestURI)
		}
		handler(writer, r)
	})
}

// === HTTP handlers

func idFromURL(reqURL *url.URL, prefix string) string {
	return reqURL.Path[len(prefix)+1:]
}

func GetAllItemsServer(menu *Menu) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(menu.Items) // Backbone wants the only the array
		w.Header().Set(CONTENT_TYPE, TYPE_JSON)
		w.Write(b)
	})
}

func GetItemByIDServer(menu *Menu, prefix string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedID := idFromURL(r.URL, prefix)
		probe, ok := menu.GetByID(requestedID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			b, _ := json.Marshal(*probe) // return only the element
			w.Header().Set(CONTENT_TYPE, TYPE_JSON)
			w.Write(b)
		}
	})
}

func PutItemServer(menu *Menu, prefix string, filepath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedID := idFromURL(r.URL, prefix)

		var item MenuItem
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&item)

		if err != nil || item.ID == "" || item.ID != requestedID {
			log.Printf("ERROR: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		appended := menu.Put(&item)

		// Tell the client if this was an update (default is 200 OK, which would be an update)
		if appended {
			w.WriteHeader(http.StatusCreated)
		}

		if filepath != "" {
			menu.Save(filepath)
		}
	})
}

// === Handle switching between Get and Put

func NewRouter(get http.HandlerFunc, put http.HandlerFunc) http.HandlerFunc {
	if get == nil && put == nil {
		log.Fatal("Supply a get and/or put method to the router")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "" || r.Method == "GET" && get != nil:
			get(w, r)
		case r.Method == "PUT" && put != nil:
			put(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
