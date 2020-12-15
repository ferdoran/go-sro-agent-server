package character

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
)

type MovementHandler struct {
}

func NewMovementHandler() server.PacketHandler {
	handler := MovementHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.EntityMovementRequest, handler)
	return handler
}

func (mh MovementHandler) Handle(data server.PacketChannelData) {
	hasDestination, err := data.ReadBool()
	if err != nil {
		logrus.Panicf("failed to read movement type")
	}

	world := model.GetSroWorldInstance()
	player := world.PlayersByUniqueId[data.UserContext.UniqueID]

	if hasDestination {
		regionId, err := data.ReadInt16()
		if err != nil {
			logrus.Panicf("failed to read regionId")
		}
		x, err := data.ReadInt16()
		if err != nil {
			logrus.Panicf("failed to read x")
		}
		y, err := data.ReadInt16()
		if err != nil {
			logrus.Panicf("failed to read y")
		}
		z, err := data.ReadInt16()
		if err != nil {
			logrus.Panicf("failed to read z")
		}

		targetPos := model.Position{
			X:       float32(x),
			Y:       float32(y),
			Z:       float32(z),
			Heading: 0,
			Region:  world.Regions[regionId],
		}

		//spawnEngine := spawn.GetSpawnEngineInstance()

		p := network.EmptyPacket()
		p.MessageID = opcode.EntityMovementResponse
		p.WriteUInt32(player.UniqueID)
		p.WriteByte(1)
		p.WriteUInt16(uint16(regionId))
		p.WriteUInt16(uint16(x) + 0xFFFF)
		p.WriteUInt16(uint16(y))
		p.WriteUInt16(uint16(z) + 0xFFFF)
		p.WriteByte(1)
		p.WriteUInt16(uint16(player.Position.Region.ID))
		p.WriteUInt16(uint16(player.Position.X) * 10)
		p.WriteFloat32(player.Position.Y)
		p.WriteUInt16(uint16(player.Position.Z) * 10)

		player.Session.Conn.Write(p.ToBytes())

		//spawnEngine.StartedMoving(player, targetPos)

		player.MoveToPosition(targetPos)

		//walkToDestination(player, targetPos, spawnEngine)

		logrus.Tracef("moving %s to position %d (%d|%d|%d)\n", player.CharName, regionId, x, y, z)
	} else {
		angleAction, err := data.ReadByte()
		if err != nil {
			logrus.Panicf("failed to read bool")
		}
		angle, err := data.ReadUInt16()
		if err != nil {
			logrus.Panicf("failed to read angle")
		}
		logrus.Tracef("MOVEMENT ANGLE %d %d\n", angleAction, angle)
	}
}

