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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go.rice"
)

const fileName = "menu.json"

func BuildStorageDir(resetContents bool, path string, box rice.Box) (string, error) {
	// Make the directory (if needed)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			return "", err
		}
	}

	// Make the data file (if needed)
	filePath := filepath.Join(path, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) || resetContents {
		data := box.MustBytes(fileName) // panics on error
		err := ioutil.WriteFile(filePath, data, os.FileMode(0755))
		if err != nil {
			return filePath, err
		}
	}

	return filePath, nil
}

func (menu *Menu) Save() error {
	if menu.Path == "" {
		return nil // nowhere to save
	}
	body, err := json.Marshal(menu)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(menu.Path, body, 0644)
}

func (menu *Menu) Load() error {
	if menu.Path == "" {
		return os.ErrNotExist // nowhere to save
	}
	body, err := ioutil.ReadFile(menu.Path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &menu)
	return err
}
