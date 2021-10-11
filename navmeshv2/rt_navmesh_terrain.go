package navmeshv2

import (
	"github.com/g3n/engine/math32"
)

const (
	BlocksX     = 6
	BlocksY     = 6
	BlocksTotal = BlocksX * BlocksY

	TilesX     = 96
	TilesY     = 96
	TilesTotal = TilesX * TilesY

	VerticesX     = TilesX + 1
	VerticesY     = TilesY + 1
	VerticesTotal = VerticesX * VerticesY

	TerrainWidth     = TilesX * TileWidth
	TerrainHeight    = TilesY * TileHeight
	TerrainWidthInt  = 1920
	TerrainHeightInt = 1920
)

type RtNavmeshTerrain struct {
	RtNavmeshBase
	Region Region

	tileMap   [TilesTotal]RtNavmeshTile
	planeMap  [BlocksTotal]RtNavmeshPlane
	heightMap [VerticesTotal]float32

	Objects []RtNavmeshInstObj
	Cells   []RtNavmeshCellQuad

	GlobalEdges   []RtNavmeshEdgeGlobal
	InternalEdges []RtNavmeshEdgeInternal
}

func (t RtNavmeshTerrain) GetNavmeshType() RtNavmeshType {
	return RtNavmeshTypeTerrain
}

func NewRtNavmeshTerrain(filename string, region Region) RtNavmeshTerrain {
	return RtNavmeshTerrain{
		RtNavmeshBase: RtNavmeshBase{Filename: filename},
		Region:        region,
		Objects:       make([]RtNavmeshInstObj, 0),
		Cells:         make([]RtNavmeshCellQuad, 0),
		GlobalEdges:   make([]RtNavmeshEdgeGlobal, 0),
		InternalEdges: make([]RtNavmeshEdgeInternal, 0),
	}
}

func (t RtNavmeshTerrain) GetCell(index int) RtNavmeshCell {
	return &t.Cells[index]
}

func (t RtNavmeshTerrain) GetTile(x, y int) RtNavmeshTile {
	return t.tileMap[y*TilesY+x]
}

func (t RtNavmeshTerrain) GetHeight(x, y int) float32 {
	return t.heightMap[y*VerticesY+x]
}

func (t RtNavmeshTerrain) GetPlane(xBlock, zBlock int) RtNavmeshPlane {
	return t.planeMap[zBlock*BlocksY+xBlock]
}

func (t RtNavmeshTerrain) ResolveCell(pos *math32.Vector3) (RtNavmeshCellQuad, error) {
	tile := t.GetTile(int(pos.X/TileWidth), int(pos.Z/TileHeight))
	return t.Cells[tile.CellIndex], nil
}

func (t RtNavmeshTerrain) ResolveHeight(pos *math32.Vector3) float32 {
	tileX := int(pos.X / TileWidth)
	tileZ := int(pos.Z / TileHeight)

	h1 := t.GetHeight(tileX, tileZ)
	h2 := t.GetHeight(tileX, tileZ+1)
	h3 := t.GetHeight(tileX+1, tileZ)
	h4 := t.GetHeight(tileX+1, tileZ+1)

	// h1--------h3
	// |   |      |
	// |   |      |
	// h5--+------h6
	// |   |      |
	// h2--------h4

	tileOffsetX := pos.X - (TileWidth * float32(tileX))
	tileOffsetXLength := tileOffsetX / TileWidth
	tileOffsetZ := pos.Z - (TileHeight * float32(tileZ))
	tileOffsetZLength := tileOffsetZ / TileHeight

	h5 := h1 + (h2-h1)*tileOffsetZLength
	h6 := h3 + (h4-h3)*tileOffsetZLength
	yHeight := h5 + (h6-h5)*tileOffsetXLength

	return yHeight
}
