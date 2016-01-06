package main

import (
	"bytes"
	"encoding/json"
	"io"
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

func TestReplaceExisting(t *testing.T) {
	const path = PREFIX + "/spaghetti"
	menu := buildMenu()

	var update MenuItem = menu.Items[0]
	update.Name = "Spaghetti with Bolognaise sauce (new)"
	payload, _ := json.Marshal(update)

	if menu.Items[0].Name == update.Name {
		t.Error("Oops! We changed the menu, not a copy. Fix the test code.")
	}

	reader := bytes.NewReader(payload)
	req, _ := http.NewRequest("PUT", path, reader)
	w := httptest.NewRecorder()
	h := PutItemServer(menu, PREFIX)
	h.ServeHTTP(w, req)

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
	const path = PREFIX + "/tiramisu"
	menu := buildMenu()

	item := MenuItem{ID: "tiramisu", Name: "Delicious Tiramisu"}
	payload, _ := json.Marshal(item)

	reader := bytes.NewReader(payload)
	req, _ := http.NewRequest("PUT", path, reader)
	w := httptest.NewRecorder()
	h := PutItemServer(menu, PREFIX)
	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("update didn't return %v, returned %v", http.StatusOK, w.Code)
	}

	// Successful execution should have changed the menu
	itemCount := len(menu.Items)
	if itemCount != 2 {
		t.Errorf("expected 2 items, found %v", itemCount)
	}

}

// === Test the switch that manages get and put

func returnMethod() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "" {
			io.WriteString(w, "GET")
		} else {
			io.WriteString(w, r.Method)
		}
	})
}

func TestGetRoutesCorrectly(t *testing.T) {
	router := NewRouter(returnMethod(), nil)

	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected OK (%v) but saw %v", http.StatusOK, w.Code)
	}

	if w.Body.String() != "GET" {
		t.Errorf("Expected GET but saw %v", w.Body.String())
	}
}

func TestMissingGetRoutesCorrectly(t *testing.T) {
	router := NewRouter(nil, returnMethod())

	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected Method Not Allowed (%v) but saw %v", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestPutRoutesCorrectly(t *testing.T) {
	router := NewRouter(nil, returnMethod())

	req, _ := http.NewRequest("PUT", "", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected OK (%v) but saw %v", http.StatusOK, w.Code)
	}

	if w.Body.String() != "PUT" {
		t.Errorf("Expected PUT but saw %v", w.Body.String())
	}
}

func TestMissingPutRoutesCorrectly(t *testing.T) {
	router := NewRouter(returnMethod(), nil)

	req, _ := http.NewRequest("PUT", "", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected Method Not Allowed (%v) but saw %v", http.StatusMethodNotAllowed, w.Code)
	}
}
