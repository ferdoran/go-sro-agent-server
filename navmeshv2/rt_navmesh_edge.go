package navmeshv2

type RtNavmeshEdge interface {
	GetType() RtNavmeshEdgeMeshType
	GetMesh() RtNavmesh
	GetIndex() int
	GetLine() LineSegment
	GetFlag() RtNavmeshEdgeFlag
	GetSrcDirection() RtNavmeshEdgeDirection
	GetDstDirection() RtNavmeshEdgeDirection
	GetSrcCellIndex() int
	GetDstCellIndex() int
	GetSrcCell() RtNavmeshCell
	GetDstCell() RtNavmeshCell
}

/*
 type: NavmeshEdgeType
    meshType: NavmeshEdgeMeshType
    sourceCellIndex: short
    destinationCellIndex: short
    sourceDirection: byte // not set for objects
    destinationDirection: byte // not set for objects
    sourceMeshIndex: short // not set for objects
    destinationMeshIndex: short // not set for objects
    flag: byte
*/

type RtNavmeshEdgeBase struct {
	RtNavmeshEdgeMeshType
	Mesh         RtNavmesh
	Index        int
	Line         LineSegment
	Flag         RtNavmeshEdgeFlag
	SrcDirection RtNavmeshEdgeDirection
	DstDirection RtNavmeshEdgeDirection
	SrcCellIndex int
	DstCellIndex int
	SrcCell      RtNavmeshCell
	DstCell      RtNavmeshCell
}

func (e RtNavmeshEdgeBase) GetType() RtNavmeshEdgeMeshType {
	return e.RtNavmeshEdgeMeshType
}

func (e RtNavmeshEdgeBase) GetMesh() RtNavmesh {
	return e.Mesh
}

func (e RtNavmeshEdgeBase) GetIndex() int {
	return e.Index
}

func (e RtNavmeshEdgeBase) GetLine() LineSegment {
	return e.Line
}

func (e RtNavmeshEdgeBase) GetFlag() RtNavmeshEdgeFlag {
	return e.Flag
}

func (e RtNavmeshEdgeBase) GetSrcDirection() RtNavmeshEdgeDirection {
	return e.SrcDirection
}

func (e RtNavmeshEdgeBase) GetDstDirection() RtNavmeshEdgeDirection {
	return e.DstDirection
}

func (e RtNavmeshEdgeBase) GetSrcCellIndex() int {
	return e.SrcCellIndex
}

func (e RtNavmeshEdgeBase) GetDstCellIndex() int {
	return e.DstCellIndex
}

func (e RtNavmeshEdgeBase) GetSrcCell() RtNavmeshCell {
	return e.SrcCell
}

func (e RtNavmeshEdgeBase) GetDstCell() RtNavmeshCell {
	return e.DstCell
}

func (e RtNavmeshEdgeBase) IsGlobalLinker() bool {
	return e.Flag.IsGlobal()
}

func (e RtNavmeshEdgeBase) IsLocalLinker() bool {
	return e.Flag.IsInternal()
}

func (e RtNavmeshEdgeBase) HasCellNeighbour() bool {
	return !(e.SrcCell == nil || e.DstCell == nil)
}

func (e RtNavmeshEdgeBase) IsBlocked(cell RtNavmeshCell) bool {
	if e.SrcCell == cell {
		return e.Flag.IsBlockedSrcToDst()
	} else {
		return e.Flag.IsBlockedDstToSrc()
	}
}

func (e RtNavmeshEdgeBase) GetRtNavmeshCell(index int) RtNavmeshCell {
	switch index {
	case 0:
		return e.SrcCell
	case 1:
		return e.DstCell
	default:
		return nil
	}
}
