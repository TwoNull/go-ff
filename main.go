package main

import (
	"bytes"
	"log"
	"os"
)

func main() {
	ff := FFData{
		filePath: "./test/mp_nuked.ff",
	}

	compressedData, err := ff.Parse()
	if err != nil {
		log.Fatal(err)
	}

	var decompressedData []byte

	// Todo: Implement Oodle
	if bytes.Equal(ff.algorithm, FFHeaderZlib) {
		decompressedData, err = decompressZlib(compressedData)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("unsupported compression algorithm")
	}

	ParseZoneData(decompressedData)

	os.WriteFile("zone", decompressedData, 0644)

	ad := ParseAssetsData(decompressedData)

	ad.SaveAllFiles("./test/mp_nuked")
}
