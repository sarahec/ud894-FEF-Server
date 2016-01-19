package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	menu     *Menu
	filepath string
}

const contentType = "Content-Type"
const appJSON = "application/json"

func (s *Server) handleGetAll(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(s.menu.Items) // Backbone wants the only the array
	w.Header().Set(contentType, appJSON)
	w.Write(b)
}

func (s *Server) handleGetByID(w http.ResponseWriter, r *http.Request) {
	requestedID := r.URL.Path[1:]
	probe, ok := s.menu.GetByID(requestedID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		b, _ := json.Marshal(*probe) // return only the element
		w.Header().Set(contentType, appJSON)
		w.Write(b)
	}
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	requestedID := r.URL.Path[1:]

	var item MenuItem
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&item)

	if err != nil || item.ID == "" || item.ID != requestedID {
		log.Printf("ERROR: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appended := s.menu.Put(&item)

	// Tell the client if this was an update (default is 200 OK, which would be an update)
	if appended {
		w.WriteHeader(http.StatusNoContent) // Not returning the object
	}

	if filepath != "" {
		s.menu.Save(filepath)
	}
}

// Adopt the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		if r.URL.Path == "" {
			s.handleGetAll(w, r)
		} else {
			s.handleGetByID(w, r)
		}
	case r.Method == "PUT":
		s.handlePut(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
