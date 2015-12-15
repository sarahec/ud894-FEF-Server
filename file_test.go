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
	_, err := os.Stat(filepath)
	if err != nil {
		return
	}
	os.Remove(filepath)
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
	if err != nil {
		t.Errorf("Error in file stat: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Errorf("Nothing saved in file")
	}
}
