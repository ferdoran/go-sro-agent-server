package navmesh

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro-framework/math"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type MeshFile struct {
	OffsetVertex          uint32
	OffsetSkin            uint32
	OffsetFace            uint32
	OffsetClothVertex     uint32
	OffsetClothEdge       uint32
	OffsetBoundingBox     uint32
	OffsetOcclusionPortal uint32
	OffsetNavMeshObj      uint32
	OffsetSkinnedNavMesh  uint32
	Offset9               uint32
	Int1                  uint32
	StructOption          uint32
	FileContent           []byte
	fileReadIndex         int
}

func LoadMeshFile(filename string, reader *pk2.Pk2Reader) *MeshFile {
	filename = strings.ReplaceAll(filename, "\\", string(os.PathSeparator))
	filename = strings.ReplaceAll(filename, "/", string(os.PathSeparator))
	logrus.Tracef("loading mesh %s", filename)
	fileContent, err := reader.ReadFile(filename)
	if err != nil {
		logrus.Panic(err)
	}

	header := string(fileContent[:12])

	if header != "JMXVBMS 0110" {
		logrus.Panicf("Header did not start with JMXVBMS 0110. Got %s\n", header)
	}

	meshFile := MeshFile{
		OffsetVertex:          utils.ByteArrayToUint32(fileContent[12:16]),
		OffsetSkin:            utils.ByteArrayToUint32(fileContent[16:20]),
		OffsetFace:            utils.ByteArrayToUint32(fileContent[20:24]),
		OffsetClothVertex:     utils.ByteArrayToUint32(fileContent[24:28]),
		OffsetClothEdge:       utils.ByteArrayToUint32(fileContent[28:32]),
		OffsetBoundingBox:     utils.ByteArrayToUint32(fileContent[32:36]),
		OffsetOcclusionPortal: utils.ByteArrayToUint32(fileContent[36:40]),
		OffsetNavMeshObj:      utils.ByteArrayToUint32(fileContent[40:44]),
		OffsetSkinnedNavMesh:  utils.ByteArrayToUint32(fileContent[44:48]),
		Offset9:               utils.ByteArrayToUint32(fileContent[48:52]),
		Int1:                  utils.ByteArrayToUint32(fileContent[52:56]),
		StructOption:          utils.ByteArrayToUint32(fileContent[56:60]),
		FileContent:           fileContent,
		fileReadIndex:         60,
	}

	return &meshFile
}

func (m *MeshFile) LoadMeshObject(object *Object) {
	m.fileReadIndex = int(m.OffsetNavMeshObj)
	m.loadObjectVertices(object)
	m.loadObjectCells(object)
	m.loadObjectGlobalEdges(object)
	m.loadObjectInternalEdges(object)
	m.loadObjectEvents(object)
	m.loadObjectGrid(object)
}

func (m *MeshFile) loadObjectVertices(object *Object) {
	// Vertices
	vertexCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	m.fileReadIndex += 4
	object.Vertices = make([]*math32.Vector3, vertexCount)
	for i := 0; i < int(vertexCount); i++ {
		object.Vertices[i] = &math32.Vector3{
			X: utils.Float32FromByteArray(m.FileContent[m.fileReadIndex : m.fileReadIndex+4]),
			Y: utils.Float32FromByteArray(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+8]),
			Z: utils.Float32FromByteArray(m.FileContent[m.fileReadIndex+8 : m.fileReadIndex+12]),
		}
		// 1 byte sin/cos cache
		m.fileReadIndex += 13
	}
}

func (m *MeshFile) loadObjectCells(object *Object) {
	// Cells
	cellCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	object.Cells = make([]*ObjectCell, cellCount)
	m.fileReadIndex += 4
	for i := 0; i < int(cellCount); i++ {
		v1Idx := binary.LittleEndian.Uint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
		v2Idx := binary.LittleEndian.Uint16(m.FileContent[m.fileReadIndex+2 : m.fileReadIndex+4])
		v3Idx := binary.LittleEndian.Uint16(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+6])
		flag := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+6 : m.fileReadIndex+8])

		tri := &math.Triangle{
			A: object.Vertices[v1Idx],
			B: object.Vertices[v2Idx],
			C: object.Vertices[v3Idx],
		}
		object.Cells[i] = &ObjectCell{
			Triangle: tri,
			Index:    i,
			Flag:     flag,
		}
		m.fileReadIndex += 8
		if m.StructOption&2 != 0 {
			// Event Zone
			m.fileReadIndex++
		}
	}
}

