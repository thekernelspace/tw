package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

var DEBUG = os.Getenv("DEBUG") != ""

var rootDirpath string

func main() {
	if DEBUG {
		f, err := os.Create("debug.log")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println("Debug mode enabled")
	} else {
		log.SetOutput(io.Discard)
	}

	// load the config
	LoadGlobalCfg()

	// if specified an arg then use that dir
	rootDirpath = "."
	if len(os.Args) > 1 {
		rootDirpath = os.Args[1]
	}
	if !filepath.IsAbs(rootDirpath) {
		cwd, _ := os.Getwd()
		rootDirpath = filepath.Join(cwd, rootDirpath)
	}

	p := tea.NewProgram(InitModel(rootDirpath))
	if _, err := p.Run(); err != nil {
		log.Panicf("Error: %v\n", err)
		os.Exit(1)
	}
}
