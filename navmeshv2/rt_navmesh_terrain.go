package navmeshv2

import (
	"fmt"
	"github.com/g3n/engine/math32"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	return t.Cells[index]
}

func (t RtNavmeshTerrain) ResolveCellAndHeight(vPos *math32.Vector3) (RtNavmeshCell, error) {
	if vPos.X < 0 || vPos.X >= TerrainWidth || vPos.Z < 0 || vPos.Z >= TerrainHeight {
		return nil, errors.New(fmt.Sprintf("position %v not in cell", vPos))
	}

	if !t.TryFindHeight(vPos) {
		logrus.Panicf("failed to find height vor %v in region %d", vPos, t.Region.ID)
	}

	tile := t.GetTile(int(vPos.X/TerrainWidth), int(vPos.Z/TerrainHeight))
	return t.GetCell(tile.GetCellIndex()), nil
}

func (t RtNavmeshTerrain) ResolvePosition(pos RtNavmeshPosition) {
	inputHeight := pos.Offset.Y
	vTerrainTest := pos.Offset.Clone()
	vObjectTest := pos.Offset.Clone()

	cell, err := t.ResolveCellAndHeight(vTerrainTest)
	if err != nil {
		logrus.Panic(err)
	}
	quad, ok := cell.(RtNavmeshCellQuad)
	if !ok {
		return
	}

	pos.Cell = quad
	pos.Instance = nil

	deltaTerrain := vTerrainTest.Y - inputHeight
	deltaTerrainAbs := math32.Abs(deltaTerrain)

	for _, inst := range quad.Objects {
		vObjectTestNew, tri := inst.GetRtNavmeshCellTri(vObjectTest)
		if vObjectTestNew.Equals(vObjectTest) {
			continue
		}
		vObjectTest = vObjectTestNew

		deltaObj := vObjectTest.Y - inputHeight
		deltaObjAbs := math32.Abs(deltaObj)

		if deltaObjAbs < deltaTerrainAbs {
			deltaTerrain = deltaObj
			pos.Cell = tri
			pos.Instance = &inst
		}
	}

	if pos.Instance == nil {
		pos.Offset = vTerrainTest
		pos.Offset.Y = deltaTerrain + inputHeight

		tileX := int(pos.Offset.X / TileWidth)
		tileZ := int(pos.Offset.Z / TileWidth)
		tile := t.GetTile(tileX, tileZ)
		if tile.Flag.IsBlocked() {
			pos.Cell = nil
		}
	}
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

func (t RtNavmeshTerrain) TryFindHeight(vPos *math32.Vector3) bool {
	if vPos.X < 0 || vPos.X >= TerrainWidth || vPos.Z < 0 || vPos.Z >= TerrainHeight {
		return false
	}

	tileX := int(vPos.X / TerrainWidth)
	tileZ := int(vPos.Z / TerrainHeight)
	offsetX := (vPos.X - (float32(tileX) * TerrainWidth)) / TerrainWidth
	offsetZ := (vPos.Z - (float32(tileZ) * TerrainHeight)) / TerrainHeight

	// https://en.wikipedia.org/wiki/Bilinear_interpolation
	y1 := (1.0 - offsetX) * (((1.0 - offsetZ) * t.GetHeight(tileX+0, tileZ+0)) + (offsetZ * t.GetHeight(tileX+0, tileZ+1)))
	y2 := offsetX * (((1.0 - offsetZ) * t.GetHeight(tileX+1, tileZ+0)) + (offsetZ * t.GetHeight(tileX+1, tileZ+1)))
	vPos.Y = y1 + y2

	xBlock := tileX / (PlaneWidth / TileWidth)
	zBlock := tileZ / (PlaneHeight / TileHeight)
	plane := t.GetPlane(xBlock, zBlock)
	if plane.SurfaceType.IsIce() {
		vPos.Y = math32.Max(vPos.Y, plane.Height)
	}

	return true
}
