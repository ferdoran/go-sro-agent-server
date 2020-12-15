package geo

import (
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/fileutils/navmesh"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/math"
)

func getObjectTileFromLocalPoint(object *navmesh.Object, point *math32.Vector3) *navmesh.ObjectTile {
	x := int(point.X - object.Grid.Origin.X/100)
	z := int(point.Z - object.Grid.Origin.Z/100)

	tileIdx := x + object.Grid.Height*z

	if tileIdx >= len(object.Grid.GridTiles) {
		return nil
	}

	return object.Grid.GridTiles[tileIdx]
}

func FindTriangleHeight(t *math.Triangle, vPos *math32.Vector3) bool {
	//From http://totologic.blogspot.se/2014/01/accurate-point-in-triangle-test.html

	//Based on Barycentric coordinates
	var denominator float32 = ((t.B.Z - t.C.Z) * (t.A.X - t.C.X)) + ((t.C.X - t.B.X) * (t.A.Z - t.C.Z))
	var a float32 = (((t.B.Z - t.C.Z) * (vPos.X - t.C.X)) + ((t.C.X - t.B.X) * (vPos.Z - t.C.Z))) / denominator
	var b float32 = (((t.C.Z - t.A.Z) * (vPos.X - t.C.X)) + ((t.A.X - t.C.X) * (vPos.Z - t.C.Z))) / denominator
	var c float32 = 1 - a - b
	vPos.Y = ((a * t.A.Y) + (b * t.B.Y) + (c * t.C.Y)) / (a + b + c)
	//return a > 0f && a < 1f && b > 0f && b < 1f && c > 0f && c < 1f; // point is within the triangle
	return a >= 0 && a <= 1 && b >= 0 && b <= 1 && c >= 0 && c <= 1 // point can be on border
}

func HasOuterCollision2(a, b *math32.Vector3, object *navmesh.Object) (bool, []Collision, bool) {
	v0 := math32.NewVector2(a.X, a.Z)
	v1 := math32.NewVector2(b.X, b.Z)
	line := math.NewLine2(v0, v1)
	colls := make([]Collision, 0)
	hasEnteredObject := false
	for _, edge := range object.GlobalEdges {
		if intersects, collision := edge.ToLine2().Intersects(line); intersects {
			diff := math32.NewVec3().SubVectors(edge.B, edge.A)
			y := edge.B.Y - diff.Y*((edge.B.X-collision.X)/diff.X)

			if math32.Abs(y-a.Y) > 50 {
				continue
			}
			if edge.Flag&(1|2|16) != 0 { // 1 = BlockDst2Src | 2 = BlockSrc2Dst | 16 = Bridge
				// Blocked

				v := math32.NewVector3(collision.X, a.Y, collision.Y)
				logrus.Tracef("collision: (%f|%f|%f)", v.X, v.Y, v.Z)
				c := FindHeightInObject(v, object)
				if c == v {
					// collision point not in object
					continue
				}
				v = c

				logrus.Tracef("vY: %f\n", v.Y)
				logrus.Tracef("edgeY: %f\n", y)
				collision := Collision{
					EdgeFlag:    edge.Flag,
					VectorLocal: v.Clone(),
				}
				collision.VectorGlobal = v.ApplyMatrix4(object.LocalToWorld)
				colls = append(colls, collision)
				logrus.Tracef("edge flag %v\n", edge.Flag)
			} else {
				hasEnteredObject = true
			}
		}
	}

	return len(colls) > 0, colls, hasEnteredObject
}

func FindCollisions(a, b *math32.Vector3, aRegionID, bRegionID int16, aObjects, bObjects []*navmesh.Object) (bool, Collision, bool, *math32.Vector3) {
	aVec := math32.NewVector3(a.X, a.Y, a.Z)
	bVec := math32.NewVector3(b.X, b.Y, b.Z)
	aObjs := filterObjectsContainingPositionInGrid(aObjects, aVec)
	bObjs := filterObjectsContainingPositionInGrid(bObjects, bVec)

	if aRegionID == bRegionID {
		logrus.Tracef("checking collision with %d objects", len(aObjs))
		for _, obj := range aObjs {
			aVecLocal := aVec.Clone().ApplyMatrix4(obj.WorldToLocal)
			bVecLocal := bVec.Clone().ApplyMatrix4(obj.WorldToLocal)

			collides, collision, inObject, objPosition := CollidesWithObject(aVecLocal, bVecLocal, obj)

			if collides || inObject {
				return collides, collision, inObject, objPosition
			}
		}
	} else {
		for _, o1 := range aObjs {
			for _, o2 := range bObjs {
				if o1.ID == o2.ID {
					aVecLocal := aVec.Clone().ApplyMatrix4(o1.WorldToLocal)
					bVecLocal := bVec.Clone().ApplyMatrix4(o2.WorldToLocal)
					// TODO: Are there any conflicts with this approach?
					collides, collision, inObject, objPosition := CollidesWithObject(aVecLocal, bVecLocal, o1)

					if collides || inObject {
						return collides, collision, inObject, objPosition
					}
				}
			}
		}

	}

	return false, Collision{}, false, nil
}

