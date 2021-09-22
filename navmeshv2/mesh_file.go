package navmeshv2

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro-framework/pk2"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type MeshFile struct {
	FileName              string
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
	StructOption          RtNavmeshStructOption
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
		FileName:              filename,
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
		StructOption:          RtNavmeshStructOption(utils.ByteArrayToUint32(fileContent[56:60])),
		FileContent:           fileContent,
		fileReadIndex:         60,
	}

	return &meshFile
}

func (m *MeshFile) LoadMeshObject() RtNavmeshObj {
	logrus.Tracef("loading object mesh: %s", m.FileName)
	m.fileReadIndex = int(m.OffsetNavMeshObj)
	object := NewNavmeshObj(m.FileName)
	vertices := m.loadObjectVertices()
	m.loadObjectCells(vertices, &object)
	m.loadObjectGlobalEdges(vertices, &object)
	m.loadObjectInternalEdges(vertices, &object)
	m.loadObjectEvents(&object)
	m.loadObjectGrid(&object)

	return object
}

func (m *MeshFile) loadObjectVertices() []*math32.Vector3 {
	// Vertices
	vertexCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	m.fileReadIndex += 4
	vertices := make([]*math32.Vector3, vertexCount)
	for i := 0; i < int(vertexCount); i++ {
		vertices[i] = &math32.Vector3{
			X: utils.Float32FromByteArray(m.FileContent[m.fileReadIndex : m.fileReadIndex+4]),
			Y: utils.Float32FromByteArray(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+8]),
			Z: utils.Float32FromByteArray(m.FileContent[m.fileReadIndex+8 : m.fileReadIndex+12]),
		}
		// 1 byte sin/cos cache
		m.fileReadIndex += 13
	}

	return vertices
}

func (m *MeshFile) loadObjectCells(vertices []*math32.Vector3, object *RtNavmeshObj) {
	// Cells
	cellCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	object.Cells = make([]RtNavmeshCellTri, cellCount)
	m.fileReadIndex += 4
	for i := 0; i < int(cellCount); i++ {
		v1Idx := binary.LittleEndian.Uint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
		v2Idx := binary.LittleEndian.Uint16(m.FileContent[m.fileReadIndex+2 : m.fileReadIndex+4])
		v3Idx := binary.LittleEndian.Uint16(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+6])
		flag := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+6 : m.fileReadIndex+8])

		tri := Triangle{
			A: vertices[v1Idx],
			B: vertices[v2Idx],
			C: vertices[v3Idx],
		}
		object.Cells[i] = RtNavmeshCellTri{
			RtNavmeshCellBase: RtNavmeshCellBase{
				Index: i,
				Mesh:  nil,
			},
			Triangle: tri,
			Flag:     int16(flag),
		}
		m.fileReadIndex += 8
		if m.StructOption.IsCell() {
			// Event Zone
			// eventZone := m.FileContent[m.fileReadIndex] // byte
			m.fileReadIndex++
		}
	}
}

func (m *MeshFile) loadObjectGlobalEdges(vertices []*math32.Vector3, object *RtNavmeshObj) {
	// Global Edges
	globalEdgeCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	m.fileReadIndex += 4
	object.GlobalEdges = make([]RtNavmeshEdgeGlobal, globalEdgeCount)
	for i := 0; i < int(globalEdgeCount); i++ {
		v1Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
		v2Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+2 : m.fileReadIndex+4])
		srcCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+6])
		destCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+6 : m.fileReadIndex+8])
		flag := m.FileContent[m.fileReadIndex+8] | 1<<3
		m.fileReadIndex += 9

		gEdge := RtNavmeshEdgeGlobal{
			RtNavmeshEdgeBase: RtNavmeshEdgeBase{
				RtNavmeshEdgeMeshType: RtNavmeshEdgeMeshTypeObject,
				Mesh:                  object,
				Index:                 i,
				Line: LineSegment{
					A: vertices[v1Idx],
					B: vertices[v2Idx],
				},
				Flag:         RtNavmeshEdgeFlag(flag),
				SrcDirection: -1,
				DstDirection: -1,
				SrcCellIndex: int(srcCellIdx),
				DstCellIndex: int(destCellIdx),
				SrcCell:      nil,
				DstCell:      nil,
			},
		}

		if m.StructOption.IsEdge() {
			// TODO
			// eventZoneFlag := m.FileContent[m.fileReadIndex+10]
			// gEdge.EventZoneFlag = eventZoneFlag
			m.fileReadIndex++
		}

		object.GlobalEdges[i] = gEdge
	}
}

