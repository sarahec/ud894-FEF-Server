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
	"flag"
	"os"
	"testing"

	"github.com/GeertJohan/go.rice"
)

const testPath = "../_testData"

//go:generate rice embed-go
var box = rice.MustFindBox("assets")

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
	dataPath, err := BuildStorageDir(true, testPath, *box)
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
	filePath, err := BuildStorageDir(true, testPath, *box)
	if err != nil {
		t.Errorf("Error in BuildStorageDir: %v", err)
	}
	m := &Menu{Path: filePath}
	m.Load()
	if _, ok := m.GetByID("strawberry-pudding"); !ok {
		t.Error("Menu not loaded")
	}
}

func TestMenuSaving(t *testing.T) {
	filePath, err := BuildStorageDir(true, testPath, *box)
	if err != nil {
		t.Errorf("Error in BuildStorageDir: %v", err)
	}
	m := &Menu{Path: filePath}
	m.Put(&MenuItem{ID: "item-one"})
	err = m.Save()
	if err != nil {
		t.Errorf("Error saving menu: %v", err)
	}

	// Now load into a new menu
	m2 := &Menu{Path: filePath}
	err = m2.Load()
	if err != nil {
		t.Errorf("Error saving menu: %v", err)
	}
	if _, ok := m2.GetByID("item-one"); !ok {
		t.Error("Menu not loaded")
	}
}
