package main

import (
	"fmt"
	"log"
	"os"
)

type Dirent struct {
	Path    string
	IsDir   bool
	Dirents []Dirent
	Level   uint // the level of the dirent in the tree; 0 = root, and so on
	Expanded    bool // whether to show the dirent in the UI
}

// if the dirent is a directory, return the directory entries
// otherwise, return an empty slice
func (d *Dirent) LoadDirents() {
	if !d.IsDir {
		return
	}
	dir, err := os.Open(d.Path)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	for _, fi := range fileInfos {
		d.Dirents = append(d.Dirents, Dirent{
			Path:  fi.Name(),
			IsDir: fi.IsDir(),
			Level: d.Level + 1,
      Expanded: false,
		})
	}
}

func (d Dirent) Print(cursor int) string {
	s := ""
	for i, dirent := range d.Dirents {
		// Render the dirent
    cursorDisplay := " "
    // If the dirent is the cursor, highlight it
    if i == cursor {
       cursorDisplay = ">"
    }
		pathname := dirent.Path
		if dirent.IsDir {
			pathname += "/"
		}
    if dirent.Expanded {
      if len(dirent.Dirents) == 0 {
        dirent.LoadDirents()
      }
      pathname += "\n" + dirent.Print(cursor)
    }
		prefixspace := ""
		for j := uint(1); j < dirent.Level; j++ {
			prefixspace += "  "
		}
		s += fmt.Sprintf("%s%s%s\n", cursorDisplay, prefixspace, pathname)
	}
	return s
}
