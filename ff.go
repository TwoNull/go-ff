package main

import (
	"errors"
	"io"
	"os"
)

var (
	FFHeaderZlib  = []byte("IWffu100")
	FFHeaderOodle = []byte("IWffa100")
)

var (
	FFVersionBO1 = []byte{217, 1, 0, 0}
)

type FFData struct {
	filePath string
	length   int

	algorithm []byte
	version   []byte
}

func (ff *FFData) Parse() ([]byte, error) {
	file, err := os.Open(ff.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ffContents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ff.length = len(ffContents)

	if ff.length < 12 {
		return nil, errors.New("invalid file")
	}

	ff.algorithm = ffContents[:8]
	ff.version = ffContents[8:12]

	return ffContents[12:], nil
}
