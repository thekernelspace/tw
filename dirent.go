package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Dirent struct {
	fi          fs.FileInfo
	Dirents     []Dirent
	Parent      *Dirent // the parent directory; nil if root
	PosInParent int     // the index in the parent's dirent list
	Level       uint    // the level of the dirent in the tree; 0 = root, and so on
	Expanded    bool    // whether to show the dirent in the UI
}

func (d Dirent) Path() string {
	if d.Parent == nil {
		abspath, _ := filepath.Abs(d.fi.Name())
		return abspath
	}
	abspath := filepath.Join(d.Parent.Path(), d.fi.Name())
	return abspath
}

func (d Dirent) Name() string {
	return d.fi.Name()
}

func (d Dirent) IsDir() bool {
	return d.fi.IsDir()
}

func (d Dirent) IsFile() bool {
	return !d.IsDir()
}

func (d Dirent) Equals(other Dirent) bool {
	return d.Path() == other.Path()
}

// if the dirent is a directory, return the directory entries
// otherwise, return an empty slice
func (d *Dirent) LoadDirents() {
	if !d.IsDir() {
		return
	}
	dir, err := os.Open(d.Path())
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	for pos, fi := range fileInfos {
		d.Dirents = append(d.Dirents, Dirent{
			fi:          fi,
			Level:       d.Level + 1,
			Parent:      d,
			PosInParent: pos,
			Expanded:    false, // don't show by default
		})
	}

	// folders first then files
	d.Dirents = sortDirents(d.Dirents)

	log.Printf("Loaded %d dirents for %s\n", len(d.Dirents), d.Path())
}

func sortDirents(dirents []Dirent) []Dirent {
	var dirs []Dirent
	var files []Dirent
	for _, dirent := range dirents {
		if dirent.IsDir() {
			dirs = append(dirs, dirent)
		} else {
			files = append(files, dirent)
		}
	}
	finalDirs := append(dirs, files...)
	for i := range finalDirs {
		finalDirs[i].PosInParent = i
	}
	return finalDirs
}

func (d Dirent) Print(state Model) string {
	s := ""
	for i := range d.Dirents {
		// Render the dirent
		dirent := d.Dirents[i]
		pathname := dirent.Name()
		if dirent.IsDir() {
			pathname += "/"
		}
		// load the children in case the dirent is expanded
		subdirtree := ""
		if dirent.Expanded {
			subdirtree += "\n" + dirent.Print(state)
		}
		prefixspace := ""
		for j := uint(1); j < dirent.Level; j++ {
			prefixspace += "  "
		}
		cursorDisplay := " "
		if state.getCurrentDirent().Equals(dirent) {
			cursorDisplay = teal.Render(">")
			pathname = teal.Render(pathname)
		}
		s += fmt.Sprintf("%s %s%s%s", cursorDisplay, prefixspace, pathname, subdirtree)
		if i < len(d.Dirents)-1 {
			s += "\n"
		}
	}
	return s
}

// get the bottom dirent in the tree
// by following the last dirent in the list
func (d Dirent) getBottomDirent() *Dirent {
	if d.IsDir() && !d.Expanded || !d.IsDir() { // file or unexpanded directory
		return &d
	}
	return d.Dirents[len(d.Dirents)-1].getBottomDirent()
}
