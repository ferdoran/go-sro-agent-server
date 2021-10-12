package navmeshv2

import (
	"github.com/g3n/engine/math32"
)

const (
	CollisionTerrain = iota
	CollisionObject
)

type CollisionFlag byte

func (f CollisionFlag) IsTerrain() bool {
	return f == CollisionTerrain
}

func (f CollisionFlag) IsObject() bool {
	return f == CollisionObject
}

type Collision struct {
	VectorLocal  *math32.Vector3
	VectorGlobal *math32.Vector3
	ContactEdge  RtNavmeshEdge
	Region
	Flag CollisionFlag
}

func FindTerrainCollisions(currentPosition, nextPosition RtNavmeshPosition,
	currentCell, nextCell RtNavmeshCellQuad) (bool, Collision) {
	collisions := make([]Collision, 0)
	line1 := LineSegment{
		A: currentPosition.Offset,
		B: nextPosition.Offset,
	}
	if currentPosition.Region.ID == nextPosition.Region.ID {

		if currentCell.Index == nextCell.Index {
			return false, Collision{}
		}

		for _, edge := range currentCell.Edges {
			if intersects, vIntersection := line1.Intersects(edge.GetLine()); intersects {
				if edge.GetFlag().IsBlocked() {
					collisions = append(collisions, Collision{
						VectorLocal:  math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						VectorGlobal: math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						ContactEdge:  edge,
						Region:       currentPosition.Region,
					})
				}
			}
		}

		for _, edge := range nextCell.Edges {
			if intersects, vIntersection := line1.Intersects(edge.GetLine()); intersects {
				if edge.GetFlag().IsBlocked() {
					collisions = append(collisions, Collision{
						VectorLocal:  math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						VectorGlobal: math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						ContactEdge:  edge,
						Region:       nextPosition.Region,
					})
				}
			}
		}

	} else {
		vCurGlobal := currentPosition.GetGlobalCoordinates()
		vNextGlobal := nextPosition.GetGlobalCoordinates()
		line1 = LineSegment{
			A: vCurGlobal,
			B: vNextGlobal,
		}

		for _, edge := range currentCell.Edges {
			a := RtNavmeshPosition{
				Cell:     &currentCell,
				Instance: nil,
				Region:   currentPosition.Region,
				Offset:   edge.GetLine().A,
			}
			b := RtNavmeshPosition{
				Cell:     &currentCell,
				Instance: nil,
				Region:   currentPosition.Region,
				Offset:   edge.GetLine().B,
			}

			line2 := LineSegment{
				A: a.GetGlobalCoordinates(),
				B: b.GetGlobalCoordinates(),
			}
			if intersects, vIntersection := line1.Intersects(line2); intersects {
				if edge.GetFlag().IsBlocked() {
					collisions = append(collisions, Collision{
						VectorLocal:  math32.NewVector3(float32(int(vIntersection.X)%int(RegionWidth)), 0, float32(int(vIntersection.Y)%int(RegionHeight))),
						VectorGlobal: math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						ContactEdge:  edge,
						Region:       currentPosition.Region,
					})
				}
			}
		}

		for _, edge := range nextCell.Edges {
			a := RtNavmeshPosition{
				Cell:     &currentCell,
				Instance: nil,
				Region:   currentPosition.Region,
				Offset:   edge.GetLine().A,
			}
			b := RtNavmeshPosition{
				Cell:     &currentCell,
				Instance: nil,
				Region:   currentPosition.Region,
				Offset:   edge.GetLine().B,
			}

			line2 := LineSegment{
				A: a.GetGlobalCoordinates(),
				B: b.GetGlobalCoordinates(),
			}
			if intersects, vIntersection := line1.Intersects(line2); intersects {
				if edge.GetFlag().IsBlocked() {
					collisions = append(collisions, Collision{
						VectorLocal:  math32.NewVector3(float32(int(vIntersection.X)%int(RegionWidth)), 0, float32(int(vIntersection.Y)%int(RegionHeight))),
						VectorGlobal: math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						ContactEdge:  edge,
						Region:       nextPosition.Region,
					})
				}
			}
		}
	}

	collision, err := FindNearestCollision(currentPosition.Offset, collisions)
	if err != nil {
		return false, Collision{}
	} else {
		return true, collision
	}
}

