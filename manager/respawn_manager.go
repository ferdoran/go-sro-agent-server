package manager

import (
	"github.com/ferdoran/go-sro-agent-server/model"
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
	world := model.GetSroWorldInstance()
	for s.started {
		select {
		case <-s.ticker.C:
			for _, region := range world.GetRegions() {
				for _, spawnArea := range region.Spawns {
					spawnArea.Respawn()
				}
			}
		}
	}
}
