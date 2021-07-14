package navmesh

import (
	"github.com/ferdoran/go-sro-framework/math"
	"github.com/g3n/engine/math32"
)

type ObjectGlobalEdge struct {
	A                    *math32.Vector3
	B                    *math32.Vector3
	SourceCellIndex      int
	DestinationCellIndex int
	SourceDirection      int
	DestinationDirection int
	SourceMeshIndex      int
	DestinationMeshIndex int
	Flag                 byte
	EventZoneFlag        byte
}

type ObjectInternalEdge struct {
	A                    *math32.Vector3
	B                    *math32.Vector3
	SourceCellIndex      int
	DestinationCellIndex int
	SourceDirection      int
	DestinationDirection int
	Flag                 byte
	EventZoneFlag        byte
}

func (e *ObjectInternalEdge) ToLine2() *math.Line2 {
	return math.NewLine2(
		math32.NewVector2(e.A.X, e.A.Z),
		math32.NewVector2(e.B.X, e.B.Z),
	)
}

func (e *ObjectGlobalEdge) ToLine2() *math.Line2 {
	return math.NewLine2(
		math32.NewVector2(e.A.X, e.A.Z),
		math32.NewVector2(e.B.X, e.B.Z),
	)
}
