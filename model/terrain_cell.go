package model

import "github.com/g3n/engine/math32"

type TerrainCell struct {
	*math32.Box2
	Min              *math32.Vector2
	Max              *math32.Vector2
	ID               int
	ObjCount         byte
	Objects          []uint16
	LinkedNeighbours []*TerrainCell
	RegionLinks      map[int16][]*TerrainCell
	RegionID         int16
}
