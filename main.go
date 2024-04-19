package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/twonull/go-ff/oodle"
)

// FFHeaderMagic and FFPatchHeaderMagic are the magic bytes for regular and patch FastFiles
var (
	FFHeaderMagic      = []byte{0xFF, 0x37, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	FFPatchHeaderMagic = []byte{0xFE, 0x37, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type BlockHeader struct {
	decompressedDataLen uint32
	compressionType     byte
}

type Block struct {
	compressedLen   uint32
	decompressedLen uint32
}

func decompressBlocks(reader io.Reader, writer io.Writer) error {
	var headerBlock BlockHeader
	if err := binary.Read(reader, binary.LittleEndian, &headerBlock.decompressedDataLen); err != nil {
		return err
	}
	reader.Read(make([]byte, 3)) // Skip 3 bytes
	if _, err := reader.Read([]byte{headerBlock.compressionType}); err != nil {
		return err
	}

	var decompressedDataCount int64
	loopCount := 1
	var block Block
	for decompressedDataCount < int64(headerBlock.decompressedDataLen) {
		if err := binary.Read(reader, binary.LittleEndian, &block.compressedLen); err != nil {
			return err
		}
		block.compressedLen = (block.compressedLen + 3) & 0xFFFFFFC // Calc alignment
		if err := binary.Read(reader, binary.LittleEndian, &block.decompressedLen); err != nil {
			return err
		}
		reader.Read(make([]byte, 4)) // Skip 4 bytes

		compressedData := make([]byte, block.compressedLen)
		if _, err := reader.Read(compressedData); err != nil {
			return err
		}

		var decompressedData []byte
		switch headerBlock.compressionType {
		case 1: // Decompress None
			decompressedData = compressedData
		case 4, 5:
			return fmt.Errorf("unimplemented compression type: %d", headerBlock.compressionType)
		case 6, 7, 12, 13, 14, 15, 16, 17: // Decompress Oodle
			decompressedData, _ = oodle.Decompress(compressedData, int64(block.decompressedLen))
		default:
			return fmt.Errorf("unknown compression type: %d", headerBlock.compressionType)
		}

		if decompressedData == nil {
			return fmt.Errorf("decompressor returned nil")
		}

		if _, err := writer.Write(decompressedData); err != nil {
			return err
		}

		decompressedDataCount += int64(block.decompressedLen)
		loopCount++
	}

	return nil
}

func decompressRegularFF(reader io.Reader, writer io.Writer) error {
	// Move to start of the signed fast file
	fastFileStart := int64(0xDC)
	if _, err := reader.(*os.File).Seek(fastFileStart, io.SeekStart); err != nil {
		return err
	}

	// Skip SHA256 hashes
	if _, err := reader.(*os.File).Seek(0x8000, io.SeekCurrent); err != nil {
		return err
	}

	return decompressBlocks(reader, writer)
}

func decompressPatchFF(reader io.Reader, writer io.Writer) error {
	if _, err := reader.(*os.File).Seek(0x28, io.SeekStart); err != nil {
		return err
	}

	var emptyPatchFlag uint32
	if err := binary.Read(reader, binary.LittleEndian, &emptyPatchFlag); err != nil {
		return err
	}

	if emptyPatchFlag == 0 { // Check if it's an empty patch file, exit if it is
		return nil
	}

	// Move to start of the signed fast file
	fastFileStart := int64(0x1EC)
	if _, err := reader.(*os.File).Seek(fastFileStart, io.SeekStart); err != nil {
		return err
	}

	// Skip SHA256 hashes
	if _, err := reader.(*os.File).Seek(0x8000, io.SeekCurrent); err != nil {
		return err
	}

	return decompressBlocks(reader, writer)
}

func decompress(path, outPath string) error {
	reader, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %w", path, err)
	}
	defer reader.Close()

	writer, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", outPath, err)
	}
	defer writer.Close()

	FFMagic := make([]byte, 8)
	if _, err := reader.Read(FFMagic); err != nil {
		return err
	}

	if bytes.Equal(FFMagic, FFHeaderMagic) {
		return decompressRegularFF(reader, writer)
	} else if bytes.Equal(FFMagic, FFPatchHeaderMagic) {
		return decompressPatchFF(reader, writer)
	} else {
		return fmt.Errorf("invalid FF magic")
	}
}

// getDiffResultSizeFromFile gets the resulting size of a fastfile when patched with a patch fastfile
// patchFilePath must be the path to a patch file
func getDiffResultSizeFromFile(patchFilePath string) (int, error) {
	reader, err := os.Open(patchFilePath)
	if err != nil {
		return 0, fmt.Errorf("unable to open file %s: %w", patchFilePath, err)
	}
	defer reader.Close()

	if _, err := reader.Seek(0x140, io.SeekStart); err != nil {
		return 0, err
	}

	var diffResultSize uint32
	if err := binary.Read(reader, binary.LittleEndian, &diffResultSize); err != nil {
		return 0, err
	}

	return int(diffResultSize), nil
}

func main() {
	// Example usage
	if err := decompress("input.ff", "output.bin"); err != nil {
		fmt.Println(err)
		return
	}

	diffResultSize, err := getDiffResultSizeFromFile("patch.ff")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Diff result size: %d\n", diffResultSize)
}
