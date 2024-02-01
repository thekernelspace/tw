package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Dirent struct {
	Path        string
	IsDir       bool
	Dirents     []Dirent
	Parent      *Dirent // the parent directory; nil if root
	PosInParent int     // the index in the parent's dirent list
	Level       uint    // the level of the dirent in the tree; 0 = root, and so on
	Expanded    bool    // whether to show the dirent in the UI
}

func (d Dirent) Equals(other Dirent) bool {
	return d.Path == other.Path
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
	for pos, fi := range fileInfos {
		path := filepath.Join(d.Path, fi.Name())
		d.Dirents = append(d.Dirents, Dirent{
			Path:        path,
			IsDir:       fi.IsDir(),
			Level:       d.Level + 1,
			Parent:      d,
			PosInParent: pos,
			Expanded:    false, // don't show by default
		})
	}
	log.Printf("Loaded %d dirents for %s\n", len(d.Dirents), d.Path)
}

func (d Dirent) Print(state Model) string {
	s := ""
	for i := range d.Dirents {
		// Render the dirent
		dirent := d.Dirents[i]
		pathname := dirent.Path
		if dirent.IsDir {
			pathname += "/"
		}
		// load the children in case the dirent is expanded
		if dirent.Expanded {
			pathname += "\n" + dirent.Print(state)
		}
		prefixspace := ""
		for j := uint(1); j < dirent.Level; j++ {
			prefixspace += "  "
		}
		cursorDisplay := " "
		if state.getCurrentDirent().Equals(dirent) {
			cursorDisplay = ">"
		}
		s += fmt.Sprintf("%s%s%s", cursorDisplay, prefixspace, pathname)
		if i < len(d.Dirents)-1 {
			s += "\n"
		}
	}
	return s
}

// get the bottom dirent in the tree
// by following the last dirent in the list
func (d Dirent) getBottomDirent() *Dirent {
	if d.IsDir && !d.Expanded || !d.IsDir { // file or unexpanded directory
		return &d
	}
	return d.Dirents[len(d.Dirents)-1].getBottomDirent()
}
