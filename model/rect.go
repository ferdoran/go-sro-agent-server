package model

import (
	"github.com/g3n/engine/math32"
)

type Rectangle struct {
	Min math32.Vector2
	Max math32.Vector2
}

func (r *Rectangle) Width() float32 {
	return math32.Abs(r.Max.X - r.Min.X)
}

func (r *Rectangle) Height() float32 {
	return math32.Abs(r.Max.Y - r.Min.Y)
}
