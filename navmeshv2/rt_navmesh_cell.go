package navmeshv2

import "github.com/sirupsen/logrus"

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
	Edges   []RtNavmeshEdge
	Rect    Rectangle
	Objects []RtNavmeshInstObj
}

func (r *RtNavmeshCellQuad) AddEdge(edge RtNavmeshEdge, direction RtNavmeshEdgeDirection) {
	if r.Edges == nil {
		r.Edges = make([]RtNavmeshEdge, 0)
	}
	if direction.IsNone() {
		logrus.Errorf("trying to add edge without direction")
		return
	}

	r.Edges = append(r.Edges, edge)
}

type RtNavmeshCellTri struct {
	RtNavmeshCellBase
	edges    []RtNavmeshEdge
	Triangle Triangle
	Flag     int16
}

func (r *RtNavmeshCellTri) AddEdge(edge RtNavmeshEdge, direction RtNavmeshEdgeDirection) {
	if len(r.edges) < 3 {
		r.edges = append(r.edges, edge)
	}
}
