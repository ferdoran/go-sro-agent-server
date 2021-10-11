package manager

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/g3n/engine/math32"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

type RespawnManager struct {
	Manager
}

var respawnManagerInstance *RespawnManager
var respawnManagerOnce sync.Once

func GetRespawnManagerInstance() *RespawnManager {
	respawnManagerOnce.Do(func() {
		respawnManagerInstance = &RespawnManager{}
		respawnManagerInstance.Name = "RespawnManager"
		respawnManagerInstance.initialDelay = time.Second
		respawnManagerInstance.rate = time.Millisecond * 100
		respawnManagerInstance.runnerFunc = respawnManagerInstance.respawn
	})
	return respawnManagerInstance
}

func (s *RespawnManager) respawn() {
	world := service.GetWorldServiceInstance()
	for s.started {
		select {
		case <-s.ticker.C:
			for _, region := range world.GetRegions() {
				for _, s := range world.GetSpawnsForRegion(region.Region.ID) {
					//spawnArea.Respawn()
					//for deadNpc := range s.NPCDied {
					//	if _, ok := s.NPCs[deadNpc]; ok {
					//		delete(s.NPCs, deadNpc)
					//		respawnTime := s.calculateRespawnTime()
					//		s.respawnTimes = append(s.respawnTimes, respawnTime)
					//	}
					//}
					logrus.Tracef("respawning %s on position R[%d] (%f|%f|%f)", s.NpcCodeName, s.Position.Region.ID, s.Position.Offset.X, s.Position.Offset.Y, s.Position.Offset.Z)
					currentTime := time.Now()
					//if len(s.respawnTimes) > 0 && len(s.NPCs) < s.MaxTotalCount {
					for i, respawnTime := range s.RespawnTimes {
						if currentTime.After(respawnTime) {
							// TODO respawn
							spawnNpc(s)
							s.RespawnTimes = removeArrayElement(s.RespawnTimes, i)
						}
					}
					//}
				}
			}
		}
	}
}

func spawnNpc(s *model.SpawnArea) {
	npc := &model.NPC{
		Type:  "NPC",
		Mutex: &sync.Mutex{},
	}
	npc.RtNavmeshPosition = generateRandomPositionInRadius(s, 0)
	npc.KnownObjectList = model.NewKnownObjectList(npc)
	npc.Name = s.NpcCodeName
	npc.RefObjectID = uint32(s.RefObjID)
	refChar, err := service.GetReferenceDataServiceInstance().GetReferenceCharacter(npc.RefObjectID)
	if err != nil {
		logrus.Error("failed to spawn npc: ref char %d has no type info")
		return
	}
	npc.TypeInfo = refChar.TypeInfo
	service.GetWorldServiceInstance().AddVisibleObject(npc)
	s.NPCs[npc.GetUniqueID()] = npc
}

func generateRandomPositionInRadius(s *model.SpawnArea, retryCount int) navmeshv2.RtNavmeshPosition {
	vWorld := s.Position.GetGlobalCoordinates()
	xWorld := int32(vWorld.X)
	zWorld := int32(vWorld.Z)

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

	newPos, err := service.GetWorldServiceInstance().NewPosFromGlobalCoordinates(math32.NewVector3(float32(spawnXWorld), 0, float32(spawnZWorld)))
	if err != nil {
		logrus.Warn(errors.Wrap(err, fmt.Sprintf("failed to generate random position (%d times), for mob spawn near %s", retryCount, s.Position.String())))
		return generateRandomPositionInRadius(s, retryCount+1)
	}
	newPos.Offset.Y = s.Position.Offset.Y

	return newPos
}

func removeArrayElement(array []time.Time, index int) []time.Time {
	if index >= len(array) {
		return array
	}
	array[len(array)-1], array[index] = array[index], array[len(array)-1]
	return array[:len(array)-1]
}
