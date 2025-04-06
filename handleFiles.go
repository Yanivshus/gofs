package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var base_dir = "/home"

type dir struct {
	Path string `json:"path"`
}

func getfiles() ([]string, error) {
	root := os.DirFS(".")

	files, err := fs.Glob(root, "*")

	if err != nil {
		return nil, err
	}

	return files, nil
}

func change_dir(dir string) error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// can only go back with ..
	// can only go forward by specifing the dir to move to.
	targetDir := filepath.Join(currentPath, dir)
	cleanDir := filepath.Clean(targetDir)

	// check if there is a perfix (home) for now
	if !strings.HasPrefix(cleanDir, base_dir) {
		err = errors.New("trying to go to forbidden directory")
		return err
	}

	fmt.Println(cleanDir)
	//check if directory exists
	_, err = os.Stat(cleanDir)
	if err != nil {
		return err
	}

	err = os.Chdir(cleanDir) // change to the directory
	if err != nil {
		return err
	}

	return nil

}

func get_wd() (string, error) {
	data, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return data, nil
}
