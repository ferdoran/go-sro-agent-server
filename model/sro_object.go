package model

import "sync"

type ISRObject interface {
	GetType() string
	GetTypeInfo() TypeInfo
	SetTypeInfo(info TypeInfo)
	SetPosition(pos Position)
	GetPosition() Position
	GetUniqueID() uint32
	SetUniqueID(uniqueId uint32)
	GetRefObjectID() uint32
	GetKnownObjectList() IKnownObjectList
	GetName() string
	GetMovementData() *MovementData
}

type SRObject struct {
	Position
	TypeInfo
	KnownObjectList IKnownObjectList
	UniqueID        uint32
	RefObjectID     uint32
	Type            string
	RWMutex         sync.RWMutex
	Name            string
}

func (o *SRObject) GetRefObjectID() uint32 {
	return o.RefObjectID
}

func (o *SRObject) GetPosition() (position Position) {
	o.RWMutex.RLock()
	defer o.RWMutex.RUnlock()
	position = o.Position
	return
}

func (o *SRObject) SetPosition(newPosition Position) {
	o.RWMutex.Lock()
	defer o.RWMutex.Unlock()

	//if oldReg := o.GetPosition().Region; oldReg != nil && oldReg.ID != newPosition.Region.ID {
	//	o.Region.RemoveVisibleObject(o)
	//}
	o.Position = newPosition
	//o.Region.AddVisibleObject(o)
}

func (o *SRObject) GetUniqueID() uint32 {
	return o.UniqueID
}

func (o *SRObject) SetUniqueID(uniqueID uint32) {
	o.RWMutex.Lock()
	defer o.RWMutex.Unlock()
	o.UniqueID = uniqueID
}

func (o *SRObject) GetType() string {
	return "SRObject"
}

func (o *SRObject) GetTypeInfo() TypeInfo {
	return o.TypeInfo
}

func (o *SRObject) SetTypeInfo(info TypeInfo) {
	o.RWMutex.Lock()
	defer o.RWMutex.Unlock()
	o.TypeInfo = info
}

func (o *SRObject) GetKnownObjectList() IKnownObjectList {
	o.RWMutex.RLock()
	defer o.RWMutex.RUnlock()
	return o.KnownObjectList
}

func (o *SRObject) GetName() string {
	return o.Name
}
