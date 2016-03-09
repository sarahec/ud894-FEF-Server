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

import "testing"

func TestInitialValue(t *testing.T) {
	m := &Menu{}
	if length := len(m.Items); length > 0 {
		t.Errorf("Expected empty menu, but counted %v elements", length)
	}
}

func TestPutSingleValue(t *testing.T) {
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one"})
	if length := len(m.Items); length != 1 {
		t.Errorf("Expected single item, but counted %v", length)
	}
}

func TestGetFromEmpty(t *testing.T) {
	m := &Menu{}
	_, ok := m.GetByID("nope")
	if ok {
		t.Error("Expected no item")
	}
}

func TestGetMissingValue(t *testing.T) {
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one"})
	_, ok := m.GetByID("nope")
	if ok {
		t.Error("Expected no item")
	}
}

func TestGetSingleValue(t *testing.T) {
	m := &Menu{}
	item := &MenuItem{ID: "item-one"}
	m.Put(item)
	v, ok := m.GetByID(item.ID)
	if v == nil || !ok {
		t.Errorf("Expected item, but received none")
	}
}

func TestPutReplacementValue(t *testing.T) {
	m := &Menu{}
	item := &MenuItem{ID: "item-one", Name: "Item one"}
	m.Put(item)
	item2 := &MenuItem{ID: "item-one", Name: "Item one's replacement"}
	m.Put(item2) // same item
	if probe, _ := m.GetByID(item.ID); probe.Name != item2.Name {
		t.Errorf("Expected replacement, but found name %s", probe.Name)
	}
}

func TestRemove(t *testing.T) {
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one", Name: "Item one"})
	m.Put(&MenuItem{ID: "item-two", Name: "Item two"})
	if length := len(m.Items); length != 2 {
		t.Errorf("Precondition error: Expected two items, but counted %v", length)
	}
	m.RemoveByID("item-one")
	if length := len(m.Items); length != 1 {
		t.Errorf("Expected one item after removal, but counted %v", length)
	}
	if _, ok := m.GetByID("item-one"); ok {
		t.Error("Expected item one to be deleted")
	}
	if _, ok := m.GetByID("item-two"); !ok {
		t.Error("Expected item two to exist")
	}
}

func TestReset(t *testing.T) {
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one", Name: "Item one"})
	m.reset()
	if length := len(m.Items); length != 0 {
		t.Errorf("Expected no items after reset, but counted %v", length)
	}
}
