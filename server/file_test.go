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
	"os"
	"testing"
)

func getTestpath() string {
	return os.TempDir() + "/testmodel.json"
}

func remove(filepath string) {
	os.RemoveAll(filepath)
}

func TestModelSavesToFile(t *testing.T) {
	filepath := getTestpath()
	m := &Menu{Path: filepath}
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
	filepath := getTestpath()
	defer remove(filepath)
	m := &Menu{Path: filepath}
	m.Put(&MenuItem{ID: "item-one", Name: "Item one"})

	err := m.Save()
	if err != nil {
		t.Errorf("Error in save: %v", err)
	}

	m.reset()
	if len(m.Items) > 0 {
		t.Errorf("Precondition failed: reset did not work")
	}

	err = m.Load()
	if err != nil {
		t.Errorf("Error in load: %v", err)
	}
	if len := len(m.Items); len != 1 {
		t.Errorf("Expected one item after load, but counted %v", len)
	}
	if _, ok := m.GetByID("item-one"); !ok {
		t.Errorf("Expected to find item-one")
	}
}
