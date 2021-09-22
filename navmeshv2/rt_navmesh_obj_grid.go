package navmeshv2

import (
	"github.com/ferdoran/go-sro-agent-server/navmesh"
	"github.com/g3n/engine/math32"
)

type RtNavmeshObjGrid struct {
	object    *RtNavmeshObj
	Tiles     []RtNavmeshObjGridTile
	X         float32
	Y         float32
	Width     int
	Height    int
	Rectangle Rectangle
}

func NewRtNavmeshObjGrid(object *RtNavmeshObj) RtNavmeshObjGrid {
	return RtNavmeshObjGrid{
		object: object,
		Tiles:  make([]RtNavmeshObjGridTile, 0),
	}
}

func (grid *RtNavmeshObjGrid) AddTile(tile RtNavmeshObjGridTile) {
	grid.Tiles = append(grid.Tiles, tile)
}

func (grid *RtNavmeshObjGrid) GetTile(index int) RtNavmeshObjGridTile {
	return grid.Tiles[index]
}

func (grid *RtNavmeshObjGrid) GetTileFromXAndY(x, y int) RtNavmeshObjGridTile {
	return grid.Tiles[y*grid.Width+x]
}

func (grid *RtNavmeshObjGrid) GetTileFromVec(v *math32.Vector3) RtNavmeshObjGridTile {
	tileX := int((v.X - grid.X) / RtNavmeshObjGridTileWidth)
	tileY := int((v.Z - grid.Y) / RtNavmeshObjGridTileHeight)

	return grid.GetTileFromXAndY(tileX, tileY)
}

func (grid *RtNavmeshObjGrid) GetTileFromVec2(v *math32.Vector2) RtNavmeshObjGridTile {
	tileX := int((v.X - grid.X) / RtNavmeshObjGridTileWidth)
	tileY := int((v.Y - grid.Y) / RtNavmeshObjGridTileHeight)

	return grid.GetTileFromXAndY(tileX, tileY)
}

func (grid *RtNavmeshObjGrid) Load(reader navmesh.Loader) {
	// TODO implement
	panic("implement me")
}
