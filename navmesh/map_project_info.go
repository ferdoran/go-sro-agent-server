package navmesh

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

const (
	RegionsX = 256
	RegionsY = 256

	RegionsTotal = RegionsX * RegionsY

	XSize   = 8
	XOffset = 0

	YSize               = 7
	YOffset             = XOffset + XSize
	DungeonOffset       = YOffset + YSize
	DungeonMask         = ((1 << 1) - 1) << DungeonOffset
	SignatureByteLength = 12
	DataByteLength      = 12
	TotalByteLength     = SignatureByteLength + DataByteLength + RegionsTotal/8
	MapFilePathDataPk2  = "Data" + string(os.PathSeparator) + "navmesh" + string(os.PathSeparator) + "mapinfo.mfo"
)

type MapProjectInfo struct {
	MapWidth           uint16
	MapHeight          uint16
	Short2             uint16
	Short3             uint16
	Short4             uint16
	Short5             uint16
	ActiveRegionsCount int
	MapRegions         []byte
	EnabledRegions     []uint16
}

func LoadMapProjectInfo(reader *pk2.Pk2Reader) MapProjectInfo {
	mapInfoBytes, err := reader.ReadFile(MapFilePathDataPk2)
	if err != nil {
		logrus.Panicf("error loading file %s. Error = %v", MapFilePathDataPk2, err)
	}

	if length := len(mapInfoBytes); length != TotalByteLength {
		logrus.Panicf("map info file has unexpected size (want = %d but got %d)", TotalByteLength, length)
	}

	signature := mapInfoBytes[:SignatureByteLength]
	if string(signature) != "JMXVMFO 1000" {
		log.Panicf("Invalid signature: %v\n", signature)
	}

	mapWidthBytes := mapInfoBytes[SignatureByteLength : SignatureByteLength+2]
	mapHeightBytes := mapInfoBytes[SignatureByteLength+2 : SignatureByteLength+4]
	short2Bytes := mapInfoBytes[SignatureByteLength+4 : SignatureByteLength+6]
	short3Bytes := mapInfoBytes[SignatureByteLength+6 : SignatureByteLength+8]
	short4Bytes := mapInfoBytes[SignatureByteLength+8 : SignatureByteLength+10]
	short5Bytes := mapInfoBytes[SignatureByteLength+10 : SignatureByteLength+12]
	totalRegionBytes := mapInfoBytes[SignatureByteLength+12 : TotalByteLength]

	mapProjectInfo := MapProjectInfo{
		MapWidth:           binary.LittleEndian.Uint16(mapWidthBytes),
		MapHeight:          binary.LittleEndian.Uint16(mapHeightBytes),
		Short2:             binary.LittleEndian.Uint16(short2Bytes),
		Short3:             binary.LittleEndian.Uint16(short3Bytes),
		Short4:             binary.LittleEndian.Uint16(short4Bytes),
		Short5:             binary.LittleEndian.Uint16(short5Bytes),
		ActiveRegionsCount: 0,
		MapRegions:         totalRegionBytes,
		EnabledRegions:     make([]uint16, 0),
	}

	for z := 0; z < int(mapProjectInfo.MapHeight); z++ {
		for x := 0; x < int(mapProjectInfo.MapWidth); x++ {
			if IsEnabled(byte(x), byte(z), mapProjectInfo) {
				mapProjectInfo.ActiveRegionsCount++
				mapProjectInfo.EnabledRegions = append(mapProjectInfo.EnabledRegions, binary.LittleEndian.Uint16([]byte{byte(x), byte(z)}))
			}
		}
	}

	return mapProjectInfo
}

func IsEnabled(x, z byte, mapProjectInfo MapProjectInfo) bool {
	regionShort := binary.LittleEndian.Uint16([]byte{x, z})
	if (regionShort&DungeonMask)>>DungeonOffset != 0 {
		// It's a dungeon
		return false
	}

	if int(x) >= int(mapProjectInfo.MapWidth) || int(z) >= int(mapProjectInfo.MapHeight) {
		return false
	}
	return (mapProjectInfo.MapRegions[regionShort>>3] & byte(uint16(128>>(regionShort%8)))) != 0
}
