package model

import (
	"fmt"
	"github.com/ferdoran/go-sro-framework/utils"
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

func (p Position) ToWorldCoordinatesInt32() (int32, int32, int32) {
	rX, rZ := utils.Int16ToXAndZ(p.Region.ID)

	x := (float32(rX) * 1920.0) + p.X
	y := p.Region.GetYAtOffset(p.X, p.Z)
	z := (float32(rZ) * 1920.0) + p.Z
	return int32(x), int32(y), int32(z)
}

func (p *Position) GetCell() *TerrainCell {
	return p.Region.GetCellAtOffset(p.X, p.Z)
}
