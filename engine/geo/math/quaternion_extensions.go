package math

import "github.com/g3n/engine/math32"

func NewQuaternion(yaw, pitch, roll float32) *math32.Quaternion {
	// In SRO y determines the height and not z
	// Therefore we have to switch y and z
	q := math32.NewQuaternion(0, 0, 0, 0)
	q = q.SetFromEuler(math32.NewVector3(0, yaw, 0))
	return q
}
