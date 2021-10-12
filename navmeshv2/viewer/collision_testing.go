package viewer

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/g3n/engine/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var collisionPoint1X = 950
var collisionPoint1Y = 0
var collisionPoint1Z = 1300
var collisionPoint2X = 950
var collisionPoint2Y = 0
var collisionPoint2Z = 1500

var collisionPoint1XText = "950"
var collisionPoint1YText = ""
var collisionPoint1ZText = "1300"
var collisionPoint2XText = "950"
var collisionPoint2YText = ""
var collisionPoint2ZText = "1500"
var cellsOnPath = 0
var objectsOnPath = 0
var objectCollisionsOnPath = 0

func drawCollisionPoints() {

	region := loader.RegionData[int16(selectedRegionId)]

	p1 := math32.NewVector3(float32(collisionPoint1X), float32(collisionPoint1Y), float32(collisionPoint1Z))
	p2 := math32.NewVector3(float32(collisionPoint2X), float32(collisionPoint2Y), float32(collisionPoint2Z))

	//p1.Y = region.ResolveHeight(p1)
	//p2.Y = region.ResolveHeight(p2)

	//collisionPoint1Y = int(p1.Y)
	//collisionPoint2Y = int(p2.Y)
	//collisionPoint1YText = fmt.Sprintf("%d", collisionPoint1Y)
	//collisionPoint2YText = fmt.Sprintf("%d", collisionPoint2Y)

	rl.DrawCircle3D(rl.NewVector3(p1.X, p1.Y, p1.Z), 10, rl.NewVector3(1, 0, 0), 90.0, rl.White)
	rl.DrawCircle3D(rl.NewVector3(p2.X, p2.Y, p2.Z), 10, rl.NewVector3(1, 0, 0), 90.0, rl.White)

	//rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), rl.NewVector3(p2.X, p2.Y, p2.Z), rl.Green)
	cellsAlongPath := make(map[int]navmeshv2.RtNavmeshCellQuad)

	direction := p2.Clone().Sub(p1.Clone())
	distance := direction.Length()
	steps := distance / 20
	ray := math32.NewRay(p1, direction.Normalize())
	lastPos := p1.Clone()
	collisions := make([]navmeshv2.Collision, 0)
	onObject := false
	for i := 1; i < int(steps); i++ {
		pos := ray.At(float32(i)*20, nil)
		cell, err := region.ResolveCell(pos)
		if err != nil {
			continue
		}
		cellsAlongPath[cell.Index] = cell
		// 1. Find terrain collisions
		terrainCollisions := make([]navmeshv2.Collision, 0)
		drawTerrainCollisions(lastPos, pos, region, &terrainCollisions)

		// 2. Find object collisions
		allObjectCollisions := make([]navmeshv2.Collision, 0)
		for _, object := range cell.Objects {
			objectCollisions := make([]navmeshv2.Collision, 0)
			enter, _, leave, _ := findObjectCollisions(lastPos, pos, object, &objectCollisions)
			// 3a. If object is being entered, check if terrain collision is before
			if !onObject && enter {
				onObject = true
			} else if onObject && leave {
				onObject = false
			}

			allObjectCollisions = append(allObjectCollisions, objectCollisions...)
		}

		if !onObject && (len(terrainCollisions) > 0 || len(allObjectCollisions) > 0) {
			collisions = append(collisions, terrainCollisions...)
			collisions = append(collisions, allObjectCollisions...)
			break
		} else if onObject && len(allObjectCollisions) > 0 {
			collisions = append(collisions, allObjectCollisions...)
			break
		}

		lastPos = pos

	}

	cellsOnPath = len(cellsAlongPath)
	//terrainCollisions := make([]navmeshv2.Collision, 0)
	//objectCollisions := make([]navmeshv2.Collision, 0)
	//drawTerrainCollisions(p1, p2, region, &terrainCollisions)
	//drawObjectCollisions(cellsAlongPath, p1, p2, &objectCollisions)

	if len(collisions) > 0 {
		objectCollisionsOnPath = len(collisions)
		closestCollision := getClosestCollision(p1, collisions)
		collisionVec := rl.NewVector3(closestCollision.VectorGlobal.X, closestCollision.VectorGlobal.Y, closestCollision.VectorGlobal.Z)
		rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), collisionVec, rl.Red)
		rl.DrawLine3D(collisionVec, rl.NewVector3(p2.X, p2.Y, p2.Z), rl.Green)
	} else {
		rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), rl.NewVector3(p2.X, p2.Y, p2.Z), rl.Green)
		objectCollisionsOnPath = 0
	}
}

