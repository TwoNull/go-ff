package main

type ZoneData struct {
	g_streamOutSize   int
	g_streamBlockSize [9]int
	listStringCount   int
	listStringOffset  int
	listStrings       []string
	assetsCount       int
	assetsListOffset  int
	assetsTypesCount  [53]int
	assetsDataOffset  int
}

func ParseZoneData(decompressedData []byte) *ZoneData {
	zd := &ZoneData{}

	zd.g_streamOutSize = getDword(decompressedData, 0)

	for i := 0; i < 9; i++ {
		zd.g_streamBlockSize[i] = getDword(decompressedData, 4*(i+1))
	}

	zd.listStringCount = getDword(decompressedData, 0x2C)

	zd.assetsCount = getDword(decompressedData, 0x24)
	zd.listStringOffset = 0x3C + zd.listStringCount*4

	zd.listStrings = make([]string, zd.listStringCount)

	// Getting the strings from list
	listStringStrStartOffset := zd.listStringOffset
	listStringStrEndOffset := zd.listStringOffset

	for i := 0; i < int(zd.listStringCount); i++ {
		strOffset := getDword(decompressedData, 0x3C+4*i)
		if strOffset == 0xFFFFFFFF {
			listStringStrEndOffset = findByte(decompressedData, 0x00, listStringStrStartOffset)
			zd.listStrings[i] = string(decompressedData[listStringStrStartOffset:listStringStrEndOffset])
			listStringStrStartOffset = listStringStrEndOffset + 1
		}
	}

	if zd.listStringCount > 0 {
		zd.assetsListOffset = listStringStrEndOffset + 5
	} else {
		zd.assetsListOffset = listStringStrEndOffset
	}

	// Getting the assets' count from assets list
	for i := 0; i < zd.assetsCount; i++ {
		assetType := getDword(decompressedData, zd.assetsListOffset+8*i)
		if assetType != 0xFFFFFFFF {
			zd.assetsTypesCount[assetType]++
		}
	}

	zd.assetsDataOffset = zd.assetsListOffset + zd.assetsCount*8

	return zd
}
