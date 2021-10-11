package navmeshv2

import "github.com/g3n/engine/math32"

type LineSegment struct {
	A *math32.Vector3
	B *math32.Vector3
}

func NewLineSegmentVec2(a, b *math32.Vector2) LineSegment {
	return LineSegment{
		A: math32.NewVector3(a.X, 0, a.Y),
		B: math32.NewVector3(b.X, 0, b.Y),
	}
}

func (ls LineSegment) Length() *math32.Vector3 {
	a := ls.A.Clone()
	b := ls.B.Clone()

	return b.Sub(a)
}

func (ls LineSegment) Center() *math32.Vector3 {
	a := ls.A.Clone()
	b := ls.B.Clone()

	return a.Add(b).MultiplyScalar(0.5)
}

func (ls LineSegment) GetPointRelation(p *math32.Vector2) float32 {
	return ((ls.A.X - p.X) * (ls.B.Z - p.Y)) - ((ls.A.Z - p.Y) * (ls.B.X - p.X))
}

func (ls LineSegment) Intersects(other LineSegment) (bool, *math32.Vector2) {

	denominator := ((other.B.Z - other.A.Z) * (ls.B.X - ls.A.X)) - ((other.B.X - other.A.X) * (ls.B.Z - ls.A.Z))

	if denominator != 0 {
		uA := (((other.B.X - other.A.X) * (ls.A.Z - other.A.Z)) - ((other.B.Z - other.A.Z) * (ls.A.X - other.A.X))) / denominator
		uB := (((ls.B.X - ls.A.X) * (ls.A.Z - other.A.Z)) - ((ls.B.Z - ls.A.Z) * (ls.A.X - other.A.X))) / denominator

		//if (uA > 0f && uA < 1f && uB > 0f && uB < 1f) // exclusive caps
		if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 { // inclusive caps
			a := ls.A.Clone()
			b := ls.B.Clone()
			result := a.Add(b.Sub(a.Clone()).MultiplyScalar(uA))
			return true, math32.NewVector2(result.X, result.Z)
		}
	}

	return false, math32.NewVec2()
}

func (ls LineSegment) Intersects3D(other LineSegment) (bool, *math32.Vector3) {

	denominator := ((other.B.Z - other.A.Z) * (ls.B.X - ls.A.X)) - ((other.B.X - other.A.X) * (ls.B.Z - ls.A.Z))

	if denominator != 0 {
		uA := (((other.B.X - other.A.X) * (ls.A.Z - other.A.Z)) - ((other.B.Z - other.A.Z) * (ls.A.X - other.A.X))) / denominator
		uB := (((ls.B.X - ls.A.X) * (ls.A.Z - other.A.Z)) - ((ls.B.Z - ls.A.Z) * (ls.A.X - other.A.X))) / denominator

		//if (uA > 0f && uA < 1f && uB > 0f && uB < 1f) // exclusive caps
		if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 { // inclusive caps
			a := ls.A.Clone()
			b := ls.B.Clone()
			result := a.Add(b.Sub(a.Clone()).MultiplyScalar(uA))

			return true, math32.NewVector3(result.X, result.Y, result.Z)
		}
	}

	return false, math32.NewVec3()
}
