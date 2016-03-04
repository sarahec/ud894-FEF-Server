package main

import (
	"flag"
	"os"
	"testing"
)

const testPath = "../_testData"

func TestMain(m *testing.M) {
	flag.Parse()
	//	os.RemoveAll(testPath)
	code := m.Run()
	os.RemoveAll(testPath)
	os.Exit(code)
}

func TestCreatingStorage(t *testing.T) {
	if _, err := os.Stat(testPath); err == nil {
		t.Errorf("Data dir exists before test")
	}
	dataPath, err := BuildStorageDir(true, testPath)
	if err != nil {
		t.Errorf("Error in BuildStorageDir: %v", err)
	}
	if _, err := os.Stat(testPath); err != nil {
		t.Errorf("Data dir not created: %v", err)
	}
	if _, err := os.Stat(dataPath); err != nil {
		t.Errorf("Backing store file not created: %v", err)
	}
}

func TestMenuLoading(t *testing.T) {
	filePath, err := BuildStorageDir(true, testPath)
	if err != nil {
		t.Errorf("Error in BuildStorageDir: %v", err)
	}
	m := &Menu{}
	m.Load(filePath)
	if _, ok := m.GetByID("strawberry-pudding"); !ok {
		t.Error("Menu not loaded")
	}
}

func TestMenuSaving(t *testing.T) {
	filePath, err := BuildStorageDir(true, testPath)
	if err != nil {
		t.Errorf("Error in BuildStorageDir: %v", err)
	}
	m := &Menu{}
	m.Put(&MenuItem{ID: "item-one"})
	err = m.Save(filePath)
	if err != nil {
		t.Errorf("Error saving menu: %v", err)
	}

	// Now load into a new menu
	m2 := &Menu{}
	err = m2.Load(filePath)
	if err != nil {
		t.Errorf("Error saving menu: %v", err)
	}
	if _, ok := m2.GetByID("item-one"); !ok {
		t.Error("Menu not loaded")
	}
}
