package viewer

import rl "github.com/gen2brain/raylib-go/raylib"

var cam = defaultCam()
var drawObjectsWithZeroHeight = true

func drawTerrain() {
	if selectedRegionId != 0 {
		rl.BeginScissorMode(int32(screenWidth*TerrainListPaneWidthFactor), ToolbarHeight, int32(screenWidth*(1-TerrainListPaneWidthFactor)), int32(screenHeight-ToolbarHeight))
		defer rl.EndScissorMode()
		rl.BeginMode3D(cam)
		defer rl.EndMode3D()
		terrain := terrainList[int16(selectedRegionId)]

		rl.DrawCubeWires(rl.NewVector3(1920/2, 0, 1920/2), 1920, 0, 1920, rl.Green)

		if showTerrainCells {
			if !showOnlyBlockedRegionCells {
				for _, cell := range terrain.Cells {
					p1 := rl.NewVector3(cell.Rect.Min.X, 0, cell.Rect.Min.Y)
					p2 := rl.NewVector3(cell.Rect.Max.X, 0, cell.Rect.Min.Y)
					p3 := rl.NewVector3(cell.Rect.Min.X, 0, cell.Rect.Max.Y)
					p4 := rl.NewVector3(cell.Rect.Max.X, 0, cell.Rect.Max.Y)

					rl.DrawLine3D(p1, p2, rl.Blue)
					rl.DrawLine3D(p1, p3, rl.Blue)
					rl.DrawLine3D(p3, p4, rl.Blue)
					rl.DrawLine3D(p2, p4, rl.Blue)
				}
			}

			for _, edge := range terrain.GlobalEdges {
				startPos := rl.NewVector3(edge.Line.A.X, edge.Line.A.Y, edge.Line.A.Z)
				endPos := rl.NewVector3(edge.Line.B.X, edge.Line.B.Y, edge.Line.B.Z)

				if edge.GetFlag().IsBlocked() {
					rl.DrawLine3D(startPos, endPos, rl.Maroon)
				}
			}

			for _, edge := range terrain.InternalEdges {
				startPos := rl.NewVector3(edge.Line.A.X, edge.Line.A.Y, edge.Line.A.Z)
				endPos := rl.NewVector3(edge.Line.B.X, edge.Line.B.Y, edge.Line.B.Z)

				if edge.GetFlag().IsBlocked() {
					rl.DrawLine3D(startPos, endPos, rl.Maroon)
				}
			}
		}

		for _, obj := range terrain.Objects {
			if showObjectsGlobalEdges {
				for _, edge := range obj.Object.GlobalEdges {
					var color = rl.Yellow

					if edge.GetFlag().IsBlocked() {
						color = rl.Red
					} else if edge.GetFlag().IsBridge() {
						color = rl.DarkGreen
					}
					startPosVec := edge.Line.A.Clone().ApplyMatrix4(obj.GetLocalToWorld())
					endPosVec := edge.Line.B.Clone().ApplyMatrix4(obj.GetLocalToWorld())

					if drawObjectsWithZeroHeight {
						startPosVec.Y = 0
						endPosVec.Y = 0
					}

					startPos := rl.NewVector3(startPosVec.X, startPosVec.Y, startPosVec.Z)
					endPos := rl.NewVector3(endPosVec.X, endPosVec.Y, endPosVec.Z)
					rl.DrawLine3D(startPos, endPos, color)
				}
			}

			if showObjectsInternalEdges {
				for _, edge := range obj.Object.InternalEdges {
					var color = rl.White

					if edge.GetFlag().IsBlocked() {
						color = rl.Red
					}
					startPosVec := edge.Line.A.Clone().ApplyMatrix4(obj.GetLocalToWorld())
					endPosVec := edge.Line.B.Clone().ApplyMatrix4(obj.GetLocalToWorld())
					if drawObjectsWithZeroHeight {
						startPosVec.Y = 0
						endPosVec.Y = 0
					}

					startPos := rl.NewVector3(startPosVec.X, startPosVec.Y, startPosVec.Z)
					endPos := rl.NewVector3(endPosVec.X, endPosVec.Y, endPosVec.Z)
					rl.DrawLine3D(startPos, endPos, color)
				}
			}
		}

		drawCollisionPoints()
	}
}

func defaultCam() rl.Camera3D {
	return rl.Camera3D{
		Position:   rl.Vector3{X: 1920 / 2, Y: 1000, Z: 1920 / 2},
		Target:     rl.NewVector3(1920/2, 0, 1920/2),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       120,
		Projection: rl.CameraPerspective,
	}
}
