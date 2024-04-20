package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	gscFormat = []byte{'.', 'g', 's', 'c', 0}
	gsxFormat = []byte{'.', 'g', 's', 'x', 0}
	rmbFormat = []byte{'.', 'r', 'm', 'b', 0}
	cfgFormat = []byte{'.', 'c', 'f', 'g', 0}
	defFormat = []byte{'.', 'd', 'e', 'f', 0}
)

type AssetData struct {
	files []FileData
}

type FileData struct {
	name           string
	nameOffset     int
	contents       string
	size           int
	originalSize   int
	contentsOffset int
}

func ParseAssetsData(decompressedData []byte) *AssetData {
	ad := &AssetData{}

	ad.AddFiles(gsxFormat, decompressedData)
	ad.AddFiles(gscFormat, decompressedData)
	ad.AddFiles(rmbFormat, decompressedData)
	ad.AddFiles(defFormat, decompressedData)
	ad.AddFiles(cfgFormat, decompressedData)

	return ad
}

func (ad *AssetData) AddFiles(extension []byte, data []byte) {
	offset := findBytes(data, extension, 0)

	for offset != -1 {
		startOfNameOffset := findByteBackward(data, 0xFF, offset+1) + 1
		endOfNameOffset := findByte(data, 0x00, offset+1)
		assetSize := getDword(data, startOfNameOffset-8)
		assetName := string(data[startOfNameOffset:endOfNameOffset])

		startOfContents := endOfNameOffset + 9
		contents := data[startOfContents : startOfContents+assetSize]

		out, err := decompressZlib(contents)
		if err != nil {
			log.Fatal(err)
		}

		ad.files = append(ad.files, FileData{
			name:           assetName,
			nameOffset:     startOfNameOffset,
			contents:       string(out[:len(out)-1]),
			size:           len(out) - 1,
			originalSize:   assetSize,
			contentsOffset: startOfContents,
		})

		offset = findBytes(data, extension, offset+1)
	}
}

func (ad *AssetData) SaveAllFiles(baseDir string) {
	for i := 0; i < len(ad.files); i++ {
		filePath := filepath.Join(baseDir, strings.ReplaceAll(ad.files[i].name, "/", string(filepath.Separator)))

		ad.SaveFile(filePath, i)
	}
}

func (ad *AssetData) SaveFile(filePath string, index int) {
	fileDir := filepath.Dir(filePath)

	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		err := os.MkdirAll(fileDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}

	err := os.WriteFile(filePath, []byte(ad.files[index].contents), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
