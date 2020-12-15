package geo

import "github.com/g3n/engine/math32"

type Collision struct {
	EdgeFlag     byte
	VectorGlobal *math32.Vector3
	VectorLocal  *math32.Vector3
}

func (c Collision) IsBridge() bool {
	return c.EdgeFlag == 16 || c.EdgeFlag == 20
}

func (c Collision) Equals(other Collision) bool {
	return c.EdgeFlag == other.EdgeFlag && c.VectorGlobal.Equals(other.VectorGlobal) && c.VectorLocal.Equals(other.VectorLocal)
}
