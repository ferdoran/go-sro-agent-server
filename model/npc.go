package model

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/sirupsen/logrus"
	"sync"
)

type NPC struct {
	SRObject
	navmeshv2.RtNavmeshPosition
	Type         string
	Mutex        *sync.Mutex
	UniqueID     uint32
	MovementData *MovementData
	BodyState
	LifeState
	MotionState
	WalkSpeed float32
	RunSpeed  float32
	HwanSpeed float32
}

func (n *NPC) GetPosition() Position {
	return n.Position
}

func (n *NPC) SetPosition(position Position) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.Position = position
	logrus.Tracef("changed position of %s to %v", n.Name, n.Position)
}

func (n *NPC) GetType() string {
	return n.Type
}

func (n *NPC) GetUniqueID() uint32 {
	return n.UniqueID
}

func (n *NPC) SetUniqueID(uniqueId uint32) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.UniqueID = uniqueId
}

func (n *NPC) GetMovementData() *MovementData {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	return n.MovementData
}

func (n *NPC) GetBodyState() BodyState {
	return n.BodyState
}

func (n *NPC) GetLifeState() LifeState {
	return n.LifeState
}

func (n *NPC) GetMotionState() MotionState {
	return n.MotionState
}

func (n *NPC) GetWalkSpeed() float32 {
	return 16
}

func (n *NPC) GetRunSpeed() float32 {
	return 50
}

func (n *NPC) GetHwanSpeed() float32 {
	return 100
}

func (n *NPC) SetBodyState(state BodyState) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.BodyState = state
}

func (n *NPC) SetLifeState(state LifeState) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.LifeState = state
}

func (n *NPC) SetMotionState(state MotionState) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.MotionState = state
}

func (n *NPC) SetName(name string) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.Name = name
}

func (n *NPC) MoveToPosition(position Position) {
	// TODO implement
}

func (n *NPC) UpdatePosition() bool {
	// TODO implement
	return true
}

func (n *NPC) SetWalkSpeed(speed float32) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.WalkSpeed = speed
}

func (n *NPC) SetRunSpeed(speed float32) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.RunSpeed = speed
}

func (n *NPC) SetHwanSpeed(speed float32) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	n.HwanSpeed = speed
}

func (n *NPC) GetMovementSpeed() float32 {
	return n.GetRunSpeed()
}

func (n *NPC) StopMovement() {
	// TODO implement
}

func (n *NPC) SendPositionUpdate() {
	// TODO not required
}
