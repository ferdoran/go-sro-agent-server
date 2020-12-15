package manager

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SpawnManager struct {
	Manager
}

var spawnManagerInstance *SpawnManager
var spawnManagerOnce sync.Once

func GetSpawnManagerInstance() *SpawnManager {
	spawnManagerOnce.Do(func() {
		spawnManagerInstance = &SpawnManager{}
		spawnManagerInstance.Name = "SpawnManager"
		spawnManagerInstance.initialDelay = time.Second
		spawnManagerInstance.rate = time.Millisecond * 100
		spawnManagerInstance.runnerFunc = spawnManagerInstance.updateSpawns
	})
	return spawnManagerInstance
}

func (s *SpawnManager) updateSpawns() {
	for s.started {
		select {
		case <-s.ticker.C:
			for _, region := range model.GetSroWorldInstance().Regions {
				for _, object := range region.VisibleObjects {

					if player, isPlayer := object.(model.IPlayer); isPlayer {
						objectsToDespawn := player.GetCharKnownObjectList().GetObjectsToDespawn()
						if len(objectsToDespawn) > 0 {
							groupSpawnObjects(player, objectsToDespawn, false)
						}

						objectsToSpawn := player.GetCharKnownObjectList().GetObjectsToSpawn()
						if len(objectsToSpawn) > 0 {
							groupSpawnObjects(player, objectsToSpawn, true)
						}

					} else {
						logrus.Tracef("object %d is not a player. it is a %s, %T with name %s", object.GetUniqueID(), object.GetType(), object, object.GetName())
					}
				}
			}
		}
	}
}

func groupSpawnObjects(player model.IPlayer, objects map[uint32]model.ISRObject, spawning bool) {
	p1 := network.EmptyPacket()
	p1.MessageID = opcode.EntityGroupSpawnBegin
	if spawning {
		p1.WriteByte(1)
	} else {
		p1.WriteByte(2)
	}
	p1.WriteUInt16(uint16(len(objects)))

	player.GetSession().Conn.Write(p1.ToBytes())
	p2 := network.EmptyPacket()
	p2.MessageID = opcode.EntityGroupSpawnData

	for _, object := range objects {
		if spawning {
			model.WriteEntitySpawnData(&p2, object)
		} else {
			p2.WriteUInt32(object.GetUniqueID())
		}
	}
	player.GetSession().Conn.Write(p2.ToBytes())

	p3 := network.EmptyPacket()
	p3.MessageID = opcode.EntityGroupSpawnEnd
	player.GetSession().Conn.Write(p3.ToBytes())
}