func (m *MeshFile) loadObjectGlobalEdges(object *Object) {
	// Global Edges
	globalEdgeCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	m.fileReadIndex += 4
	object.GlobalEdges = make([]*ObjectGlobalEdge, globalEdgeCount)
	for i := 0; i < int(globalEdgeCount); i++ {
		v1Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
		v2Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+2 : m.fileReadIndex+4])
		srcCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+6])
		destCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+6 : m.fileReadIndex+8])
		flag := m.FileContent[m.fileReadIndex+8]
		m.fileReadIndex += 9

		gEdge := ObjectGlobalEdge{
			A:                    object.Vertices[v1Idx],
			B:                    object.Vertices[v2Idx],
			SourceCellIndex:      int(srcCellIdx),
			DestinationCellIndex: int(destCellIdx),
			SourceMeshIndex:      -1,
			DestinationMeshIndex: -1,
			SourceDirection:      -1,
			DestinationDirection: -1,
			Flag:                 flag,
			EventZoneFlag:        0,
		}

		if m.StructOption&1 != 0 {
			eventZoneFlag := m.FileContent[m.fileReadIndex+10]
			gEdge.EventZoneFlag = eventZoneFlag
			m.fileReadIndex++
		}

		object.GlobalEdges[i] = &gEdge
	}
}

func (m *MeshFile) loadObjectInternalEdges(object *Object) {
	// Internal Edges
	internalEdgeCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	m.fileReadIndex += 4
	object.InternalEdges = make([]*ObjectInternalEdge, internalEdgeCount)
	for i := 0; i < int(internalEdgeCount); i++ {
		v1Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
		v2Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+2 : m.fileReadIndex+4])
		srcCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+6])
		destCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+6 : m.fileReadIndex+8])
		flag := m.FileContent[m.fileReadIndex+8]
		m.fileReadIndex += 9

		iEdge := ObjectInternalEdge{
			A:                    object.Vertices[v1Idx],
			B:                    object.Vertices[v2Idx],
			SourceCellIndex:      int(srcCellIdx),
			DestinationCellIndex: int(destCellIdx),
			SourceDirection:      -1,
			DestinationDirection: -1,
			Flag:                 flag,
			EventZoneFlag:        0,
		}

		if m.StructOption&1 != 0 {
			eventZoneFlag := m.FileContent[m.fileReadIndex]
			iEdge.EventZoneFlag = eventZoneFlag
			m.fileReadIndex++
		}

		object.InternalEdges[i] = &iEdge
	}
}

func (m *MeshFile) loadObjectEvents(object *Object) {
	// Events
	if m.StructOption&4 != 0 {
		eventCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
		m.fileReadIndex += 4
		object.Events = make([]string, eventCount)
		for i := 0; i < int(eventCount); i++ {
			strLen := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
			object.Events[i] = string(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+4+int(strLen)])

			m.fileReadIndex += 4 + int(strLen)
		}
	}
}

func (m *MeshFile) loadObjectGrid(object *Object) {
	// Grid
	// Grid origin
	originX := utils.Float32FromByteArray(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	originZ := utils.Float32FromByteArray(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+8])
	// tile width
	width := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex+8 : m.fileReadIndex+12])
	height := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex+12 : m.fileReadIndex+16])
	gridTileCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex+16 : m.fileReadIndex+20])
	m.fileReadIndex += 20
	gridTiles := make([]*ObjectTile, gridTileCount)
	for i := 0; i < int(gridTileCount); i++ {
		gEdgeCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
		gridTile := NewObjectTile(originX, originZ, i, int(gEdgeCount))
		m.fileReadIndex += 4
		for j := 0; j < int(gEdgeCount); j++ {
			gEdgeIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
			gridTile.GlobalEdges[j] = object.GlobalEdges[gEdgeIdx]
			m.fileReadIndex += 2
		}

		gridTiles[i] = gridTile
	}

	for i := 0; i < len(gridTiles); i++ {
		tile := gridTiles[i]

		for _, cell := range object.Cells {
			intersects, _ := cell.Triangle.ToTriangle2().IntersectsRect(tile.Rectangle)
			if intersects {
				tile.Cells = append(tile.Cells, cell)
			}
		}

		for _, edge := range object.InternalEdges {
			intersects, _ := tile.IntersectsLine(edge.ToLine2())

			if intersects {
				tile.InternalEdges = append(tile.InternalEdges, edge)
			}
		}
	}

	object.Grid = &ObjectGrid{
		Origin:    math32.NewVector3(originX, 0, originZ),
		Width:     int(width),
		Height:    int(height),
		GridTiles: gridTiles,
	}
}
