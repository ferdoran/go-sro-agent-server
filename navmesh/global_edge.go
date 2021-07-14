package navmesh

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro-framework/utils"
	"golang.org/x/image/math/f64"
	"log"
)

type NavMeshGlobalEdges struct {
	GlobalEdgeCount uint32
	GlobalEdges     []NavMeshGlobalEdge
}

type NavMeshGlobalEdge struct {
	Min             f64.Vec2
	Max             f64.Vec2
	EdgeFlag        byte // OR it with EdgeFlag.Global = 8, + check if Bit5 and Bit6 are set (EdgeFlag & Bit5) and print it
	AssocDirection0 byte
	AssocDirection1 byte
	AssocCell0      uint16
	AssocCell1      uint16
	AssocRegion0    uint16
	AssocRegion1    uint16
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

func loadGlobalEdges(content []byte, readIndex *int) NavMeshGlobalEdges {
	globalEdges := NavMeshGlobalEdges{
		GlobalEdgeCount: binary.LittleEndian.Uint32(content[*readIndex : *readIndex+4]),
		GlobalEdges:     make([]NavMeshGlobalEdge, 0),
	}

	*readIndex += 4

	for i := 0; i < int(globalEdges.GlobalEdgeCount); i++ {
		minVectorOffset := *readIndex
		maxVectorOffset := minVectorOffset + 8
		edgeFlagOffset := maxVectorOffset + 8
		assocDir0Offset := edgeFlagOffset + 1
		assocDir1Offset := assocDir0Offset + 1
		assocCell0Offset := assocDir1Offset + 1
		assocCell1Offset := assocCell0Offset + 2
		assocRgn0Offset := assocCell1Offset + 2
		assocRgn1Offset := assocRgn0Offset + 2

		*readIndex += 27

		if val := content[edgeFlagOffset] & 32; val != 0 {
			log.Printf("Internal EdgeType & Bit5 %v", val)
		}
		if val := content[edgeFlagOffset] & 64; val != 0 {
			log.Printf("Internal EdgeType & Bit6 %v", val)
		}

		globalEdge := NavMeshGlobalEdge{
			Min: f64.Vec2{
				float64(utils.Float32FromByteArray(content[minVectorOffset : minVectorOffset+4])),
				float64(utils.Float32FromByteArray(content[minVectorOffset+4 : maxVectorOffset])),
			},
			Max: f64.Vec2{
				float64(utils.Float32FromByteArray(content[maxVectorOffset : maxVectorOffset+4])),
				float64(utils.Float32FromByteArray(content[maxVectorOffset+4 : edgeFlagOffset])),
			},
			EdgeFlag:        content[edgeFlagOffset],
			AssocDirection0: content[assocDir0Offset],
			AssocDirection1: content[assocDir1Offset],
			AssocCell0:      binary.LittleEndian.Uint16(content[assocCell0Offset:assocCell1Offset]),
			AssocCell1:      binary.LittleEndian.Uint16(content[assocCell1Offset:assocRgn0Offset]),
			AssocRegion0:    binary.LittleEndian.Uint16(content[assocRgn0Offset:assocRgn1Offset]),
			AssocRegion1:    binary.LittleEndian.Uint16(content[assocRgn1Offset : assocRgn1Offset+2]),
		}

		globalEdges.GlobalEdges = append(globalEdges.GlobalEdges, globalEdge)
	}
	return globalEdges
}
