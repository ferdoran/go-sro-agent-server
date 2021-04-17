package model

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

type SpawnArea struct {
	Spawn
	NPCs         map[uint32]*NPC
	NPCDied      chan uint32
	respawnTimes []time.Time
	tickRate     int
}

func InitSpawnAreaFromSpawnNest(spawn Spawn) *SpawnArea {
	area := &SpawnArea{
		Spawn:        spawn,
		NPCs:         make(map[uint32]*NPC),
		NPCDied:      make(chan uint32, 100),
		respawnTimes: make([]time.Time, 0),
	}

	initialSpawnTime := area.calculateRespawnTime()

	for i := 0; i < spawn.MaxTotalCount; i++ {
		area.respawnTimes = append(area.respawnTimes, initialSpawnTime)
	}

	//go area.Respawn()
	return area
}

func (s *SpawnArea) Respawn() {
	//for deadNpc := range s.NPCDied {
	//	if _, ok := s.NPCs[deadNpc]; ok {
	//		delete(s.NPCs, deadNpc)
	//		respawnTime := s.calculateRespawnTime()
	//		s.respawnTimes = append(s.respawnTimes, respawnTime)
	//	}
	//}
	logrus.Tracef("respawning %s on position R[%d] (%f|%f|%f)", s.NpcCodeName, s.Position.Region.ID, s.Position.X, s.Position.Y, s.Position.Z)
	currentTime := time.Now()
	//if len(s.respawnTimes) > 0 && len(s.NPCs) < s.MaxTotalCount {
	for i, respawnTime := range s.respawnTimes {
		if currentTime.After(respawnTime) {
			// TODO respawn
			s.spawnNpc()
			s.respawnTimes = removeArrayElement(s.respawnTimes, i)
		}
	}
	//}
}

func (s *SpawnArea) spawnNpc() {
	npc := &NPC{
		Type:  "NPC",
		Mutex: &sync.Mutex{},
	}
	npc.Position = s.generateRandomPositionInRadius()
	npc.KnownObjectList = NewKnownObjectList(npc)
	npc.Name = s.NpcCodeName
	npc.RefObjectID = uint32(s.RefObjID)
	npc.TypeInfo = RefChars[npc.RefObjectID].TypeInfo

	GetSroWorldInstance().AddVisibleObject(npc)
	s.NPCs[npc.GetUniqueID()] = npc
}

func (s *SpawnArea) generateRandomPositionInRadius() Position {
	xWorld, _, zWorld := s.Position.ToWorldCoordinatesInt32()

	var spawnXWorld, spawnZWorld int32

	minXWorld := xWorld - int32(s.GenerateRadius)
	maxXWorld := xWorld + int32(s.GenerateRadius)
	xWorldDiff := maxXWorld - minXWorld

	minZWorld := zWorld - int32(s.GenerateRadius)
	maxZWorld := zWorld + int32(s.GenerateRadius)
	zWorldDiff := maxZWorld - minZWorld

	if xWorldDiff > 0 {
		spawnXWorld = rand.Int31n(xWorldDiff) + minXWorld
	} else if s.GenerateRadius > 0 {
		spawnXWorld = rand.Int31n(int32(s.GenerateRadius*2)) + minXWorld
	} else {
		spawnXWorld = xWorld
	}

	if zWorldDiff > 0 {
		spawnZWorld = rand.Int31n(zWorldDiff) + minZWorld
	} else if s.GenerateRadius > 0 {
		spawnZWorld = rand.Int31n(int32(s.GenerateRadius*2)) + minZWorld
	} else {
		spawnZWorld = zWorld
	}

	newPos, err := NewPosFromWorldCoordinates(float32(spawnXWorld), float32(spawnZWorld))
	newPos.Y = s.Position.Y
	if err != nil {
		newPos = s.Position
		logrus.Warn(errors.Wrap(err, "failed to generate random position for mob spawn near "+s.Position.String()))
	}

	return newPos
}

func (s *SpawnArea) calculateRespawnTime() time.Time {
	deathTime := time.Now()
	delayDiff := s.DelayTimeMax - s.DelayTimeMin
	var respawnDelay time.Duration
	if delayDiff > 0 {
		respawnDelay = time.Second * time.Duration((rand.Intn(delayDiff))+s.DelayTimeMin)
	} else {
		respawnDelay = time.Second
	}

	return deathTime.Add(respawnDelay)
}

func removeArrayElement(array []time.Time, index int) []time.Time {
	if index >= len(array) {
		return array
	}
	array[len(array)-1], array[index] = array[index], array[len(array)-1]
	return array[:len(array)-1]
}
