package navmeshv2

import "github.com/g3n/engine/math32"

type Rectangle struct {
	Min *math32.Vector2
	Max *math32.Vector2
}

func NewRectangle(x, y, width, height float32) Rectangle {
	min := math32.NewVector2(x, y)
	max := math32.NewVector2(x+width, y+height)

	return Rectangle{
		Min: min,
		Max: max,
	}
}

func (r Rectangle) X() float32 {
	return r.Min.X
}

func (r Rectangle) Y() float32 {
	return r.Min.Y
}

func (r Rectangle) Width() float32 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Height() float32 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Center() *math32.Vector2 {
	min := math32.NewVector2(r.Min.X, r.Min.Y)
	max := math32.NewVector2(r.Max.X, r.Max.Y)

	return min.Add(max).MultiplyScalar(0.5)
}

func (r Rectangle) Contains(x, y float32) bool {
	return x >= r.Min.X && x <= r.Max.X && y >= r.Min.Y && y <= r.Max.Y
}

func (r Rectangle) ContainsVec2(p *math32.Vector2) bool {
	return r.Contains(p.X, p.Y)
}

func (r Rectangle) ContainsVec3(p *math32.Vector3) bool {
	return r.Contains(p.X, p.Z)
}

func (r Rectangle) IntersectsTriangle(triangle Triangle) bool {
	triV1 := math32.NewVector2(triangle.A.X, triangle.A.Z)
	triV2 := math32.NewVector2(triangle.B.X, triangle.B.Z)
	triV3 := math32.NewVector2(triangle.C.X, triangle.C.Z)

	if r.ContainsVec2(triV1) || r.ContainsVec2(triV2) || r.ContainsVec2(triV3) {
		return true
	}

	recV1 := math32.NewVector3(r.Min.X, 0, r.Min.Y)
	recV2 := math32.NewVector3(r.Min.X, 0, r.Max.Y)
	recV3 := math32.NewVector3(r.Max.X, 0, r.Min.Y)
	recV4 := math32.NewVector3(r.Max.X, 0, r.Max.Y)

	h1, _ := triangle.FindHeight(recV1)
	h2, _ := triangle.FindHeight(recV2)
	h3, _ := triangle.FindHeight(recV3)
	h4, _ := triangle.FindHeight(recV4)

	if h1 || h2 || h3 || h4 {
		return true
	}

	line1 := NewLineSegmentVec2(triV1, triV2)
	line2 := NewLineSegmentVec2(triV1, triV3)
	line3 := NewLineSegmentVec2(triV2, triV3)

	return r.IntersectsLine(line1) || r.IntersectsLine(line2) || r.IntersectsLine(line3)
}

func (r Rectangle) IntersectsLine(line LineSegment) bool {
	if r.ContainsVec3(line.A) || r.ContainsVec3(line.B) {
		return true
	}

	recV1 := math32.NewVector3(r.Min.X, 0, r.Min.Y)
	recV2 := math32.NewVector3(r.Min.X, 0, r.Max.Y)
	recV3 := math32.NewVector3(r.Max.X, 0, r.Min.Y)
	recV4 := math32.NewVector3(r.Max.X, 0, r.Max.Y)

	l1 := LineSegment{recV1, recV2}
	l2 := LineSegment{recV1, recV3}
	l3 := LineSegment{recV3, recV4}
	l4 := LineSegment{recV2, recV4}

	if i1, _ := line.Intersects(l1); i1 {
		return i1
	} else if i2, _ := line.Intersects(l2); i2 {
		return i2
	} else if i3, _ := line.Intersects(l3); i3 {
		return i3
	} else if i4, _ := line.Intersects(l4); i4 {
		return i4
	}

	return false
}