func filterObjectsContainingPositionInGrid(objects []*navmesh.Object, pos *math32.Vector3) []*navmesh.Object {
	filteredObjects := make([]*navmesh.Object, 0)

	for _, o := range objects {
		if o.Grid.ContainsPoint(pos.Clone().ApplyMatrix4(o.WorldToLocal)) {
			filteredObjects = append(filteredObjects, o)
		}
	}

	return filteredObjects
}

func IsNextPositionTooHigh(curPos, newPos *math32.Vector3) bool {
	//x, _, z := curPos.ToWorldCoordinates()
	//x1, _, z1 := newPos.ToWorldCoordinates()
	v0 := math32.NewVector3(curPos.X, 0, curPos.Z)
	v1 := math32.NewVector3(newPos.X, 0, newPos.Z)

	adjacent := v0.DistanceTo(v1)
	opposite := math32.Abs(curPos.Y - newPos.Y)
	angle := math32.RadToDeg(math32.Atan(opposite / adjacent))
	logrus.Tracef("next pos angle is %f°\n", angle)
	return math32.Abs(newPos.Y-curPos.Y) > 100
}

func CollidesWithObject(a, b *math32.Vector3, object *navmesh.Object) (bool, Collision, bool, *math32.Vector3) {
	//srcInObject, destInObject := false, false

	srcInObject := IsPointInObject(a, object)
	b2 := b.Clone()
	b2.Y = a.Y
	destInObject := IsPointInObject(b2, object)
	//for _, cell := range object.Cells {
	//	if cell.ToTriangle2().PointInTriangle(a) {
	//
	//		srcInObject = true
	//	}
	//
	//	if cell.ToTriangle2().PointInTriangle(b) {
	//		destInObject = true
	//	}
	//}

	if srcInObject && !destInObject {
		// Leaving object
		// Check inner & outer collisions
		hasInnerColl, innerColls, inObject := HasInnerCollision2(a, b, object)
		logrus.Tracef("leaving object %d", object.ID)

		if hasInnerColl {
			collisionPoint := GetClosestCollisionPoint2(a, innerColls)
			logrus.Debugf("inner collision at (%f|%f|%f) with object %d", collisionPoint.VectorGlobal.X, collisionPoint.VectorGlobal.Y, collisionPoint.VectorGlobal.Z, object.ID)
			return true, collisionPoint, true, nil
		}

		hasOuterColl, outerColls, inObject2 := HasOuterCollision2(a, b, object)
		if hasOuterColl {
			collisionPoint := GetClosestCollisionPoint2(a, outerColls)
			logrus.Debugf("outer collision at (%f|%f|%f) with object %d", collisionPoint.VectorGlobal.X, collisionPoint.VectorGlobal.Y, collisionPoint.VectorGlobal.Z, object.ID)
			return true, collisionPoint, true, nil
		}

		bObjPos := b.Clone()
		bObjPos.Y = a.Y
		bObjPos = FindHeightInObject(bObjPos, object)
		bObjPos = bObjPos.ApplyMatrix4(object.LocalToWorld)
		return false, Collision{}, inObject || inObject2, bObjPos

	} else if !srcInObject && destInObject {
		// Entering object
		// check outer collisions
		logrus.Tracef("entering object %d", object.ID)
		hasOuterColl, outerColls, inObject := HasOuterCollision2(a, b, object)
		if hasOuterColl {
			collisionPoint := GetClosestCollisionPoint2(a, outerColls)
			logrus.Debugf("outer collision at (%f|%f|%f) with object %d", collisionPoint.VectorGlobal.X, collisionPoint.VectorGlobal.Y, collisionPoint.VectorGlobal.Z, object.ID)
			return true, collisionPoint, inObject, nil
		}
		hasInnerColl, innerColls, inObject2 := HasInnerCollision2(a, b, object)

		if hasInnerColl {
			collisionPoint := GetClosestCollisionPoint2(a, innerColls)
			logrus.Debugf("inner collision at (%f|%f|%f) with object %d", collisionPoint.VectorGlobal.X, collisionPoint.VectorGlobal.Y, collisionPoint.VectorGlobal.Z, object.ID)

			return true, collisionPoint, true, nil
		}
		bObjPos := b.Clone()
		bObjPos.Y = a.Y
		bObjPos = FindHeightInObject(bObjPos, object)
		bObjPos = bObjPos.ApplyMatrix4(object.LocalToWorld)
		return false, Collision{}, inObject || inObject2, bObjPos
	} else if srcInObject && destInObject {
		// Staying in object
		// check inner collisions
		logrus.Tracef("staying in object %d", object.ID)
		hasInnerColl, innerColls, _ := HasInnerCollision2(a, b, object)

		if hasInnerColl {
			collisionPoint := GetClosestCollisionPoint2(a, innerColls)
			logrus.Tracef("inner collision at (%f|%f|%f) with object %d", collisionPoint.VectorGlobal.X, collisionPoint.VectorGlobal.Y, collisionPoint.VectorGlobal.Z, object.ID)
			return true, collisionPoint, true, nil
		}

		hasOuterColl, outerColls, _ := HasOuterCollision2(a, b, object)
		if hasOuterColl {
			collisionPoint := GetClosestCollisionPoint2(a, outerColls)
			logrus.Tracef("outer collision at (%f|%f|%f) with object %d", collisionPoint.VectorGlobal.X, collisionPoint.VectorGlobal.Y, collisionPoint.VectorGlobal.Z, object.ID)
			return true, collisionPoint, true, nil
		}

		bObjPos := b.Clone()
		bObjPos.Y = a.Y
		bObjPos = FindHeightInObject(bObjPos, object)
		bObjPos = bObjPos.ApplyMatrix4(object.LocalToWorld)
		return false, Collision{}, true, bObjPos
	}

	logrus.Tracef("no collision with object %d", object.ID)

	return false, Collision{}, false, nil
}

