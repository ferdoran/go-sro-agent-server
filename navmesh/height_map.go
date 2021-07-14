package navmesh

import (
	"github.com/ferdoran/go-sro-framework/utils"
)

type HeightMap struct {
	Heights []float32
}

func loadHeightMap(content []byte, readIndex *int) HeightMap {
	heightMap := HeightMap{
		Heights: make([]float32, 0),
	}

	for i := 0; i < 97*97; i++ {
		heightMap.Heights = append(heightMap.Heights, utils.Float32FromByteArray(content[*readIndex:*readIndex+4]))
		*readIndex += 4
	}

	return heightMap
}
