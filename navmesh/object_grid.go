package navmesh

import (
	"github.com/ferdoran/go-sro-framework/math"
	"github.com/g3n/engine/math32"
)

type ObjectGrid struct {
	Origin    *math32.Vector3
	Width     int
	Height    int
	GridTiles []*ObjectTile
}

func (og *ObjectGrid) ContainsPoint(vPos *math32.Vector3) bool {

	minVec := math32.NewVector2(og.Origin.X, og.Origin.Z)
	maxVec := math32.NewVector2(minVec.X+float32(og.Width)*100, minVec.Y+float32(og.Height)*100)

	return vPos.X >= minVec.X && vPos.X <= maxVec.X && vPos.Z >= minVec.Y && vPos.Z <= maxVec.Y
}

func (og *ObjectGrid) Rect() *math.Rectangle {
	minVec := math32.NewVector2(og.Origin.X, og.Origin.Z)
	maxVec := math32.NewVector2(minVec.X+float32(og.Width)*100, minVec.Y+float32(og.Height)*100)

	rect := math.Rectangle{
		Min: minVec,
		Max: maxVec,
	}

	return &rect
}
