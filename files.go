package main

import (
	"io/fs"
	"log"
	"os"
)

// these are files like .git/ that are not included in the .gitignore
// but should be ignored
func isObliviousPattern(fi fs.FileInfo) bool {
	if fi.IsDir() && fi.Name() == ".git" {
		return true
	}
	return false
}

// retrieve file info from the path
func getFi(path string) fs.FileInfo {
	// get the directory entries for the current directory
	dir, err := os.Open(path)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}

	defer dir.Close()

	fi, err := dir.Stat()
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	return fi
}
