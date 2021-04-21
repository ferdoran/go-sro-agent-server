package service

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"sync"
	"time"
)

type MovementService struct {
}

var movementServiceInstance *MovementService
var movementServiceOnce sync.Once

func GetMovementServiceInstance() *MovementService {
	worldServiceOnce.Do(func() {
		movementServiceInstance = &MovementService{}
	})

	return movementServiceInstance
}

func (m *MovementService) MoveToPosition(player *model.Player, newPosition model.Position) {
	// TODO implement
	currentTime := time.Now()
	movementData := &model.MovementData{
		StartTime:      currentTime,
		UpdateTime:     currentTime,
		TargetPosition: newPosition,
		HasDestination: true,
		DirectionAngle: 0,
	}
	pPos := player.GetPosition()
	packet := network.EmptyPacket()
	packet.MessageID = opcode.EntityMovementResponse
	packet.WriteUInt32(player.GetUniqueID())
	packet.WriteBool(movementData.HasDestination)
	packet.WriteUInt16(uint16(movementData.TargetPosition.Region.ID))
	packet.WriteUInt16(uint16(movementData.TargetPosition.X) + 0xFFFF)
	packet.WriteUInt16(uint16(movementData.TargetPosition.Y))
	packet.WriteUInt16(uint16(movementData.TargetPosition.Z) + 0xFFFF)
	packet.WriteByte(1)
	packet.WriteUInt16(uint16(pPos.Region.ID))
	packet.WriteUInt16(uint16(pPos.X) * 10)
	packet.WriteFloat32(pPos.Y)
	packet.WriteUInt16(uint16(pPos.Z) * 10)

	player.MovementData = movementData
	player.SetMotionState(model.Running) // TODO check if walking or sitting
	worldServiceInstance.RegisterMovingCharacter(player)
	// Broadcast movement Update to known objects around

	player.Broadcast(&packet)
}

func (m *MovementService) WalkToDirection(player *model.Player, heading float32) {
	currentTime := time.Now()
	movementData := &model.MovementData{
		StartTime:      currentTime,
		UpdateTime:     currentTime,
		HasDestination: false,
		DirectionAngle: heading,
	}

	player.MovementData = movementData
	player.SetMotionState(model.Running) // TODO check if walking or sitting
	worldServiceInstance.RegisterMovingCharacter(player)
}
