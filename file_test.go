package main

import (
		"testing"
		"os"
	)

func setTestpath(m *Menu) (filepath string) {
	filepath = os.TempDir() + "testmodel.json"
	m.testpath = filepath
	return
}

func remove(filepath string) {
	os.RemoveAll(filepath)
}

func TestModelSavesToFile(t *testing.T) {
	m := &Menu{}
	filepath := setTestpath(m)
	defer remove(filepath)

	err := m.Save()
	if err != nil {
		t.Errorf("Error in save: %v", err)
	}
	fileInfo, err := os.Stat(filepath)
	if err != nil || fileInfo == nil {
		t.Errorf("Error in file stat: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Errorf("Nothing saved in file")
	}
}

func TestModelReadsFile(t *testing.T) {
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one", Name: "Item one"})

	filepath := setTestpath(m)
	defer remove(filepath)

	err := m.Save()
	if err != nil {
		t.Errorf("Error in save: %v", err)
	}

	m.reset();
	if len(m.Items) > 0 {
		t.Errorf("Precondition failed: reset did not work")
	}

	err = m.Load();
	if (err != nil) {
		t.Errorf("Error in load: %v", err)
	}
	if len := len(m.Items); len != 1 {
		t.Errorf("Expected one item after load, but counted %v", len)
	}
	if m.Get("item-one") == nil {
		t.Errorf("Expected to find item-one")
	}
}