func HasInnerCollision2(a, b *math32.Vector3, object *navmesh.Object) (bool, []Collision, bool) {
	colls := make([]Collision, 0)
	inObjSpace := false
	v0 := math32.NewVector2(a.X, a.Z)
	v1 := math32.NewVector2(b.X, b.Z)
	line := math.NewLine2(v0, v1)
	for _, edge := range object.InternalEdges {
		if intersects, collision := edge.ToLine2().Intersects(line); intersects {
			diff := math32.NewVec3().SubVectors(edge.B, edge.A)
			y := edge.B.Y - diff.Y*((edge.B.X-collision.X)/diff.X)

			if math32.Abs(y-a.Y) > 50 {
				logrus.Tracef("edge intersection diff too big. newY %f, oldY %f, flag %d", y, a.Y, edge.Flag)
				continue
			}
			if edge.Flag == 7 { // 1 = BlockedDst2Src, 2 = BlockedSrc2Dst, 7 == Blocked, Internal
				// Blocked
				v := math32.NewVector3(collision.X, a.Y, collision.Y)
				logrus.Tracef("collision: (%f|%f|%f)", v.X, v.Y, v.Z)
				c := FindHeightInObject(v, object)
				if c == v {
					// collision point not in object
					continue
				}
				v = c

				logrus.Tracef("vY: %f\n", v.Y)
				logrus.Tracef("edgeY: %f\n", y)
				collision := Collision{
					EdgeFlag:    edge.Flag,
					VectorLocal: v.Clone(),
				}
				collision.VectorGlobal = v.ApplyMatrix4(object.LocalToWorld)
				colls = append(colls, collision)
				logrus.Tracef("edge flag %v\n", edge.Flag)
			} else {
				inObjSpace = true
			}
		}
	}

	return len(colls) > 0, colls, inObjSpace
}

func GetClosestCollisionPoint2(pos *math32.Vector3, collisions []Collision) Collision {
	closestCollision := Collision{}
	closestDistance := math32.Infinity

	for _, c := range collisions {
		line := math32.NewLine3(pos, c.VectorGlobal)
		if line.Distance() <= closestDistance {
			closestDistance = line.Distance()
			closestCollision = c
		}
	}

	return closestCollision
}

func FindHeightInObject(pos *math32.Vector3, obj *navmesh.Object) *math32.Vector3 {
	colls := make([]*math32.Vector3, 0)
	for _, cell := range obj.Cells {
		if cell.ToTriangle2().PointInTriangle(pos) {
			p := pos.Clone()
			abovePos := math32.NewVector3(pos.X, pos.Y+1000, pos.Z)
			direction := math32.NewVector3(0, -1, 0)
			ray := math32.NewRay(abovePos, direction)
			if ray.IntersectTriangle(cell.A, cell.B, cell.C, false, p) {
				colls = append(colls, p.ApplyMatrix4(obj.LocalToWorld))
			}
		}
	}
	p := pos.Clone().ApplyMatrix4(obj.LocalToWorld)
	closestCollision := p
	closestYDistance := math32.Inf(1)
	if size := len(colls); size > 1 {
		logrus.Tracef("having multiple vertical positions in object %d: %d", obj.ID, size)
	}
	// FIXME somehow the calculated distance is lower for a cell that actually is further away
	for _, coll := range colls {
		if distance := math32.Abs(coll.Y - p.Y); distance <= closestYDistance {
			logrus.Tracef("Distance: %f, collisionY = %f, positionY: %f", distance, coll.Y, p.Y)
			closestCollision = coll
			closestYDistance = distance
		}
	}

	return closestCollision.ApplyMatrix4(obj.WorldToLocal)
}

func IsPointInObject(point *math32.Vector3, obj *navmesh.Object) bool {
	p := point.Clone()
	for _, cell := range obj.Cells {
		if FindTriangleHeight(cell.Triangle, p) && math32.Abs(p.Y-point.Y) <= 100 {
			return true
		}
	}

	return false
}
