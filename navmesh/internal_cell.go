package navmesh

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"log"
)

type NavigationCells struct {
	CellCount     uint32
	OpenCellCount uint32
	Cells         []NavigationCell
}

type NavigationCell struct {
	Min      math32.Vector2
	Max      math32.Vector2
	ObjCount byte
	Objects  []uint16
}

type NavMeshInternalEdges struct {
	InternalEdgeCount uint32
	InternalEdges     []NavMeshInternalEdge
}

type NavMeshInternalEdge struct {
	Min             math32.Vector2
	Max             math32.Vector2
	EdgeFlag        byte // OR it with EdgeFlag.Internal = 4, + check if Bit5 and Bit6 are set (EdgeFlag & Bit5) and print it
	AssocDirection0 byte
	AssocDirection1 byte
	AssocCell0      uint16
	AssocCell1      uint16
}

//enum EdgeFlag
//{
//	None = 0,
//	BlockDst2Src = 1,
//	BlockSrc2Dst = 2,
//	Blocked = BlockDst2Src | BlockSrc2Dst,
//	Internal = 4,
//	Global = 8,
//	Bridge = 16,
//	Bit5 = 32,
//	Bit6 = 64,
//	Siege = 128,
//}

func loadNavigationCells(content []byte, readIndex *int) NavigationCells {
	navCells := NavigationCells{
		CellCount:     binary.LittleEndian.Uint32(content[*readIndex : *readIndex+4]),
		OpenCellCount: binary.LittleEndian.Uint32(content[*readIndex+4 : *readIndex+8]),
		Cells:         make([]NavigationCell, 0),
	}
	*readIndex += 8

	for i := 0; i < int(navCells.CellCount); i++ {
		minVectorOffset := *readIndex
		maxVectorOffset := minVectorOffset + 8
		objCountOffset := maxVectorOffset + 8
		*readIndex += 17

		cell := NavigationCell{
			Min: math32.Vector2{
				X: utils.Float32FromByteArray(content[minVectorOffset : minVectorOffset+4]),
				Y: utils.Float32FromByteArray(content[minVectorOffset+4 : maxVectorOffset]),
			},
			Max: math32.Vector2{
				X: utils.Float32FromByteArray(content[maxVectorOffset : maxVectorOffset+4]),
				Y: utils.Float32FromByteArray(content[maxVectorOffset+4 : objCountOffset]),
			},
			ObjCount: content[objCountOffset],
			Objects:  make([]uint16, 0),
		}

		for j := 0; j < int(cell.ObjCount); j++ {
			cell.Objects = append(cell.Objects, binary.LittleEndian.Uint16(content[*readIndex:*readIndex+2]))
			*readIndex += 2
		}

		navCells.Cells = append(navCells.Cells, cell)
	}
	return navCells
}

func loadInternalEdges(content []byte, readIndex *int) NavMeshInternalEdges {
	internalEdges := NavMeshInternalEdges{
		InternalEdgeCount: binary.LittleEndian.Uint32(content[*readIndex : *readIndex+4]),
		InternalEdges:     make([]NavMeshInternalEdge, 0),
	}

	*readIndex += 4

	for i := 0; i < int(internalEdges.InternalEdgeCount); i++ {
		minVectorOffset := *readIndex
		maxVectorOffset := minVectorOffset + 8
		edgeFlagOffset := maxVectorOffset + 8
		assocDir0Offset := edgeFlagOffset + 1
		assocDir1Offset := assocDir0Offset + 1
		assocCell0Offset := assocDir1Offset + 1
		assocCell1Offset := assocCell0Offset + 2

		*readIndex += 23

		if val := content[edgeFlagOffset] & 32; val != 0 {
			log.Printf("Internal EdgeType & Bit5 %v", val)
		}
		if val := content[edgeFlagOffset] & 64; val != 0 {
			log.Printf("Internal EdgeType & Bit6 %v", val)
		}

		internalEdge := NavMeshInternalEdge{
			Min: math32.Vector2{
				X: utils.Float32FromByteArray(content[minVectorOffset : minVectorOffset+4]),
				Y: utils.Float32FromByteArray(content[minVectorOffset+4 : maxVectorOffset]),
			},
			Max: math32.Vector2{
				X: utils.Float32FromByteArray(content[maxVectorOffset : maxVectorOffset+4]),
				Y: utils.Float32FromByteArray(content[maxVectorOffset+4 : edgeFlagOffset]),
			},
			EdgeFlag:        content[edgeFlagOffset],
			AssocDirection0: content[assocDir0Offset],
			AssocDirection1: content[assocDir1Offset],
			AssocCell0:      binary.LittleEndian.Uint16(content[assocCell0Offset:assocCell1Offset]),
			AssocCell1:      binary.LittleEndian.Uint16(content[assocCell1Offset : assocCell1Offset+2]),
		}

		internalEdges.InternalEdges = append(internalEdges.InternalEdges, internalEdge)
	}
	return internalEdges
}
