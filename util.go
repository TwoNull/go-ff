package main

import (
	"bytes"
	"compress/zlib"
)

func decompressZlib(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var decompressed bytes.Buffer
	_, err = decompressed.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return decompressed.Bytes(), nil
}
