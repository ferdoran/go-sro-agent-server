package model

import (
	"github.com/ferdoran/go-sro-agent-server/engine/geo"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const PlayerType = "Player"

// Player Body State
type BodyState int

const (
	NoStatus BodyState = iota
	Hwan
	Untouchable
	GmInvincible
	GmInvisible
	Berserk
	Stealth
	Invisible
)

// Life state
type LifeState byte

const (
	Spawning LifeState = iota
	Alive
	Dead
)

// MotionState
type MotionState int

const (
	NoMotion MotionState = iota
	_
	Walking
	Running
	Sitting
)

// PVP State
type PKState int

const (
	PvpWhite PKState = iota
	PvpPurple
	PvpRed
)

// PVP Flag (For CTF)
type CTFTeam int

const (
	CtfRed CTFTeam = iota
	CtfBlue
	CtfNone = 0xFF
)

// PVP Cape
type PVPCape int

const (
	PvpCapeNone PVPCape = iota
	PvpCapeRed
	PvpCapeGray
	PvpCapeBlue
	PvpCapeWhite
	PvpCapeGold
)

type IPlayer interface {
	ICharacter
	GetPKState() PKState
	GetInventory() Inventory
	GetSession() *server.Session
	GetScale() byte
	GetPVPCape() PVPCape
	GetCharKnownObjectList() *CharKnownObjectList
}

type Player struct {
	SRObject
	MotionState
	LifeState
	BodyState
	CharKnownObjectList *CharKnownObjectList
	MovementData        *MovementData
	Session             *server.Session
	Mutex               sync.Mutex
	Scale               byte
	ID                  int
	CharName            string
	Inventory           Inventory
	BaseStats           BaseStats
	BaseAttackStats     AttackStats
	BaseDefenseStats    DefenseStats
	BonusStats          BonusBaseStats
	BonusAttackStats    BonusAttackStats
	BonusDefenseStats   BonusDefenseStats
	PhyAbsorbPercent    int
	MagAbsorbPercent    int
	PhyBalancePercent   int
	MagBalancePercent   int
	SkillPoints         int
	StatPoints          int
	ExpOffset           uint64
	SkillExpOffset      uint
	Level               int
	MaxLevel            int
	PKState             PKState
	CTFTeam             CTFTeam
	PVPCape             PVPCape
	IsGm                bool
	Party               Party
	WalkSpeed           float32
	RunSpeed            float32
	HwanSpeed           float32

	/* TODO
	- Teleport Position
	- Gold
	- Skills
	- Masteries
	- JobInfo
	- TeleportLocation
	- Active Buffs / Effects
	*/
}

func (p *Player) GetLifeState() LifeState {
	return p.LifeState
}

func (p *Player) SetLifeState(state LifeState) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.LifeState = state
}

func (p *Player) GetBodyState() BodyState {
	return p.BodyState
}

func (p *Player) SetBodyState(state BodyState) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.BodyState = state
}

func (p *Player) GetWalkSpeed() float32 {
	return p.WalkSpeed
}

func (p *Player) SetWalkSpeed(speed float32) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.WalkSpeed = speed
}

func (p *Player) GetRunSpeed() float32 {
	return p.RunSpeed
}

func (p *Player) SetRunSpeed(speed float32) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.RunSpeed = speed
}

func (p *Player) GetHwanSpeed() float32 {
	return p.HwanSpeed
}

func (p *Player) SetHwanSpeed(speed float32) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.HwanSpeed = speed
}

func (p *Player) SetName(name string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.Name = name
}

func (p *Player) SetUniqueID(uniqueId uint32) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.UniqueID = uniqueId
}

func (p *Player) SendPositionUpdate() {
	packet := network.EmptyPacket()
	packet.MessageID = opcode.MovementPositionUpdate
	packet.WriteUInt32(p.UniqueID)
	packet.WriteUInt16(uint16(p.Position.Region.ID))
	packet.WriteFloat32(p.Position.X)
	packet.WriteFloat32(p.Position.Y)
	packet.WriteFloat32(p.Position.Z)
	packet.WriteUInt16(uint16(p.Position.Heading))
	p.Session.Conn.Write(packet.ToBytes())
}

