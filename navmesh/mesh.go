package navmesh

import "github.com/g3n/engine/math32"

type Mesh struct {
	TriangleCells []bool
	GlobalEdges   []bool
	InternalEdges []bool
	Events        []string
	Grid          ObjectGrid
}

type RTNavMeshCellTri struct {
	RTNavMeshCell
	Triangle       math32.Triangle
	RTNavMeshEdges [3]RTNavMeshEdge
	Flag           uint16
}

type RTNavMeshCell struct {
	RTNavMesh
	EdgeCount int
	Index     int
}

type RTNavMeshCellQuad struct {
	RTNavMeshCell
	math32.Box2
	Edges   []RTNavMeshEdge
	Objects []RTNavMeshInstanceObject
}

type RTNavMesh struct {
	FileName      string
	RTNavMeshType byte
}

type RTNavMeshObj struct {
	Cells       []RTNavMeshCellTri
	GlobalEdges []RTNavMeshEdgeGlobal
}

type RTNavMeshInstance struct {
	RTNavMesh
	RTNavMeshObj
	ID           int
	Position     math32.Vector3
	Rotation     math32.Quaternion
	Scale        math32.Vector3
	LocalToWorld math32.Matrix4
	WorldToLocal math32.Matrix4
}

type RTNavMeshInstanceObject struct {
	RTNavMeshInstance
	WorldID int
	Region  int
}

type RTNavMeshEdge struct {
	RTNavMesh
	Index                int
	Flag                 byte
	Line                 math32.Line3
	SourceDirection      int8
	DestinationCell      RTNavMeshCell
	DestinationCellIndex uint16
}

type RTNavMeshEdgeGlobal struct {
	RTNavMeshEdge
	SourceMeshIndex      uint16
	DestinationMeshIndex uint16
}

type RTNavMeshEdgeInternal struct {
	RTNavMeshEdge
}