func drawObjectCollisions(cells map[int]navmeshv2.RtNavmeshCellQuad, p1 *math32.Vector3, p2 *math32.Vector3, collisions *[]navmeshv2.Collision) {
	numObjects := 0
	prevObjectId := int16(0)
	for _, cell := range cells {
		for _, object := range cell.Objects {
			if prevObjectId != object.ID {
				numObjects++
			}
			findObjectCollisions(p1, p2, object, collisions)
		}
	}
	objectsOnPath = numObjects
}

func getClosestCollision(from *math32.Vector3, collisions []navmeshv2.Collision) navmeshv2.Collision {
	distance := math32.Inf(1)
	var closestCollision navmeshv2.Collision

	for _, collision := range collisions {
		if d := from.DistanceToSquared(collision.VectorGlobal); d <= distance {
			distance = d
			closestCollision = collision
		}
	}

	return closestCollision
}

func findObjectCollisions(p1 *math32.Vector3, p2 *math32.Vector3, object navmeshv2.RtNavmeshInstObj, objectCollisions *[]navmeshv2.Collision) (bool, *math32.Vector3, bool, *math32.Vector3) {
	p1Local := p1.Clone().ApplyMatrix4(object.GetWorldToLocal())
	p2Local := p2.Clone().ApplyMatrix4(object.GetWorldToLocal())
	line := navmeshv2.LineSegment{
		A: p1Local,
		B: p2Local,
	}

	p1InObj := isPointInObject(p1Local, object)
	p2InObj := isPointInObject(p2Local, object)

	enteringObject := false
	var enteringObjectPosition *math32.Vector3
	leavingObject := false
	var leavingObjectPosition *math32.Vector3

	for _, edge := range object.Object.GlobalEdges {
		edgeLine := edge.GetLine()
		if intersects, collision := edgeLine.Intersects3D(line); intersects {
			if math32.Abs(p1Local.Y-collision.Y) > 50 {
				continue
			}
			if edge.GetFlag().IsBlocked() || (p1InObj && edge.GetFlag().IsBridge()) {
				coll := navmeshv2.Collision{
					VectorLocal:  collision,
					VectorGlobal: collision.Clone().ApplyMatrix4(object.GetLocalToWorld()),
					ContactEdge:  &edge,
					Region:       navmeshv2.Region{},
				}
				*objectCollisions = append(*objectCollisions, coll)
			} else if !p1InObj && p2InObj {
				// entering object
				enteringObject = true
				enteringObjectPosition = collision
			} else if p1InObj && !p2InObj {
				// leaving object
				leavingObject = true
				leavingObjectPosition = collision
			}
		}
	}

	for _, edge := range object.Object.InternalEdges {
		edgeLine := edge.GetLine()

		if intersects, collision := edgeLine.Intersects3D(line); intersects {
			if collision.Y < p1Local.Y && p1Local.Y-collision.Y > 50 {
				continue
			}
			if edge.GetFlag().IsBlocked() {
				coll := navmeshv2.Collision{
					VectorLocal:  collision,
					VectorGlobal: collision.Clone().ApplyMatrix4(object.GetLocalToWorld()),
					ContactEdge:  &edge,
					Region:       navmeshv2.Region{},
				}
				*objectCollisions = append(*objectCollisions, coll)
			}
		}
	}

	return enteringObject, enteringObjectPosition, leavingObject, leavingObjectPosition
}

func isPointInObject(pLocal *math32.Vector3, object navmeshv2.RtNavmeshInstObj) bool {
	inObject := false
	_, err := object.Object.ResolveCell(pLocal)
	if err == nil {

		inCell, y := object.Object.FindHeight(pLocal)
		if inCell && math32.Abs(pLocal.Y-y) < 50 {
			// p1 is in object
			inObject = true
		}
	}
	return inObject
}

func drawTerrainCollisions(p1, p2 *math32.Vector3, region navmeshv2.RtNavmeshTerrain, collisions *[]navmeshv2.Collision) {
	line := navmeshv2.LineSegment{
		A: p1,
		B: p2,
	}
	for _, edge := range region.InternalEdges {
		if intersects, collisionPoint := edge.GetLine().Intersects(line); intersects {
			if edge.GetFlag().IsBlocked() {
				y := region.ResolveHeight(math32.NewVector3(collisionPoint.X, 0, collisionPoint.Y))
				coll := navmeshv2.Collision{
					VectorLocal:  math32.NewVector3(collisionPoint.X, y, collisionPoint.Y),
					VectorGlobal: math32.NewVector3(collisionPoint.X, y, collisionPoint.Y),
					ContactEdge:  &edge,
					Region:       navmeshv2.Region{},
				}
				*collisions = append(*collisions, coll)
			}
		}
	}
}