func (p *Player) SendStatsUpdate() {
	// Player update stats
	// TODO: actually calculate value
	packet := network.EmptyPacket()
	packet.MessageID = opcode.EntityUpdateStats
	packet.WriteUInt32(uint32(utils.BaseMinAttack(p.BaseStats.Str))) // PhyAttackMin
	packet.WriteUInt32(uint32(utils.BaseMaxAttack(p.BaseStats.Str))) // PhyAttackMax
	packet.WriteUInt32(uint32(utils.BaseMinAttack(p.BaseStats.Int))) // MagAttackMin
	packet.WriteUInt32(uint32(utils.BaseMaxAttack(p.BaseStats.Int))) // MagAttackMax
	packet.WriteUInt16(uint16(utils.BaseDef(p.BaseStats.Str)))       // Phy Def
	packet.WriteUInt16(uint16(utils.BaseDef(p.BaseStats.Int)))       // Mag Def
	packet.WriteUInt16(11)                                           // Hit Rate
	packet.WriteUInt16(11)                                           // Parry Rate
	packet.WriteUInt32(uint32(p.BaseStats.HP))                       // HP
	packet.WriteUInt32(uint32(p.BaseStats.MP))                       // MP
	packet.WriteUInt16(uint16(p.BaseStats.Str))                      // Str
	packet.WriteUInt16(uint16(p.BaseStats.Int))                      // Int
	p.Session.Conn.Write(packet.ToBytes())
}

func (p *Player) SendMovementStateUpdate() {
	packet := network.EmptyPacket()
	packet.MessageID = opcode.EntityUpdateMovementState
	packet.WriteUInt32(p.UniqueID)
	packet.WriteByte(4)
	packet.WriteByte(2)
	p.Session.Conn.Write(packet.ToBytes())
}

func (p *Player) SendEquipItemPacket(item Item, slot byte) {
	packet := network.EmptyPacket()
	packet.MessageID = opcode.EntityEquipItem
	packet.WriteUInt32(p.UniqueID)
	packet.WriteByte(slot)
	packet.WriteUInt32(item.GetRefObjectID())
	packet.WriteBool(item.IsOneHandedWeapon())

	p.Session.Conn.Write(packet.ToBytes())
}

func (p *Player) SendUnequipItemPacket(item Item, slot byte) {
	packet := network.EmptyPacket()
	packet.MessageID = opcode.EntityUnequipItem
	packet.WriteUInt32(p.UniqueID)
	packet.WriteByte(slot)
	packet.WriteUInt32(item.GetRefObjectID())

	p.Session.Conn.Write(packet.ToBytes())
}

func (p *Player) IsChinese() bool {
	return p.RefObjectID >= 1907 && p.RefObjectID <= 1932
}

func (p *Player) IsEuropean() bool {
	return p.RefObjectID >= 14875 && p.RefObjectID <= 14900
}

func (p *Player) IsMale() bool {
	return (p.RefObjectID >= 1907 && p.RefObjectID <= 1919) || (p.RefObjectID >= 14875 && p.RefObjectID <= 14887)
}

func (p *Player) IsFemale() bool {
	return (p.RefObjectID >= 1920 && p.RefObjectID <= 1932) || (p.RefObjectID >= 14888 && p.RefObjectID <= 14900)
}

func (p *Player) GetPosition() Position {
	return p.Position
}

func (p *Player) GetUniqueID() uint32 {
	return p.UniqueID
}

func (p *Player) GetPKState() PKState {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	s := p.PKState
	return s
}

func (p *Player) GetInventory() Inventory {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	i := p.Inventory
	return i
}

func (p *Player) GetName() string {
	return p.CharName
}

func (p *Player) GetSession() *server.Session {
	return p.Session
}

