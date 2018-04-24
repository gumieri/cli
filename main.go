package openInEditor

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

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
