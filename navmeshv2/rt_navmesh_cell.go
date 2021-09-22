package navmeshv2

type RtNavmeshCell interface {
	GetIndex() int
	GetMesh() RtNavmesh
	AddEdge(edge RtNavmeshEdge, direction RtNavmeshEdgeDirection)
}

type RtNavmeshCellBase struct {
	Index int
	Mesh  RtNavmesh
}

func (r RtNavmeshCellBase) GetIndex() int {
	return r.Index
}

func (r RtNavmeshCellBase) GetMesh() RtNavmesh {
	return r.Mesh
}

type RtNavmeshCellQuad struct {
	RtNavmeshCellBase
	edges   []RtNavmeshEdge
	Rect    Rectangle
	Objects []RtNavmeshInstObj
}

func (r RtNavmeshCellQuad) AddEdge(edge RtNavmeshEdge, direction RtNavmeshEdgeDirection) {
	panic("implement me")
}

type RtNavmeshCellTri struct {
	RtNavmeshCellBase
	edges    []RtNavmeshEdge
	Triangle Triangle
	Flag     int16
}

func (r RtNavmeshCellTri) AddEdge(edge RtNavmeshEdge, direction RtNavmeshEdgeDirection) {
	panic("implement me")
}
