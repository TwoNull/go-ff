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

func getDword(allBytes []byte, startOffset int, reverseBytes bool) int {
	tableIndexBytes := make([]byte, 4)
	copy(tableIndexBytes, allBytes[startOffset:startOffset+4])

	if !reverseBytes {
		reverse(tableIndexBytes)
	}

	return int(binary.LittleEndian.Uint32(tableIndexBytes))
}

func getString(bytes []byte, offset, endOffset int) string {
	byteSelection := bytes[offset:endOffset]
	return string(byteSelection)
}

func reverse(bytes []byte) {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
}

/*func removeBytes(bytes []byte, startOffset, endOffset int) []byte {
	result := make([]byte, 0, len(bytes)-(endOffset-startOffset))
	result = append(result, bytes[:startOffset]...)
	result = append(result, bytes[endOffset:]...)
	return result
}

func addBytes(bytes, add []byte, startOffset int) []byte {
	newLength := len(bytes) + len(add)
	result := make([]byte, newLength)
	copy(result, bytes[:startOffset])
	copy(result[startOffset:], add)
	copy(result[startOffset+len(add):], bytes[startOffset:])
	return result
}

func replaceBytes(bytes []byte, offset int, replacement []byte, finishWithNullByte bool) []byte {
	copy(bytes[offset:], replacement)
	if finishWithNullByte {
		bytes[offset+len(replacement)] = 0
	}
	return bytes
}

func getBytes(bytes []byte, startOffset, endOffset int) []byte {
	end := len(bytes)
	if endOffset != -1 {
		end = endOffset
	}
	return bytes[startOffset:end]
}

func setBytes(bytes []byte, startOffset, length int, setByte byte) []byte {
	for i := startOffset; i < startOffset+length; i++ {
		bytes[i] = setByte
	}
	return bytes
}

func countBytes(bytes, seekBytes []byte, offset, endOffset int) int {
	endPos := len(bytes)
	if endOffset != -1 {
		endPos = endOffset
	}
	currentByte := 0
	count := 0
	for i := offset; i < endPos; i++ {
		if bytes[i] == seekBytes[currentByte] {
			currentByte++
		} else {
			currentByte = 0
		}

		if currentByte >= len(seekBytes) {
			count++
			currentByte = 0
		}
	}
	return count
}*/