func (p *Player) UpdatePosition() bool {
	// TODO implement
	if p.MovementData == nil {
		return true
	}
	currentTime := time.Now()
	movementSpeed := p.GetMovementSpeed()
	deltaTime := currentTime.Sub(p.MovementData.UpdateTime)

	if deltaTime <= 0 {
		// Position was just updated
		return false
	}

	curPos := p.GetPosition()
	curWorldX, _, curWorldZ := curPos.ToWorldCoordinates()
	curWorldVec := math32.NewVector3(curWorldX, 0, curWorldZ)
	var walkVector *math32.Vector3
	nextPosIsTarget := false

	if p.MovementData.HasDestination {
		targetWorldX, _, targetWorldZ := p.MovementData.TargetPosition.ToWorldCoordinates()
		targetWorldVec := math32.NewVector3(targetWorldX, 0, targetWorldZ)
		walkVector = targetWorldVec.Clone().Sub(curWorldVec.Clone()).Normalize()
	} else {
		x := math32.Cos(math32.DegToRad(p.MovementData.DirectionAngle))
		z := math32.Sin(math32.DegToRad(p.MovementData.DirectionAngle))

		walkVector = math32.NewVector3(x, 0, z) // already normalized
	}
	nextPosVec := curWorldVec.Clone().Add(walkVector.MultiplyScalar(movementSpeed * float32(deltaTime.Seconds())))

	newPos := NewPosFromWorldCoordinates(nextPosVec.X, nextPosVec.Z)

	if p.MovementData.HasDestination && curPos.DistanceToSquared(newPos) > curPos.DistanceToSquared(p.MovementData.TargetPosition) {
		newPos = p.MovementData.TargetPosition
		nextPosIsTarget = true
	}

	curCell := curPos.Region.GetCellAtOffset(curPos.X, curPos.Z)
	newCell := newPos.Region.GetCellAtOffset(newPos.X, newPos.Z)
	heading := math32.Atan2(walkVector.Z, walkVector.X)
	newPos.Heading = heading

	if curPos.Region.ID != newPos.Region.ID {
		logrus.Tracef("new position is in new region (%d) -> (%d)\n", curPos.Region.ID, newPos.Region.ID)
		if !curPos.Region.CanEnter(curCell, newCell) {
			p.StopMovement()
			logrus.Tracef("Cell collision between R(%d)[%d] and R(%d)[%d]\n", curCell.RegionID, curCell.ID, newCell.RegionID, newCell.ID)
			return true
		}
	}
	hasCollision, _, inObj, objPos := geo.FindCollisions(
		math32.NewVector3(curPos.X, curPos.Y, curPos.Z),
		math32.NewVector3(newPos.X, newPos.Y, newPos.Z),
		curPos.Region.ID,
		newPos.Region.ID,
		curPos.Region.Objects,
		newPos.Region.Objects)
	if hasCollision {
		p.StopMovement()
		return true
	}

	if inObj && objPos != nil && !geo.IsNextPositionTooHigh(curWorldVec, nextPosVec) {
		newPos.Y = objPos.Y
		logrus.Tracef("Changing position to obj position: %v", newPos)
		objPos = nil
	}

	if curCell.ID != newCell.ID && !inObj {
		logrus.Tracef("cell %d has %d objects\n", curCell.ID, curCell.ObjCount)
		if !p.Position.Region.CanEnter(curCell, newCell) {
			p.StopMovement()
			logrus.Debugf("Cell collision between R(%d)[%d] and R(%d)[%d]\n", curCell.RegionID, curCell.ID, newCell.RegionID, newCell.ID)
			return true
		}
	}
	logrus.Tracef("setting new position to %v\n", newPos)
	if diff := math32.Abs(p.GetPosition().Y - newPos.Y); diff > 10 {
		logrus.Tracef("y-pos difference greater 10: %v\n", diff)
	}
	p.SetPosition(newPos)
	if curPos.Region != nil && newPos.Region != nil && curPos.Region.ID != newPos.Region.ID {
		curPos.Region.RemoveVisibleObject(p)
		newPos.Region.AddVisibleObject(p)
	}
	if nextPosIsTarget {
		p.StopMovement()
		return true
	}
	p.MovementData.UpdateTime = currentTime
	return false
	// TODO
}

func (p *Player) MoveToPosition(newPosition Position) {
	// TODO implement
	currentTime := time.Now()
	movementData := &MovementData{
		StartTime:      currentTime,
		UpdateTime:     currentTime,
		TargetPosition: newPosition,
		HasDestination: true,
		DirectionAngle: 0,
	}
	pPos := p.GetPosition()
	packet := network.EmptyPacket()
	packet.MessageID = opcode.EntityMovementResponse
	packet.WriteUInt32(p.UniqueID)
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

	p.MovementData = movementData
	p.SetMotionState(Running) // TODO check if walking or sitting
	sroWorldInstance.RegisterMovingCharacter(p)
	// Broadcast movement Update to known objects around

	p.Broadcast(&packet)
}

func (p *Player) WalkToDirection(heading float32) {
	currentTime := time.Now()
	movementData := &MovementData{
		StartTime:      currentTime,
		UpdateTime:     currentTime,
		HasDestination: false,
		DirectionAngle: heading,
	}

	p.MovementData = movementData
	p.SetMotionState(Running) // TODO check if walking or sitting
	sroWorldInstance.RegisterMovingCharacter(p)
}

func (p *Player) Broadcast(packet *network.Packet) {
	packetBuffer := packet.ToBytes()
	playerCount := 0
	for _, object := range p.CharKnownObjectList.GetKnownObjects() {
		if player, isPlayer := object.(IPlayer); isPlayer {
			playerCount++
			player.GetSession().Conn.Write(packetBuffer)
		}
	}
	logrus.Debugf("broadcasted message %02X to %d known players of %s", packet.MessageID, playerCount, p.GetName())
}

func (p *Player) StopMovement() {
	if p.MovementData != nil {
		p.MovementData = nil
	}
	p.SetMotionState(NoMotion)
}

func (p *Player) GetCharKnownObjectList() *CharKnownObjectList {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	return p.CharKnownObjectList
}

func (p *Player) GetKnownObjectList() IKnownObjectList {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	return p.CharKnownObjectList
}

func (p *Player) SetMotionState(newState MotionState) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.MotionState = newState
}

func (p *Player) GetMotionState() MotionState {
	return p.MotionState
}

func (p *Player) GetScale() byte {
	return p.Scale
}

func (p *Player) GetPVPCape() PVPCape {
	return p.PVPCape
}

func (p *Player) GetMovementData() *MovementData {
	return p.MovementData
}
