package navmeshv2

import (
	"github.com/g3n/engine/math32"
	"github.com/pkg/errors"
)

type RtNavmeshObj struct {
	RtNavmeshBase
	Cells         []RtNavmeshCellTri
	GlobalEdges   []RtNavmeshEdgeGlobal
	InternalEdges []RtNavmeshEdgeInternal
	Events        []string
	Grid          RtNavmeshObjGrid
}

func NewNavmeshObj(filename string) RtNavmeshObj {
	obj := RtNavmeshObj{
		RtNavmeshBase: RtNavmeshBase{Filename: filename},
	}
	obj.Grid = NewRtNavmeshObjGrid(&obj)
	return obj
}

func (obj RtNavmeshObj) GetNavmeshType() RtNavmeshType {
	return RtNavmeshTypeObject
}

func (obj RtNavmeshObj) GetCell(index int) RtNavmeshCell {
	return obj.Cells[index]
}

func (obj RtNavmeshObj) ResolveCellAndHeight(vPos *math32.Vector3) (RtNavmeshCell, error) {
	var result RtNavmeshCellTri

	if !obj.Grid.Rectangle.ContainsVec3(vPos) {
		return result, errors.New("cell not in grid")
	}

	tile := obj.Grid.GetTileFromVec(vPos)
	minDeltaY := math32.Inf(1)
	y := math32.Inf(1)
	vTest := vPos.Clone()

	for _, cell := range tile.GetCells() {
		if cell.Triangle.FindHeight(vTest) {
			deltaY := math32.Abs(vTest.Y - vPos.Y)
			if deltaY < minDeltaY {
				minDeltaY = deltaY
				result = cell
				y = vTest.Y
			}
		}
	}

	vPos.Y = y
	return result, nil
}

func (obj RtNavmeshObj) ResolvePosition(pos RtNavmeshPosition) {
	// TODO
	panic("implement me")
}

func (obj RtNavmeshObj) TestOutlineIntersection(line LineSegment) bool {
	if obj.Grid.Rectangle.IntersectsLine(line) {
		return false
	}

	srcTileX := int((line.A.X - obj.Grid.X) / RtNavmeshObjGridTileWidth)
	srcTileZ := int((line.A.Z - obj.Grid.Y) / RtNavmeshObjGridTileHeight)

	dstTileX := int((line.B.X - obj.Grid.X) / RtNavmeshObjGridTileWidth)
	dstTileZ := int((line.B.Z - obj.Grid.Y) / RtNavmeshObjGridTileHeight)

	// TODO stay within grid
	//srcTileX = UnityEngine.Mathf.Clamp(srcTileX, 0, this.Grid.Width - 1);
	//srcTileY = UnityEngine.Mathf.Clamp(srcTileY, 0, this.Grid.Height - 1);
	//
	//dstTileX = UnityEngine.Mathf.Clamp(dstTileX, 0, this.Grid.Width - 1);
	//dstTileY = UnityEngine.Mathf.Clamp(dstTileY, 0, this.Grid.Height - 1);

	// swap if direction is negative
	tileX := srcTileX
	if srcTileX > dstTileX {
		srcTileX = dstTileX
		dstTileX = tileX
	}

	// swap if direction is negative
	tileZ := srcTileZ
	if srcTileZ > dstTileZ {
		srcTileZ = dstTileZ
		dstTileZ = tileZ
	}

	vSrc := math32.NewVector2(line.A.X, line.A.Z)

	var intersectionDistanceSquared float32 = 0.0
	var intersectionEdge RtNavmeshEdge = nil

	for tileZ = srcTileZ; tileZ <= dstTileZ; tileZ++ {
		for tileX = srcTileX; tileX <= dstTileX; tileX++ {
			tile := obj.Grid.GetTileFromXAndY(tileX, tileZ)

			for _, edge := range tile.GlobalEdges {
				if edge.GetFlag().IsBridge() {
					continue
				}

				if intersects, intersPos := edge.GetLine().Intersects(line); intersects {
					vDelta := vSrc.Sub(intersPos)
					if intersectionEdge == nil || intersectionDistanceSquared > vDelta.LengthSq() {
						intersectionDistanceSquared = vDelta.LengthSq()
						intersectionEdge = &edge
					}
				}
			}
		}
	}

	return intersectionEdge != nil
}

/*
GetNavmeshType() byte
	GetFilename() string
	GetCell(index int) RtNavmeshCell
	ResolveCellAndHeight(vPos *math32.Vector3)
	ResolvePosition(pos RtNavmeshPosition)
*/
