package navmeshv2

type RtNavmeshPlane struct {
	Height      float32
	SurfaceType RtNavmeshSurfaceType
}

const (
	PlaneWidth               = 320.0
	PlaneHeight              = 320.0
	RtNavmeshSurfaceTypeNone = iota
	RtNavmeshSurfaceTypeWater
	RtNavmeshSurfaceTypeIce
)

type RtNavmeshSurfaceType byte

func (st RtNavmeshSurfaceType) IsNone() bool {
	return st == RtNavmeshSurfaceTypeNone
}

func (st RtNavmeshSurfaceType) IsWater() bool {
	return st == RtNavmeshSurfaceTypeWater
}

func (st RtNavmeshSurfaceType) IsIce() bool {
	return st == RtNavmeshSurfaceTypeIce
}
