package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	OS_READ = 04
	OS_WRITE = 02
	OS_EX = 01
	OS_USER_SHIFT = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT = 0

	OS_USER_R = OS_READ<<OS_USER_SHIFT
	OS_USER_W = OS_WRITE<<OS_USER_SHIFT
	OS_USER_X = OS_EX<<OS_USER_SHIFT
	OS_USER_RW = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R = OS_READ<<OS_GROUP_SHIFT
	OS_GROUP_W = OS_WRITE<<OS_GROUP_SHIFT
	OS_GROUP_X = OS_EX<<OS_GROUP_SHIFT
	OS_GROUP_RW = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R = OS_READ<<OS_OTH_SHIFT
	OS_OTH_W = OS_WRITE<<OS_OTH_SHIFT
	OS_OTH_X = OS_EX<<OS_OTH_SHIFT
	OS_OTH_RW = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
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
