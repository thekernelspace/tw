package main

import (
	"os"
	"os/exec"
)

func EditorOpenFile(filepath string) *exec.Cmd {
	// Open file with the default $EDITOR
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	c := exec.Command(editor, filepath)
	return c
}
