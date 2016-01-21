package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildMenu() *Menu {
	item1 := MenuItem{ID: "spaghetti", Name: "Spaghetti and Meatballs"}
	items := []MenuItem{item1}
	return &Menu{Items: items}
}

func TestGetAll(t *testing.T) {
	menu := buildMenu()
	server := &Server{menu: menu}
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	server.handleGetAll(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET all didn't return %v, returned %v", http.StatusOK, w.Code)
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("GET all didn't return application/json, returned %v", contentType)
	}
	b, _ := json.Marshal(menu)
	if w.Body.String() != string(b) {
		t.Errorf("GET all didn't return menu items array as JSON, returned %v", w.Body.String())
	}
}

func TestGetByID(t *testing.T) {
	menu := buildMenu()
	server := &Server{menu: menu}

	req, _ := http.NewRequest("GET", "/spaghetti", nil)
	w := httptest.NewRecorder()
	server.handleGetByID(w, req)
	if w.Code != http.StatusOK {
		t.Errorf(".../spaghetti didn't return %v, returned %v", http.StatusOK, w.Code)
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(".../spaghetti didn't return application/json, returned %v", contentType)
	}
	b, _ := json.Marshal(menu.Items[0]) // only the item
	if w.Body.String() != string(b) {
		t.Errorf(".../spaghetti didn't return menu items array as JSON, returned %v", w.Body.String())
	}
}

func TestGetMissingByID(t *testing.T) {
	menu := buildMenu()
	server := &Server{menu: menu}
	req, _ := http.NewRequest("GET", "/NOT-THE-DROIDS-YOU-ARE-LOOKING-FOR", nil)
	w := httptest.NewRecorder()
	server.handleGetByID(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Missing item didn't return %v, returned %v", http.StatusNotFound, w.Code)
	}
}

func TestReplaceExisting(t *testing.T) {
	const path = "/spaghetti"
	menu := buildMenu()
	server := &Server{menu: menu}

	var update MenuItem = menu.Items[0]
	update.Name = "Spaghetti with Bolognaise sauce (new)"
	payload, _ := json.Marshal(update)

	if menu.Items[0].Name == update.Name {
		t.Error("Oops! We changed the menu, not a copy. Fix the test code.")
	}

	reader := bytes.NewReader(payload)
	req, _ := http.NewRequest("PUT", path, reader)
	w := httptest.NewRecorder()
	server.handlePut(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("update didn't return %v, returned %v", http.StatusOK, w.Code)
	}

	// Successful execution should have changed the menu
	newName := menu.Items[0].Name
	if newName != update.Name {
		t.Errorf("expected changed name %v, found %v", update.Name, newName)
	}
}

func TestAddNewItem(t *testing.T) {
	const path = "/tiramisu"
	menu := buildMenu()
	server := &Server{menu: menu}

	item := MenuItem{ID: "tiramisu", Name: "Delicious Tiramisu"}
	payload, _ := json.Marshal(item)

	reader := bytes.NewReader(payload)
	req, _ := http.NewRequest("PUT", path, reader)
	w := httptest.NewRecorder()
	server.handlePut(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("update didn't return %v, returned %v", http.StatusNoContent, w.Code)
	}

	// Successful execution should have changed the menu
	itemCount := len(menu.Items)
	if itemCount != 2 {
		t.Errorf("expected 2 items, found %v", itemCount)
	}

}
