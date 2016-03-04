package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/GeertJohan/go.rice"
)

const fileName = "menu.json"

//go:generate rice embed-go
var (
	box = rice.MustFindBox("assets")
)

func BuildStorageDir(resetContents bool, path string) (string, error) {
	// Make the directory (if needed)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			return "", err
		}
	}

	// Make the data file (if needed)
	filePath := path + "/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) || resetContents {
		data := box.MustBytes(fileName) // panics on error
		err := ioutil.WriteFile(filePath, data, os.FileMode(0755))
		if err != nil {
			return filePath, err
		}
	}

	return filePath, nil
}

func (menu *Menu) Save(path string) error {
	body, err := json.Marshal(menu)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, body, 0644)
}

func (menu *Menu) Load(path string) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &menu)
	return err
}
