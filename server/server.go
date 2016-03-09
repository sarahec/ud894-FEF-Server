/*
Copyright 2016 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	menu *Menu
}

const contentType = "Content-Type"
const appJSON = "application/json"

func (s *Server) handleGetAll(w http.ResponseWriter, r *http.Request) {
	// Backbone wants the only the array (by default), but sending top-level
	// JSON arrays is insecure. Sent the whole object and patch the Backbone
	// app instead.
	b, _ := json.Marshal(s.menu)
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

	s.menu.Save()
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
