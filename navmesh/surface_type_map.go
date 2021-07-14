package navmesh

import (
	"github.com/ferdoran/go-sro-framework/utils"
)

const (
	BlocksX = 6
	BlocksY = 6
)

type SurfaceTypeMap struct {
	SurfaceTypes []byte // 0 = Solid, 1 = Water, 2 = Ice
}

type SurfaceHeightMap struct {
	SurfaceHeights []float32
}

func loadSurfaceMaps(content []byte, readIndex *int) (SurfaceTypeMap, SurfaceHeightMap) {
	surfaceTypeMap := SurfaceTypeMap{SurfaceTypes: content[*readIndex : *readIndex+36]}
	*readIndex += 36

	surfaceHeightMap := SurfaceHeightMap{SurfaceHeights: make([]float32, 36)}

	for i := 0; i < 6*6; i++ {
		surfaceHeightMap.SurfaceHeights = append(surfaceHeightMap.SurfaceHeights, utils.Float32FromByteArray(content[*readIndex:*readIndex+4]))
		*readIndex += 4
	}

	return surfaceTypeMap, surfaceHeightMap

}

//enum EdgeDirection
//{
//	Invalid = -1,
//	North = 0,
//	East = 1,
//	South = 2,
//	West = 3,
//}
