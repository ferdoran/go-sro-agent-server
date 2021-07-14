package navmesh

import (
	"encoding/binary"
)

const (
	MaxTilesX = 96
	MaxTilesY = 96
)

type TileMap struct {
	Map []Tile
}

type Tile struct {
	CellID       uint32
	Flag         uint16
	TextureIndex uint16
}

func loadTileMap(content []byte, readIndex *int) TileMap {
	tileMap := TileMap{Map: make([]Tile, 0)}
	for i := 0; i < MaxTilesX*MaxTilesY; i++ {
		cellIdOffset := *readIndex
		flagOffset := cellIdOffset + 4
		textureIndexOffset := flagOffset + 2

		*readIndex += 8

		tile := Tile{
			CellID:       binary.LittleEndian.Uint32(content[cellIdOffset:flagOffset]),
			Flag:         binary.LittleEndian.Uint16(content[flagOffset:textureIndexOffset]),
			TextureIndex: binary.LittleEndian.Uint16(content[textureIndexOffset : textureIndexOffset+2]),
		}

		tileMap.Map = append(tileMap.Map, tile)
	}
	return tileMap
}
