package main

import (
	"fmt"
	"os"
)

type Dirent struct {
	Path    string
	IsDir   bool
	Dirents []Dirent
}

// if the dirent is a directory, return the directory entries
// otherwise, return an empty slice
func (d *Dirent) LoadDirents() {
	if !d.IsDir {
		return
	}
	dir, err := os.Open(d.Path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	for _, fi := range fileInfos {
		d.Dirents = append(d.Dirents, Dirent{
			Path:  fi.Name(),
			IsDir: fi.IsDir(),
		})
	}
}


