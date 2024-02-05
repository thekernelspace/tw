package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/denormal/go-gitignore"
	"github.com/thekernelspace/tw/icons"
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

	ignore, errGitignore := gitignore.NewFromFile(".gitignore")

	for pos, fi := range fileInfos {
		// don't take the ones matched by ignore patterns
		fidirent := Dirent{
			fi:          fi,
			Level:       d.Level + 1,
			Parent:      d,
			PosInParent: pos,
			Expanded:    false, // don't show by default
		}

		if !globalConfig.ShowHidden {
			if isObliviousPattern(fidirent.fi) {
				continue
			}
			if errGitignore == nil {
				match := ignore.Match(fidirent.Path())
				if match != nil {
					if match.Ignore() {
						log.Printf("Ignoring %s\n", fidirent.Path())
						continue
					}
				}
			}
		}
		d.Dirents = append(d.Dirents, fidirent)
	}

	// folders first then files, alphabetically
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

	// sort alphabetically each dirs and files
	alphabetNameSort(dirs)
	alphabetNameSort(files)

	finalDirs := append(dirs, files...)
	for i := range finalDirs {
		finalDirs[i].PosInParent = i
	}
	return finalDirs
}

func alphabetNameSort(dirents []Dirent) {
	sort.Slice(dirents, func(i, j int) bool {
		return dirents[i].Name() < dirents[j].Name()
	})
}

// get the bottom dirent in the tree
// by following the last dirent in the list
func (d Dirent) getBottomDirent() *Dirent {
	if d.IsDir() && !d.Expanded || !d.IsDir() { // file or unexpanded directory
		return &d
	}
	return d.Dirents[len(d.Dirents)-1].getBottomDirent()
}

func (d Dirent) getIcon(iconMode string) string {
	icon, color := icons.GetIcon(
		d.fi.Name(),
		filepath.Ext(d.fi.Name()),
		icons.GetIndicator(d.fi.Mode()),
	)
	var fileIcon string
	switch iconMode {
	case ICONS_COLOR:
		fileIcon = lipgloss.NewStyle().Width(2).Render(fmt.Sprintf("%s%s\033[0m ", color, icon))
	case ICONS_MONO:
		fileIcon = lipgloss.NewStyle().Width(2).Render(fmt.Sprintf("%s\033[0m ", icon))
	}
	return fileIcon
}

type WalkFunc func(d Dirent)

// walk through the dirent and its children
// and call the walk function on each dirent
// exclude the dirent if it's not expanded and is a directory
func (d Dirent) Walk(walkFn WalkFunc) {
	walkFn(d)
	if !d.Expanded && d.IsDir() {
		return
	}
	for i := range d.Dirents {
		d.Dirents[i].Walk(walkFn)
	}
}

// to print the dirent and its children if any
func (d Dirent) Print(state Model) string {
	config := GetGlobalCfg()
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
		fileIcon := ""
		if config.Icons != ICONS_NONE {
			fileIcon = dirent.getIcon(config.Icons)
		}

		sniperAnnotation := ""
		direntKeyCombo := getKeyFromDirent(dirent)
		samePrefix := direntKeyCombo[:len(state.sniperKeyBuffer)] == state.sniperKeyBuffer
		if state.sniper && samePrefix {
			sniperAnnotation = fmt.Sprintf("\t[%s]", direntKeyCombo)
			sniperAnnotation = orange.Italic(true).Render(sniperAnnotation)
			log.Printf("sniper annotation for %s: %s\n", dirent.Path(), sniperAnnotation)
		}

		s += fmt.Sprintf("%s %s%s%s%s%s", cursorDisplay, prefixspace, fileIcon, pathname, sniperAnnotation, subdirtree)
		if i < len(d.Dirents)-1 {
			s += "\n"
		}
	}
	return s
}
