package navmeshv2

import (
	"github.com/ferdoran/go-sro-agent-server/navmesh"
)

type RtNavmeshInstObj struct {
	RtNavmeshInstBase
	Region  Region
	WorldID int
}

func (inst *RtNavmeshInstObj) Read(reader *navmesh.Loader) {
	// TODO Implement
	panic("implement me")
}

//func NewRtNavmeshInstObj(mesh RtNavmesh) *RtNavmeshInstObj {
//	// TODO fill all fields?
//	return &RtNavmeshInstObj{
//		RtNavmeshInstBase: RtNavmeshInstBase{
//			Mesh:         mesh,
//			Object:       nil,
//			ID:           0,
//			Position:     nil,
//			Rotation:     nil,
//			Scale:        nil,
//			LocalToWorld: nil,
//			WorldToLocal: nil,
//		},
//		Region:            nil,
//		WorldID:           0,
//	}
//}
