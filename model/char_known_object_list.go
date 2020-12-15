package model

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type ICharKnownObjectList interface {
	IKnownObjectList
	GetObjectsToSpawn() map[uint32]ISRObject
	GetObjectsToDespawn() map[uint32]ISRObject
}

type CharKnownObjectList struct {
	KnownObjectList
	ObjectsToSpawn   map[uint32]ISRObject
	ObjectsToDespawn map[uint32]ISRObject
	mutex            sync.Mutex
}

func (k *CharKnownObjectList) AddObject(object ISRObject) bool {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	logrus.Debugf("adding object %d to CharKnownObjectList with type %s", object.GetUniqueID(), object.GetType())
	if k.KnownObjectList.AddObject(object) {
		if _, exists := k.ObjectsToDespawn[object.GetUniqueID()]; exists {
			// despawn list should not contain an object that was just added
			delete(k.ObjectsToDespawn, object.GetUniqueID())
		}

		if _, exists := k.ObjectsToSpawn[object.GetUniqueID()]; !exists {
			// only add to spawn list if not yet included
			k.ObjectsToSpawn[object.GetUniqueID()] = object
			return true
		}

	}

	return false
}

func (k *CharKnownObjectList) RemoveObject(object ISRObject) bool {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	logrus.Debugf("removing object %d from CharKnownObjectList with type %s", object.GetUniqueID(), object.GetType())
	if k.KnownObjectList.RemoveObject(object) {
		if _, exists := k.ObjectsToSpawn[object.GetUniqueID()]; exists {
			// spawn list should not contain an object that was just removed
			delete(k.ObjectsToSpawn, object.GetUniqueID())
		}

		if _, exists := k.ObjectsToDespawn[object.GetUniqueID()]; !exists {
			// only add to spawn list if not yet included
			k.ObjectsToDespawn[object.GetUniqueID()] = object
			return true
		}

	}

	return false
}

func (k *CharKnownObjectList) GetObjectsToSpawn() map[uint32]ISRObject {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	m := k.ObjectsToSpawn
	k.ObjectsToSpawn = make(map[uint32]ISRObject)
	return m
}

func (k *CharKnownObjectList) GetObjectsToDespawn() map[uint32]ISRObject {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	m := k.ObjectsToDespawn
	k.ObjectsToDespawn = make(map[uint32]ISRObject)
	return m
}

func NewCharKnownObjectList(owner ISRObject) *CharKnownObjectList {
	return &CharKnownObjectList{
		KnownObjectList: KnownObjectList{
			KnownObjects: make(map[uint32]ISRObject),
			Owner:        owner,
			mutex:        sync.Mutex{},
		},
		ObjectsToSpawn:   make(map[uint32]ISRObject),
		ObjectsToDespawn: make(map[uint32]ISRObject),
		mutex:            sync.Mutex{},
	}
}

func (k *CharKnownObjectList) GetKnownObjects() map[uint32]ISRObject {
	return k.KnownObjects
}
