package model

import (
	"fmt"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/utils"
)

type Position struct {
	X       float32
	Y       float32
	Z       float32
	Heading float32
	Region  *Region
}

func (p Position) String() string {
	return fmt.Sprintf("(X=%f, Y=%f, Z=%f, H=%d, R=%d)", p.X, p.Y, p.Z, p.Heading, p.Region.ID)
}

func (p Position) ToWorldCoordinates() (float32, float32, float32) {
	rX, rZ := utils.Int16ToXAndZ(p.Region.ID)

	if p.X > RegionWidth {
		rX++
		p.X -= RegionWidth
	}

	if p.Z > RegionHeight {
		rZ++
		p.Z -= RegionHeight
	}

	regionId := utils.XAndZToInt16(byte(rX), byte(rZ))
	region := sroWorldInstance.Regions[regionId]

	p.Region = region
	x := (float32(rX) * 1920.0) + p.X
	y := p.Region.GetYAtOffset(p.X, p.Z)
	z := (float32(rZ) * 1920.0) + p.Z
	return x, y, z
}

func (p Position) ToWorldCoordinatesInt32() (int32, int32, int32) {
	rX, rZ := utils.Int16ToXAndZ(p.Region.ID)

	x := (float32(rX) * 1920.0) + p.X
	y := p.Region.GetYAtOffset(p.X, p.Z)
	z := (float32(rZ) * 1920.0) + p.Z
	return int32(x), int32(y), int32(z)
}

func (p Position) DistanceToSquared(other Position) float32 {
	x1, y1, z1 := p.ToWorldCoordinates()
	x2, y2, z2 := other.ToWorldCoordinates()

	dx := x1 - x2
	dy := y1 - y2
	dz := z1 - z2

	return dx*dx + dy*dy + dz*dz
}

func (p Position) DistanceTo(other Position) float32 {
	return math32.Sqrt(p.DistanceToSquared(other))
}

func (p *Position) GetCell() *TerrainCell {
	return p.Region.GetCellAtOffset(p.X, p.Z)
}

func NewPosFromWorldCoordinates(x, z float32) Position {
	regionX := byte(x / RegionWidth)
	regionZ := byte(z / RegionHeight)

	regionId := utils.XAndZToInt16(regionX, regionZ)
	region := sroWorldInstance.Regions[regionId]
	pX := x - float32(regionX)*float32(RegionWidth)
	pZ := z - float32(regionZ)*float32(RegionHeight)
	pY := region.GetYAtOffset(pX, pZ)

	return Position{
		X:       pX,
		Y:       pY,
		Z:       pZ,
		Heading: 0,
		Region:  region,
	}
}

func (p *Position) nextPositionIsSlope(nextPos Position) bool {

	x1, y1, z1 := p.ToWorldCoordinates()
	x2, y2, z2 := nextPos.ToWorldCoordinates()
	dx := x1 - x2
	dy := y1 - y2
	dz := z1 - z2

	nonYDistance := math32.Sqrt(dx*dx + dz*dz)

	angle := dy / nonYDistance
	logrus.Tracef("slope angle is %f\n", angle)
	return angle >= 1 || angle <= -1
}
