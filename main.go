package openineditor

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// File executes a text editor to edit a file
func File(editorCommand string, filePath string) (err error) {
	editor := exec.Command(editorCommand, filePath)

	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr

	err = editor.Start()
	if err != nil {
		return
	}

	err = editor.Wait()
	if err != nil {
		return
	}

	return
}

// GetContentFromTemporaryFile executes a text editor
// to edit a temporary file with an specified name
// and return its content
func GetContentFromTemporaryFile(editorCommand string, fileName string) (text string, err error) {
	filePath := filepath.Join(os.TempDir(), fileName)

	tmpFile, err := os.Create(filePath)
	if err != nil {
		return
	}

	tmpFile.Close()

	err = File(editorCommand, filePath)
	if err != nil {
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	text = string(content)

	return
}
