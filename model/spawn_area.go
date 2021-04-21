package model

import (
	"math/rand"
	"time"
)

type SpawnArea struct {
	Spawn
	NPCs         map[uint32]*NPC
	NPCDied      chan uint32
	RespawnTimes []time.Time
	tickRate     int
}

func InitSpawnAreaFromSpawnNest(spawn Spawn) *SpawnArea {
	area := &SpawnArea{
		Spawn:        spawn,
		NPCs:         make(map[uint32]*NPC),
		NPCDied:      make(chan uint32, 100),
		RespawnTimes: make([]time.Time, 0),
	}

	initialSpawnTime := area.calculateRespawnTime()

	for i := 0; i < spawn.MaxTotalCount; i++ {
		area.RespawnTimes = append(area.RespawnTimes, initialSpawnTime)
	}

	//go area.Respawn()
	return area
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
