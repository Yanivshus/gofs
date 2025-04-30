package files

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

var Proj_dir string

type Dir struct {
	Path string `json:"path"`
}

type File_t struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isdir"`
	Type  uint32 `json:"fm"`
	Perm  uint32 `json:"perm"`
}

func Getfiles() (string, error) {
	var dir []File_t

	c, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}
	for _, entry := range c {
		dir = append(dir, File_t{
			entry.Name(),
			entry.IsDir(),
			uint32(entry.Type().Type()),
			uint32(entry.Type().Perm()),
		})
	}

	json_bbl, _ := json.Marshal(dir) // get the list as json
	return string(json_bbl), nil
}

func Change_Dir(Dir string) error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// can only go back with ..
	// can only go forward by specifing the Dir to move to.
	targetDir := filepath.Join(currentPath, Dir)
	cleanDir := filepath.Clean(targetDir)

	// check if there is a perfix (home) for now
	if !strings.HasPrefix(cleanDir, Proj_dir) {
		err = errors.New("trying to go to forbidden Directory")
		return err
	}

	//check if Directory exists
	_, err = os.Stat(cleanDir)
	if err != nil {
		return err
	}

	err = os.Chdir(cleanDir) // change to the Directory
	if err != nil {
		return err
	}

	return nil

}

func Get_wd() (string, error) {
	data, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return data, nil
}
