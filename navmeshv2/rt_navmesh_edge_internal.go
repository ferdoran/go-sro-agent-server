package navmeshv2

type RtNavmeshEdgeInternal struct {
	RtNavmeshEdgeBase
}

func (e *RtNavmeshEdgeInternal) Link() {
	e.SrcCell = e.GetMesh().GetCell(e.SrcCellIndex)
	e.SrcCell.AddEdge(e, e.SrcDirection)

	if !e.GetFlag().IsBlocked() {

		e.DstCell = e.GetMesh().GetCell(e.DstCellIndex)
		e.DstCell.AddEdge(e, e.DstDirection)
	}
}