//func walkToDestination(player *model.Player, target model.Position, spawnEngine *spawn.SpawnEngine) {
//	if player.MovementTimer != nil {
//		player.MovementTimer.Stop()
//	}
//	if player.MovementTicker != nil {
//		player.MovementTicker.Stop()
//	}
//
//	distance := player.Position.DistanceTo(target)
//	walkTimeInSeconds := distance / player.GetRunSpeed()
//
//	logrus.Tracef("MovementSpeed %f\n", player.GetRunSpeed())
//	logrus.Tracef("Distance %f\n", distance)
//
//	startPosWorldX, _, startPosWorldZ := player.Position.ToWorldCoordinates()
//	endPosWorldX, _, endPosWorldZ := target.ToWorldCoordinates()
//
//	startPos := math32.NewVector3(startPosWorldX, 0, startPosWorldZ)
//	endPos := math32.NewVector3(endPosWorldX, 0, endPosWorldZ)
//	diffPos := endPos.Sub(startPos)
//	angle := math.AngleToEastInDeg(*diffPos)
//	heading := (angle / 360) * 0xFFFF
//	player.MovementTicker = time.NewTicker(time.Second / 60)
//	player.MovementTimer = time.NewTimer(time.Duration(walkTimeInSeconds*1000) * time.Millisecond)
//	go func() {
//		logrus.Tracef("Starting Movement. Walking for %f seconds\n", walkTimeInSeconds)
//		curNumOfTicks := 0
//		maxNumOfTicks := int(walkTimeInSeconds * 60)
//		walkStepVector := diffPos.MultiplyScalar(1 / float32(maxNumOfTicks))
//		world := model.GetSroWorldInstance()
//		for {
//			select {
//			case <-player.MovementTimer.C:
//				player.MovementTicker.Stop()
//				player.MovementTimer.Stop()
//				player.MotionState = model.NoMotion
//				logrus.Tracef("Stopping Movement")
//				break
//			case <-player.MovementTicker.C:
//				curNumOfTicks++
//				if curNumOfTicks <= maxNumOfTicks {
//					player.MotionState = model.Running
//					worldX, _, worldZ := player.Position.ToWorldCoordinates()
//					curPosVec := math32.NewVector3(worldX, player.Position.Y, worldZ)
//
//					newPosVec := curPosVec.Clone().Add(walkStepVector)
//					newPos := model.NewPosFromWorldCoordinates(newPosVec.X, newPosVec.Z)
//					newPos.Heading = heading
//					curCell := player.Position.Region.GetCellAtOffset(player.Position.X, player.Position.Z)
//					newCell := newPos.Region.GetCellAtOffset(newPos.X, newPos.Z)
//
//					objects := make([]*navmesh.Object, 0)
//					for _, objId := range curCell.Objects {
//						obj := player.Position.Region.Objects[objId]
//						objects = append(objects, obj)
//					}
//
//					if curCell.ID != newCell.ID {
//						for _, objId := range newCell.Objects {
//							obj := newPos.Region.Objects[objId]
//							objects = append(objects, obj)
//						}
//					}
//
//					if player.Position.Region.ID != newPos.Region.ID {
//						logrus.Tracef("new position is in new region (%d) -> (%d)\n", player.Position.Region.ID, newPos.Region.ID)
//						if !player.Position.Region.CanEnter(curCell, newCell) {
//							player.MotionState = model.NoMotion
//							player.MovementTicker.Stop()
//							player.MovementTimer.Stop()
//							player.SendPositionUpdate()
//							logrus.Tracef("Cell collision between R(%d)[%d] and R(%d)[%d]\n", curCell.RegionID, curCell.ID, newCell.RegionID, newCell.ID)
//							return
//						}
//					}
//					hasCollision, _, inObj, objPos := geo.FindCollisions(player.Position, newPos)
//					if hasCollision {
//						player.MovementTimer.Stop()
//						player.MovementTicker.Stop()
//						player.MotionState = model.NoMotion
//						player.SetPosition(player.GetPosition())
//						player.SendPositionUpdate()
//						spawnEngine.PositionChanged(player)
//						spawnEngine.PlayerObjectCollision(player)
//						return
//					}
//
//					if inObj && objPos != nil && !IsNextPositionTooHigh(player.GetPosition(), model.Position{
//						X:       objPos.X,
//						Y:       objPos.Y,
//						Z:       objPos.Z,
//						Heading: heading,
//						Region:  world.Regions[newCell.RegionID],
//					}) {
//						newPos.Y = objPos.Y
//						logrus.Tracef("Changing position to obj position: %v", newPos)
//						objPos = nil
//					}
//
//					if curCell.ID != newCell.ID && !inObj {
//						logrus.Tracef("cell %d has %d objects\n", curCell.ID, curCell.ObjCount)
//						if !player.Position.Region.CanEnter(curCell, newCell) {
//							player.MovementTicker.Stop()
//							player.MovementTimer.Stop()
//							player.MotionState = model.NoMotion
//							logrus.Debugf("Cell collision between R(%d)[%d] and R(%d)[%d]\n", curCell.RegionID, curCell.ID, newCell.RegionID, newCell.ID)
//							player.SendPositionUpdate()
//							spawnEngine.PlayerObjectCollision(player)
//							return
//						}
//					}
//					inObj = false
//					objPos = nil
//					logrus.Tracef("setting new position to %v\n", newPos)
//					if diff := math32.Abs(player.GetPosition().Y - newPos.Y); diff > 10 {
//						logrus.Tracef("y-pos difference greater 10: %v\n", diff)
//					}
//					player.SetPosition(newPos)
//					spawnEngine.PositionChanged(player)
//				}
//			}
//
//		}
//	}()
//}

func IsNextPositionTooHigh(curPos, newPos model.Position) bool {
	x, _, z := curPos.ToWorldCoordinates()
	x1, _, z1 := newPos.ToWorldCoordinates()
	v0 := math32.NewVector3(x, 0, z)
	v1 := math32.NewVector3(x1, 0, z1)

	adjacent := v0.DistanceTo(v1)
	opposite := math32.Abs(curPos.Y - newPos.Y)
	angle := math32.RadToDeg(math32.Atan(opposite / adjacent))
	logrus.Tracef("next pos angle is %f°\n", angle)
	return math32.Abs(newPos.Y-curPos.Y) > 100
}
