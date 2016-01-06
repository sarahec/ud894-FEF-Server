package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const PREFIX = "/api/items"

func buildMenu() *Menu {
	item1 := MenuItem{ID: "spaghetti", Name: "Spaghetti and Meatballs"}
	items := []MenuItem{item1}
	return &Menu{Items: items}
}

func TestGetAll(t *testing.T) {
	menu := buildMenu()
	h := GetAllItemsServer(menu)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("API root didn't return %v, returned %v", http.StatusOK, w.Code)
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("API root didn't return application/json, returned %v", contentType)
	}
	b, _ := json.Marshal(menu.Items) // Backbone wants the only the array
	if w.Body.String() != string(b) {
		t.Errorf("API root didn't return menu items array as JSON, returned %v", w.Body.String())
	}
}

func TestGetByID(t *testing.T) {
	menu := buildMenu()
	h := GetItemByIDServer(menu, PREFIX)
	req, _ := http.NewRequest("GET", PREFIX+"/spaghetti", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
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
	h := GetItemByIDServer(menu, PREFIX)
	req, _ := http.NewRequest("GET", PREFIX+"/NOT-THE-DROIDS-YOU-ARE-LOOKING-FOR", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Missing item didn't return %v, returned %v", http.StatusNotFound, w.Code)
	}
}
