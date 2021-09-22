package navmeshv2

const (
	RtNavmeshObjGridTileWidth  = 100.0
	RtNavmeshObjGridTileHeight = 100.0
)

type RtNavmeshObjGridTile struct {
	Index         int
	X             int
	Y             int
	Rectangle     Rectangle
	grid          RtNavmeshObjGrid
	GlobalEdges   []RtNavmeshEdgeGlobal
	InternalEdges []RtNavmeshEdgeInternal
	Cells         []RtNavmeshCellTri
}

func (tile *RtNavmeshObjGridTile) AddGlobalEdge(edge RtNavmeshEdgeGlobal) {
	tile.GlobalEdges = append(tile.GlobalEdges, edge)
}

func (tile *RtNavmeshObjGridTile) AddInternalEdge(edge RtNavmeshEdgeInternal) {
	tile.InternalEdges = append(tile.InternalEdges, edge)
}

func (tile *RtNavmeshObjGridTile) AddCell(cell RtNavmeshCellTri) {
	tile.Cells = append(tile.Cells, cell)
}

func (tile *RtNavmeshObjGridTile) GetCells() []RtNavmeshCellTri {
	return tile.Cells
}
