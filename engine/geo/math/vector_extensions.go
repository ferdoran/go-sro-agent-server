package math

import "github.com/g3n/engine/math32"

func Determinant(v1, v2 *math32.Vector2) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Calculates the angle between two vectors with regard to East in degrees
func AngleToEastInDeg(diff math32.Vector3) float32 {
	normalizedDiff := *diff.Normalize()
	v1 := math32.NewVector2(-1, 0)
	v2 := math32.NewVector2(normalizedDiff.X, normalizedDiff.Z)
	dot := v1.Dot(v2)
	det := Determinant(v1, v2)
	return math32.RadToDeg(math32.Atan2(det, dot))
}

func SlopeAngleInDeg(diff *math32.Vector3) float32 {
	normalizedDiff := *diff.Normalize()
	v1 := math32.NewVector2(1, 0)
	v2 := math32.NewVector2(normalizedDiff.X, normalizedDiff.Z)
	dot := v1.Dot(v2)
	det := Determinant(v1, v2)
	return math32.RadToDeg(math32.Atan2(det, dot))
}
