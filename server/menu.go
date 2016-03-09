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

type Menu struct {
	Items []MenuItem `json:"menu"`
	Path  string     `json:"-"`
}

type MenuItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Calories     uint    `json:"calories"`
	Rating       float32 `json:"rating"`
	Description  string  `json:"description"`
	Source       string  `json:"source"`
	Photographer string  `json:"photographer"`
}

// Reset the mnu to its empty state. Used only in init testing
func (menu *Menu) reset() {
	menu.Items = make([]MenuItem, 0)
}

// Finds the menu item with a given ID and returns its index
// and true if found
func (menu *Menu) findIndex(id string) (i int, found bool) {
	for i, v := range menu.Items {
		if v.ID == id {
			return i, true
		}
	}
	return 0, false
}

// Put this menu item into the collection, overwriting the one with
// the same ID if it exists. Returns a bool indicating whether this
// item was appended

func (menu *Menu) Put(item *MenuItem) (appended bool) {
	index, found := menu.findIndex(item.ID)
	if !found {
		menu.Items = append(menu.Items, *item)
	} else {
		menu.Items[index] = *item
	}
	return !found
}

// Get a menu item from the collection, matching by ID.
// Return the found item and true if found (false if not)
func (menu *Menu) GetByID(id string) (*MenuItem, bool) {
	index, found := menu.findIndex(id)
	if !found {
		return nil, false
	}
	return &menu.Items[index], true
}

// Remove the menu item with the specified ID from the collection
func (menu *Menu) RemoveByID(id string) (ok bool) {
	i, ok := menu.findIndex(id)
	if ok {
		// delete the ith item
		menu.Items = menu.Items[:i+copy(menu.Items[i:], menu.Items[i+1:])]
	}
	return
}