func FindObjectCollisions(currentPosition, nextPosition RtNavmeshPosition,
	currentObjects, nextObjects []RtNavmeshInstObj) (bool, Collision, bool, *math32.Vector3) {
	collides1, coll1, inObject1, objPos1 := FindCollisionsInObjects(currentPosition, nextPosition, currentObjects)
	collides2, coll2, inObject2, objPos2 := FindCollisionsInObjects(currentPosition, nextPosition, nextObjects)

	if collides1 && collides2 {
		distance1 := currentPosition.Offset.DistanceToSquared(coll1.VectorGlobal)
		distance2 := currentPosition.Offset.DistanceToSquared(coll2.VectorGlobal)
		if distance1 <= distance2 {
			return true, coll1, inObject1, coll1.VectorGlobal
		} else {
			return true, coll2, inObject2, coll2.VectorGlobal
		}
	} else if collides1 {
		return true, coll1, inObject1, coll1.VectorGlobal
	} else if collides2 {
		return true, coll2, inObject2, coll2.VectorGlobal
	} else if inObject1 {
		return false, Collision{}, inObject1, objPos1
	} else if inObject2 {
		return false, Collision{}, inObject2, objPos2
	} else {
		return false, Collision{}, false, nil
	}
}

func FindCollisionsInObjects(cur, next RtNavmeshPosition, objects []RtNavmeshInstObj) (bool, Collision, bool, *math32.Vector3) {
	inObj := false
	inObj2 := false
	var objPos *math32.Vector3
	for _, object := range objects {
		vCurLocal := cur.Offset.Clone().ApplyMatrix4(object.GetWorldToLocal())
		vNextLocal := next.Offset.Clone().ApplyMatrix4(object.GetWorldToLocal())

		curInObject := object.Object.IsPositionInObjectCell(vCurLocal)
		nextInObject := object.Object.IsPositionInObjectCell(vNextLocal)
		if !curInObject && !nextInObject {
			// skip
			continue
		}

		line1 := LineSegment{
			A: vCurLocal.Clone(),
			B: vNextLocal.Clone(),
		}

		gCollisions := make([]Collision, 0)
		iCollisions := make([]Collision, 0)

		for _, gEdge := range object.Object.GlobalEdges {
			line2 := gEdge.GetLine()
			if intersects, vIntersection := line1.Intersects(line2); intersects {
				if gEdge.GetFlag().IsBlocked() || (curInObject && gEdge.GetFlag().IsBridge()) {
					gCollisions = append(gCollisions, Collision{
						VectorLocal:  math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						VectorGlobal: math32.NewVector3(vIntersection.X, 0, vIntersection.Y).ApplyMatrix4(object.GetLocalToWorld()),
						ContactEdge:  &gEdge,
						Region:       object.Region,
					})
				} else {
					inObj = true
				}
			}
		}

		for _, iEdge := range object.Object.InternalEdges {
			line2 := iEdge.GetLine()
			if intersects, vIntersection := line1.Intersects(line2); intersects {
				if iEdge.GetFlag().IsBlocked() {
					iCollisions = append(iCollisions, Collision{
						VectorLocal:  math32.NewVector3(vIntersection.X, 0, vIntersection.Y),
						VectorGlobal: math32.NewVector3(vIntersection.X, 0, vIntersection.Y).ApplyMatrix4(object.GetLocalToWorld()),
						ContactEdge:  &iEdge,
						Region:       object.Region,
					})
				} else {
					inObj2 = true
				}
			}
		}

		filteredGlobalCollisions := make([]Collision, 0)
		for _, collision := range gCollisions {
			flag := collision.ContactEdge.GetFlag()
			if curInObject && flag.IsBridge() {
				filteredGlobalCollisions = append(filteredGlobalCollisions, collision)
			}
		}

		finalCollisions := make([]Collision, 0)
		nearestGlobalCollision, err1 := FindNearestCollision(vCurLocal, filteredGlobalCollisions)
		nearestLocalCollision, err2 := FindNearestCollision(vCurLocal, iCollisions)

		if err1 == nil {
			finalCollisions = append(finalCollisions, nearestGlobalCollision)
		}
		if err2 == nil {
			finalCollisions = append(finalCollisions, nearestLocalCollision)
		}

		finalCollision, err3 := FindNearestCollision(vCurLocal, finalCollisions)

		if inObj || inObj2 {
			_, height := object.Object.FindHeight(vNextLocal)
			objPos = vNextLocal
			objPos.Y = height
			objPos.ApplyMatrix4(object.LocalToWorld)
		}
		if err3 == ErrEmptyCollisionList {
			// no collision!
			continue
		}

		return true, finalCollision, inObj || inObj2, objPos
	}

	return false, Collision{}, inObj || inObj2, objPos
}

func FindNearestCollision(localPosition *math32.Vector3, collisions []Collision) (Collision, error) {
	delta := math32.Inf(1)
	var nearestCollision Collision

	if len(collisions) == 0 {
		return nearestCollision, ErrEmptyCollisionList
	}

	for _, collision := range collisions {
		if distance := localPosition.DistanceToSquared(collision.VectorLocal); distance <= delta {
			delta = distance
			nearestCollision = collision
		}
	}

	return nearestCollision, nil
}
