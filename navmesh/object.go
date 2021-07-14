package navmesh

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	log "github.com/sirupsen/logrus"
)

type ObjectList struct {
	ObjectCount uint16
	Objects     []*Object
}

type Object struct {
	ID                  uint32
	Position            *math32.Vector3
	Type                uint16 // 0xFFFF = Static, 0x0000 = Skinned
	Yaw                 float32
	LocalUID            uint16
	Short0              uint16
	IsLarge             bool
	IsStructure         bool
	RegionID            uint16
	GlobalEdgeLinkCount uint16
	GlobalEdgeLinks     []*GlobalEdgeLink
	Vertices            []*math32.Vector3
	Cells               []*ObjectCell
	GlobalEdges         []*ObjectGlobalEdge
	InternalEdges       []*ObjectInternalEdge
	Events              []string
	Grid                *ObjectGrid
	Rotation            *math32.Quaternion
	LocalToWorld        *math32.Matrix4
	WorldToLocal        *math32.Matrix4
}

type GlobalEdgeLink struct {
	LinkObjID     uint16
	LinkObjEdgeID uint16
	EdgeID        uint16
}

func ParseNavMeshFile(filename string, fileContent []byte) NavMeshData {
	readIndex := 0
	signature := string(fileContent[:12])

	if signature != "JMXVNVM 1000" {
		log.Panic("Invalid signature")
	}
	readIndex += 12
	objList := loadObjectList(fileContent, &readIndex)
	if filename == "nv_62a8.nvm" {
		log.Infof("Found %d objects for Region 25256", objList.ObjectCount)
	}
	navCells := loadNavigationCells(fileContent, &readIndex)
	globalEdges := loadGlobalEdges(fileContent, &readIndex)
	internalEdges := loadInternalEdges(fileContent, &readIndex)
	tileMap := loadTileMap(fileContent, &readIndex)
	heightMap := loadHeightMap(fileContent, &readIndex)
	surfaceTypeMap, surfaceHeightMap := loadSurfaceMaps(fileContent, &readIndex)

	log.Tracef("Finished loading %s\n", filename)

	return NavMeshData{
		ObjectList:           objList,
		NavigationCells:      navCells,
		NavMeshGlobalEdges:   globalEdges,
		NavMeshInternalEdges: internalEdges,
		TileMap:              tileMap,
		HeightMap:            heightMap,
		SurfaceTypeMap:       surfaceTypeMap,
		SurfaceHeightMap:     surfaceHeightMap,
	}
}

func loadObjectList(content []byte, readIndex *int) ObjectList {
	objList := ObjectList{
		ObjectCount: binary.LittleEndian.Uint16(content[*readIndex : *readIndex+2]),
		Objects:     make([]*Object, binary.LittleEndian.Uint16(content[*readIndex:*readIndex+2])),
	}
	*readIndex += 2
	for i := 0; i < int(objList.ObjectCount); i++ {
		objIdOffset := *readIndex
		positionOffset := objIdOffset + 4
		objTypeOffset := positionOffset + 12
		yawOffset := objTypeOffset + 2
		localUidOffset := yawOffset + 4
		short0Offset := localUidOffset + 2
		isLargeOffset := short0Offset + 2
		isStructureOffset := isLargeOffset + 1
		regionIdOffset := isStructureOffset + 1
		globalEdgeLinkCountOffset := regionIdOffset + 2
		*readIndex += 32
		object := Object{
			ID: binary.LittleEndian.Uint32(content[objIdOffset:positionOffset]),
			Position: math32.NewVector3(
				utils.Float32FromByteArray(content[positionOffset:positionOffset+4]),
				utils.Float32FromByteArray(content[positionOffset+4:positionOffset+8]),
				utils.Float32FromByteArray(content[positionOffset+8:objTypeOffset]),
			),
			Type:                binary.LittleEndian.Uint16(content[objTypeOffset:yawOffset]),
			Yaw:                 utils.Float32FromByteArray(content[yawOffset:localUidOffset]),
			LocalUID:            binary.LittleEndian.Uint16(content[localUidOffset:short0Offset]),
			Short0:              binary.LittleEndian.Uint16(content[short0Offset:isLargeOffset]),
			IsLarge:             content[isLargeOffset] != 0,
			IsStructure:         content[isStructureOffset] != 0,
			RegionID:            binary.LittleEndian.Uint16(content[regionIdOffset:globalEdgeLinkCountOffset]),
			GlobalEdgeLinkCount: binary.LittleEndian.Uint16(content[globalEdgeLinkCountOffset : globalEdgeLinkCountOffset+2]),
			GlobalEdgeLinks:     make([]*GlobalEdgeLink, 0),
		}

		for j := 0; j < int(object.GlobalEdgeLinkCount); j++ {
			linkedObjIdOffset := *readIndex
			linkedObjEdgeIdOffset := linkedObjIdOffset + 2
			edgeIdOffset := linkedObjEdgeIdOffset + 2
			*readIndex += 6
			globalEdgeLink := GlobalEdgeLink{
				LinkObjID:     binary.LittleEndian.Uint16(content[linkedObjIdOffset:linkedObjEdgeIdOffset]),
				LinkObjEdgeID: binary.LittleEndian.Uint16(content[linkedObjEdgeIdOffset:edgeIdOffset]),
				EdgeID:        binary.LittleEndian.Uint16(content[edgeIdOffset : edgeIdOffset+2]),
			}
			object.GlobalEdgeLinks = append(object.GlobalEdgeLinks, &globalEdgeLink)
		}
		objList.Objects[i] = &object
	}
	return objList
}