func (m *MeshFile) loadObjectInternalEdges(vertices []*math32.Vector3, object *RtNavmeshObj) {
	// Internal Edges
	internalEdgeCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	m.fileReadIndex += 4
	object.InternalEdges = make([]RtNavmeshEdgeInternal, internalEdgeCount)
	for i := 0; i < int(internalEdgeCount); i++ {
		v1Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
		v2Idx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+2 : m.fileReadIndex+4])
		srcCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+6])
		destCellIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex+6 : m.fileReadIndex+8])
		flag := m.FileContent[m.fileReadIndex+8] | 1<<2
		m.fileReadIndex += 9

		iEdge := RtNavmeshEdgeInternal{
			RtNavmeshEdgeBase: RtNavmeshEdgeBase{
				RtNavmeshEdgeMeshType: RtNavmeshEdgeMeshTypeObject,
				Mesh:                  object,
				Index:                 i,
				Line: LineSegment{
					A: vertices[v1Idx],
					B: vertices[v2Idx],
				},
				Flag:         RtNavmeshEdgeFlag(flag),
				SrcDirection: -1,
				DstDirection: -1,
				SrcCellIndex: int(srcCellIdx),
				DstCellIndex: int(destCellIdx),
				SrcCell:      nil,
				DstCell:      nil,
			},
		}

		if m.StructOption.IsEdge() {
			// TODO
			//eventZoneFlag := m.FileContent[m.fileReadIndex]
			//iEdge.EventZoneFlag = eventZoneFlag
			m.fileReadIndex++
		}

		object.InternalEdges[i] = iEdge
	}
}

func (m *MeshFile) loadObjectEvents(object *RtNavmeshObj) {
	// Events
	if m.StructOption.IsEvent() {
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

func (m *MeshFile) loadObjectGrid(object *RtNavmeshObj) {
	// Grid
	// Grid origin
	originX := utils.Float32FromByteArray(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
	originZ := utils.Float32FromByteArray(m.FileContent[m.fileReadIndex+4 : m.fileReadIndex+8])
	// tile width
	width := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex+8 : m.fileReadIndex+12])
	height := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex+12 : m.fileReadIndex+16])
	gridRect := NewRectangle(originX, originZ, float32(width)*RtNavmeshObjGridTileWidth, float32(height)*RtNavmeshObjGridTileHeight)
	gridTileCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex+16 : m.fileReadIndex+20])
	grid := RtNavmeshObjGrid{
		object:    object,
		X:         originX,
		Y:         originZ,
		Width:     int(width),
		Height:    int(height),
		Rectangle: gridRect,
	}
	m.fileReadIndex += 20

	gridTiles := make([]RtNavmeshObjGridTile, gridTileCount)
	for i := 0; i < int(gridTileCount); i++ {
		gEdgeCount := utils.ByteArrayToUint32(m.FileContent[m.fileReadIndex : m.fileReadIndex+4])
		gridTile := RtNavmeshObjGridTile{
			Index:         i,
			X:             i % int(width),
			Y:             i / int(height),
			grid:          grid,
			GlobalEdges:   make([]RtNavmeshEdgeGlobal, gEdgeCount),
			InternalEdges: make([]RtNavmeshEdgeInternal, 0),
			Cells:         make([]RtNavmeshCellTri, 0),
		}
		gridTile.Rectangle = NewRectangle(grid.X+(float32(gridTile.X)*RtNavmeshObjGridTileWidth), grid.Y+(float32(gridTile.Y)*RtNavmeshObjGridTileHeight), RtNavmeshObjGridTileWidth, RtNavmeshObjGridTileHeight)
		m.fileReadIndex += 4
		for j := 0; j < int(gEdgeCount); j++ {
			gEdgeIdx := utils.ByteArrayToUint16(m.FileContent[m.fileReadIndex : m.fileReadIndex+2])
			gridTile.GlobalEdges[j] = object.GlobalEdges[gEdgeIdx]
			m.fileReadIndex += 2
		}

		gridTiles[i] = gridTile
	}
	grid.Tiles = gridTiles

	for i := 0; i < len(grid.Tiles); i++ {
		tile := gridTiles[i]

		for _, cell := range object.Cells {
			if tile.Rectangle.IntersectsTriangle(cell.Triangle) {
				tile.AddCell(cell)
			}
		}

		for _, edge := range object.InternalEdges {
			if tile.Rectangle.IntersectsLine(edge.Line) {
				tile.AddInternalEdge(edge)
			}
		}
	}

	object.Grid = grid
}
