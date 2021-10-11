package navmeshv2

import (
	"github.com/g3n/engine/math32"
)

type RtNavmeshObj struct {
	RtNavmeshBase
	Cells         []RtNavmeshCellTri
	GlobalEdges   []RtNavmeshEdgeGlobal
	InternalEdges []RtNavmeshEdgeInternal
	Events        []string
	Grid          RtNavmeshObjGrid
}

func (obj RtNavmeshObj) IsPositionInObjectGrid(vPosLocal *math32.Vector3) bool {
	return obj.Grid.Rectangle.ContainsVec3(vPosLocal)
}

func (obj RtNavmeshObj) IsPositionInObjectCell(vPosLocal *math32.Vector3) bool {
	if !obj.IsPositionInObjectGrid(vPosLocal) {
		return false
	}

	for _, cell := range obj.Cells {
		if inTri, _ := cell.Triangle.FindHeight(vPosLocal); inTri {
			return true
		}
	}

	return false
}

func (obj RtNavmeshObj) FindHeight(vPosLocal *math32.Vector3) (bool, float32) {
	distance := math32.Inf(1)
	var nearestPos *math32.Vector3
	rayOrigin := vPosLocal.Clone()
	rayOrigin.Y += 1000
	for _, cell := range obj.Cells {
		ray := math32.NewRay(rayOrigin, math32.NewVector3(0, -1, 0))
		result := math32.NewVec3()
		intersects := ray.IntersectTriangle(cell.Triangle.A, cell.Triangle.B, cell.Triangle.C, true, result)
		if intersects && vPosLocal.DistanceToSquared(result) <= distance {
			distance = vPosLocal.DistanceToSquared(result)
			nearestPos = result
		}
	}

	if nearestPos == nil {
		return false, 0
	}

	return true, nearestPos.Y

}

func (obj RtNavmeshObj) ResolveCell(vPosLocal *math32.Vector3) (RtNavmeshCellTri, error) {
	if !obj.IsPositionInObjectGrid(vPosLocal) {
		return RtNavmeshCellTri{}, ErrCellNotInObject
	}
	distance := math32.Inf(1)
	var c *RtNavmeshCellTri
	for _, cell := range obj.Cells {
		if inTri, y := cell.Triangle.FindHeight2(vPosLocal); inTri && math32.Abs(vPosLocal.Y-y) < distance {
			distance = math32.Abs(vPosLocal.Y - y)
			c = &cell
		}
	}
	if c == nil {
		return RtNavmeshCellTri{}, ErrNoObjectCellForPosition
	} else {
		return *c, nil
	}
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
	return &obj.Cells[index]
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
