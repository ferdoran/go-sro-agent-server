package navmeshv2

type RtNavmeshEdgeGlobal struct {
	RtNavmeshEdgeBase
	SrcMeshIndex int
	DstMeshIndex int
}

func (e *RtNavmeshEdgeGlobal) GetMesh() RtNavmesh {
	panic("implement me")
}

func (e *RtNavmeshEdgeGlobal) Link(terrainIndex map[uint16]RtNavmeshTerrain) {
	e.SrcCell = e.Mesh.GetCell(e.SrcCellIndex)
	e.SrcCell.AddEdge(e, e.SrcDirection)

	if !e.Flag.IsBlocked() && e.DstMeshIndex > 0 {
		pDstNavMesh, exists := terrainIndex[uint16(e.DstMeshIndex)]
		if !exists {
			return
		}

		e.DstCell = pDstNavMesh.GetCell(e.DstCellIndex)
		e.DstCell.AddEdge(e, e.DstDirection)
	}
}
