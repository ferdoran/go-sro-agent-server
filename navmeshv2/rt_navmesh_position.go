package navmeshv2

import "github.com/g3n/engine/math32"

type RtNavmeshPosition struct {
	Cell     RtNavmeshCell
	Instance RtNavmeshInst
	Region   Region
	Offset   *math32.Vector3
}

func (pos RtNavmeshPosition) GetMesh() RtNavmesh {
	return pos.Cell.GetMesh()
}
