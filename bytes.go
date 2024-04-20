package main

import (
	"encoding/binary"
)

func findByte(bytes []byte, seekByte byte, offset int) int {
	if offset > len(bytes) {
		return -1
	}

	for i := offset; i < len(bytes); i++ {
		if bytes[i] == seekByte {
			return i
		}
	}
	return -1
}

func getDword(bytes []byte, startOffset int) int {
	tableIndexBytes := make([]byte, 4)
	copy(tableIndexBytes, bytes[startOffset:startOffset+4])

	return int(binary.LittleEndian.Uint32(tableIndexBytes))
}
