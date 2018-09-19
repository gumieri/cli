package openineditor

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Editor "object"
type Editor struct {
	Command     string
	Flags       []string
	OpenedFiles []File
}

// File "object"
type File struct {
	FileName string
	FilePath string
	Content  []byte
	File     *os.File
}

// NewTempFile create a File in the OS temporary directory
func NewTempFile(fileName string, content []byte) (f File, err error) {
	f.FileName = fileName
	f.Content = content

	err = f.CreateInTempDir()

	return
}

// CreateInTempDir create the File in the OS temporary directory
func (f *File) CreateInTempDir() (err error) {
	f.FilePath = filepath.Join(os.TempDir(), f.FileName)

	f.File, err = os.Create(f.FilePath)
	if err != nil {
		return
	}

	if len(f.Content) != 0 {
		_, err = f.File.Write(f.Content)
		if err != nil {
			return
		}
	}

	f.File.Close()

	return
}

// OpenFile executes the Editor to open the specified File
func (e *Editor) OpenFile(f File) (err error) {
	switch e.Command {
	case "subl", "code":
		e.Flags = append(e.Flags, "--wait")
	}

	args := append(e.Flags, f.FilePath)

	process := exec.Command(e.Command, args...)

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	err = process.Start()
	if err != nil {
		return
	}

	err = process.Wait()
	if err != nil {
		return
	}

	e.OpenedFiles = append(e.OpenedFiles, f)

	return
}

// OpenTempFile executes the Editor to open a TempFile
func (e *Editor) OpenTempFile(f File) (err error) {
	if f.File == nil {
		if f.CreateInTempDir() != nil {
			return
		}
	}

	err = e.OpenFile(f)
	if err != nil {
		return
	}

	f.Content, err = ioutil.ReadFile(f.FilePath)

	return
}

// LastFile return the last opened file or return error if no file was opened
func (e *Editor) LastFile() (f File, err error) {
	totalFiles := len(e.OpenedFiles)
	if totalFiles == 0 {
		err = errors.New("no files were opened")
		return
	}

	f = e.OpenedFiles[totalFiles-1]

	return
}
