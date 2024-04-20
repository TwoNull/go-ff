package main

import (
	"bytes"
	"compress/zlib"
	"unicode"
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

func isASCII(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
