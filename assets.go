package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	ad.GetRawFiles(decompressedData)

	return ad
}

func (ad *AssetData) GetRawFiles(data []byte) {
	index := 0
	for index < len(data) {
		sep := bytes.Index(data[index:], []byte{0xFF, 0xFF, 0xFF, 0xFF}) + index
		if sep < index {
			break
		}
		if bytes.Equal(data[sep+8:sep+12], []byte{0xFF, 0xFF, 0xFF, 0xFF}) {
			endOfName := bytes.IndexByte(data[sep+12:], 0x0) + sep + 12
			if endOfName < (sep + 12) {
				break
			}
			if data[endOfName+9] == 0x78 {
				if isASCII(data[sep+12 : endOfName]) {
					assetSize := getDword(data, sep+4)
					assetName := string(data[sep+12 : endOfName])
					startOfContents := endOfName + 9

					contents := data[startOfContents : startOfContents+assetSize]

					out, err := decompressZlib(contents)
					if err != nil {
						log.Fatal(err)
					}

					ad.files = append(ad.files, FileData{
						name:           assetName,
						nameOffset:     sep + 12,
						contents:       string(out[:len(out)-1]),
						size:           len(out) - 1,
						originalSize:   assetSize,
						contentsOffset: startOfContents,
					})
				}
			}
		}
		index = sep + 4
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
			fmt.Println("error creating directory", err)
			return
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("error deleting file", err)
			return
		}
	}

	err := os.WriteFile(filePath, []byte(ad.files[index].contents), 0644)
	if err != nil {
		fmt.Println("error writing file", err)
		return
	}
}
