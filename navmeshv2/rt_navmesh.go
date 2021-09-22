package navmeshv2

import "github.com/g3n/engine/math32"

type RtNavmesh interface {
	GetNavmeshType() RtNavmeshType
	GetFilename() string
	GetCell(index int) RtNavmeshCell
	ResolveCellAndHeight(vPos *math32.Vector3) (RtNavmeshCell, error)
	ResolvePosition(pos RtNavmeshPosition)
}

type RtNavmeshBase struct {
	Filename string
}

func (base RtNavmeshBase) GetFilename() string {
	return base.Filename
}
