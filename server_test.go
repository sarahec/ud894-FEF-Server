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
	v := m.Get("nope")
	if v != nil {
		t.Errorf("Expected nil, but received %v", v)
	}
}

func TestGetMissingValue(t *testing.T) {
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one"})
	v := m.Get("nope")
	if v != nil {
		t.Errorf("Expected nil, but received %v", v)
	}
}

func TestGetSingleValue(t *testing.T) {
	m := &Menu{}
	item := &MenuItem{ID: "item-one"}
	m.Put(item)
	v := m.Get(item.ID)
	if v == nil {
		t.Errorf("Expected item, but received nil")
	}
}

func TestPutReplacementValue(t *testing.T) {
	m := &Menu{}
	item := &MenuItem{ID: "item-one", Name: "Item one"}
	m.Put(item)
	item2 := &MenuItem{ID: "item-one", Name: "Item one's replacement"}
	m.Put(item2) // same item
	if probe := m.Get(item.ID); probe.Name != item2.Name {
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
	m.Remove("item-one")
	if length := len(m.Items); length != 1 {
		t.Errorf("Expected one item after removal, but counted %v", length)
	}
	if probe := m.Get("item-one"); probe != nil {
		t.Errorf("Expected item one to be deleted, but found %v", probe)
	}
	if m.Get("item-two") == nil {
		t.Errorf("Expected item two to exist")
	}
}
