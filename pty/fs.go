package pty

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func FetchDir(path string) (*[]File, error) {
	files := []File{}
	entries, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fileType := "file"
		if entry.IsDir() {
			fileType = "dir"
		}

		file := File{
			Type: fileType,
			Name: entry.Name(),
			Path: filepath.Join(wd, entry.Name()),
		}
		files = append(files, file)
	}
	return &files, err
}
