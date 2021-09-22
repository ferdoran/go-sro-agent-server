package navmeshv2

import "github.com/g3n/engine/math32"

type Triangle struct {
	A *math32.Vector3
	B *math32.Vector3
	C *math32.Vector3
}

func (t Triangle) Center() *math32.Vector3 {
	a := t.A.Clone()
	b := t.B.Clone()
	c := t.C.Clone()

	return a.Add(b).Add(c).DivideScalar(3)
}

func (t Triangle) FindHeight(pos *math32.Vector3) bool {
	denominator := ((t.B.Z - t.C.Z) * (t.A.X - t.C.X)) + ((t.C.X - t.B.X) * (t.A.Z - t.C.Z))

	a := (((t.B.Z - t.C.Z) * (pos.X - t.C.X)) + ((t.C.X - t.B.X) * (pos.Z - t.C.Z))) / denominator
	b := (((t.C.Z - t.A.Z) * (pos.X - t.C.X)) + ((t.A.X - t.C.X) * (pos.Z - t.C.Z))) / denominator
	c := 1 - a - b

	pos.Y = ((a * t.A.Y) + (b * t.B.Y) + (c * t.C.Y)) / (a + b + c)
	// return a > 0 && a < 1 && b > 0 && b < 1 && c > 0 && c < 1 // point can only be within triangle
	return a >= 0 && a <= 1 && b >= 0 && b <= 1 && c >= 0 && c <= 1 // point can be on border
}

func (t Triangle) OffsetTowardsCenter(pos *math32.Vector2) {
	var tolerance float32 = 0.199999988079071
	center := t.Center()
	centerV2 := math32.NewVector2(center.X, center.Z)
	vDelta := centerV2.Sub(pos)

	if vDelta.Length() > tolerance {
		vDelta = vDelta.Normalize().MultiplyScalar(tolerance)
	}

	pos = pos.Add(vDelta)
}
