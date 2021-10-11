package viewer

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/g3n/engine/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sirupsen/logrus"
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

	rl.DrawCircle3D(rl.NewVector3(float32(collisionPoint1X), float32(collisionPoint1Y), float32(collisionPoint1Z)), 5, rl.NewVector3(1, 0, 0), 90.0, rl.White)
	rl.DrawCircle3D(rl.NewVector3(float32(collisionPoint2X), float32(collisionPoint2Y), float32(collisionPoint2Z)), 5, rl.NewVector3(1, 0, 0), 90.0, rl.White)

	p1 := math32.NewVector3(float32(collisionPoint1X), float32(collisionPoint1Y), float32(collisionPoint1Z))
	p2 := math32.NewVector3(float32(collisionPoint2X), float32(collisionPoint2Y), float32(collisionPoint2Z))

	//rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), rl.NewVector3(p2.X, p2.Y, p2.Z), rl.Green)
	region := loader.RegionData[int16(selectedRegionId)]
	cellsAlongPath := make(map[int]navmeshv2.RtNavmeshCellQuad)

	direction := p2.Clone().Sub(p1.Clone())
	distance := direction.Length()
	steps := distance / 20
	ray := math32.NewRay(p1, direction.Normalize())
	for i := 0; i < int(steps); i++ {
		pos := ray.At(float32(i)*20, nil)
		cell, err := region.ResolveCell(pos)
		if err != nil {
			logrus.Error("failed to resolve cell", err)
			continue
		}

		cellsAlongPath[cell.Index] = cell
	}
	cellsOnPath = len(cellsAlongPath)

	//drawTerrainCollisions(p1, p2, region)

	drawObjectCollisions(cellsAlongPath, p1, p2)
}

func drawObjectCollisions(cells map[int]navmeshv2.RtNavmeshCellQuad, p1 *math32.Vector3, p2 *math32.Vector3) {
	objectCollisions := make([]navmeshv2.Collision, 0)
	numObjects := 0
	prevObjectId := int16(0)
	for _, cell := range cells {
		for _, object := range cell.Objects {
			if prevObjectId != object.ID {
				numObjects++
			}
			findObjectCollisions(p1, object, p2, &objectCollisions)
		}
	}
	objectsOnPath = numObjects

	if len(objectCollisions) > 0 {
		objectCollisionsOnPath = len(objectCollisions)
		collisionPoint := getClosestObjectCollision(p1, objectCollisions)
		rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), collisionPoint, rl.Red)
		rl.DrawLine3D(collisionPoint, rl.NewVector3(p2.X, p2.Y, p2.Z), rl.Green)
	} else {
		rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), rl.NewVector3(p2.X, p2.Y, p2.Z), rl.Green)
		objectCollisionsOnPath = 0
	}
}

func getClosestObjectCollision(from *math32.Vector3, collisions []navmeshv2.Collision) rl.Vector3 {
	distance := math32.Inf(1)
	var closestCollision *math32.Vector3

	for _, collision := range collisions {
		if d := from.DistanceToSquared(collision.VectorGlobal); d <= distance {
			distance = d
			closestCollision = collision.VectorGlobal
		}
	}

	return rl.Vector3{
		X: closestCollision.X,
		Y: closestCollision.Y,
		Z: closestCollision.Z,
	}
}

func findObjectCollisions(p1 *math32.Vector3, object navmeshv2.RtNavmeshInstObj, p2 *math32.Vector3, objectCollisions *[]navmeshv2.Collision) {
	p1Local := p1.Clone().ApplyMatrix4(object.GetWorldToLocal())
	p2Local := p2.Clone().ApplyMatrix4(object.GetWorldToLocal())
	line := navmeshv2.LineSegment{
		A: p1Local,
		B: p2Local,
	}

	p1InObj := false
	_, err := object.Object.ResolveCell(p1Local)
	if err == nil {

		inCell, y := object.Object.FindHeight(p1Local)
		if inCell && math32.Abs(p1Local.Y-y) < 50 {
			// p1 is in object
			p1InObj = true
		} else {
			p := math32.NewVector3(p1Local.X, y, p1Local.Z)
			p.ApplyMatrix4(object.GetLocalToWorld())
			logrus.Infof("p1 not in cell: %t, y1: %f, y2: %f", inCell, p1.Y, p.Y)
		}
	}

	for _, edge := range object.Object.GlobalEdges {
		edgeLine := edge.GetLine()
		if intersects, collision := edgeLine.Intersects3D(line); intersects {
			if edge.GetFlag().IsBlocked() || (p1InObj && edge.GetFlag().IsBridge()) {
				if edge.GetFlag().IsBridge() {
					//logrus.Infof("y1: %f, y2: %f", p1.Y, collision.Clone().ApplyMatrix4(object.GetLocalToWorld()).Y)
					logrus.Infof("y1: %f, y2: %f", p1Local.Y, collision.Y)
				}
				if math32.Abs(p1Local.Y-collision.Y) > 50 {
					logrus.Infof("no collision: y1: %f, y2: %f", p1Local.Y, collision.Y)
					continue
				}
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

	for _, edge := range object.Object.InternalEdges {
		edgeLine := edge.GetLine()

		if intersects, collision := edgeLine.Intersects3D(line); intersects {
			if edge.GetFlag().IsBlocked() {
				if collision.Y < p1Local.Y && p1Local.Y-collision.Y > 50 {
					continue
				}
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
}

func drawTerrainCollisions(p1, p2 *math32.Vector3, region navmeshv2.RtNavmeshTerrain) {
	line := navmeshv2.LineSegment{
		A: p1,
		B: p2,
	}
	terrainCollisions := make([]*math32.Vector2, 0)
	for _, edge := range region.InternalEdges {
		if intersects, collisionPoint := edge.GetLine().Intersects(line); intersects {
			if edge.GetFlag().IsBlocked() {
				terrainCollisions = append(terrainCollisions, collisionPoint)
				rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), rl.NewVector3(collisionPoint.X, 0, collisionPoint.Y), rl.Red)
			}
		}
	}
	if len(terrainCollisions) > 0 {
		collision := getClosestTerrainCollision(p1, terrainCollisions)
		rl.DrawLine3D(rl.NewVector3(p1.X, p1.Y, p1.Z), collision, rl.Red)
	}
}

func getClosestTerrainCollision(from *math32.Vector3, collisions []*math32.Vector2) rl.Vector3 {
	fromV2 := math32.NewVector2(from.X, from.Z)
	distance := math32.Inf(1)
	var closestCollision *math32.Vector2

	for _, collision := range collisions {
		if d := fromV2.DistanceToSquared(collision); d <= distance {
			distance = d
			closestCollision = collision
		}
	}

	return rl.Vector3{
		X: closestCollision.X,
		Y: 0,
		Z: closestCollision.Y,
	}
}
