package main

import (
	"encoding/binary"
)

func findBytes(bytes, seekBytes []byte, offset int) int {
	currentByte := 0
	for i := offset; i < len(bytes); i++ {
		if bytes[i] == seekBytes[currentByte] {
			currentByte++
		} else {
			currentByte = 0
		}

		if currentByte >= len(seekBytes) {
			return i - (len(seekBytes) - 1)
		}
	}
	return -1
}

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

func findByteBackward(bytes []byte, seekByte byte, offset int) int {
	if offset > len(bytes) {
		return -1
	}

	for i := offset; i >= 0; i-- {
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
