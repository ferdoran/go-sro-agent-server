package manager

import (
	"github.com/ferdoran/go-sro-agent-server/config"
	"github.com/ferdoran/go-sro-agent-server/model"
	"sync"
	"time"
)

type GameTimeManager struct {
	Manager
	referenceTime time.Time
}

var gameTimeManagerInstance *GameTimeManager
var gameTimeManagerOnce sync.Once

func GetGameTimeManagerInstance() *GameTimeManager {
	gameTimeManagerOnce.Do(func() {
		now := time.Now()
		gameTimeManagerInstance = &GameTimeManager{
			referenceTime: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		}
		gameTimeManagerInstance.Name = "GameTimeManager"
		gameTimeManagerInstance.initialDelay = time.Second
		gameTimeManagerInstance.rate = time.Second / time.Duration(config.GlobalConfig.GameTimeConfig.TicksPerSecond)
		gameTimeManagerInstance.runnerFunc = gameTimeManagerInstance.moveObjects
	})
	return gameTimeManagerInstance
}

func (gtm *GameTimeManager) moveObjects() {
	for gtm.started {
		world := model.GetSroWorldInstance()
		for _, obj := range world.GetMovingObjects() {
			if obj.UpdatePosition() {
				world.RemoveMovingObject(obj.GetUniqueID())
			}
		}
	}
}

func (gtm *GameTimeManager) GetCurrentTick() int64 {
	return int64(time.Now().Sub(gtm.referenceTime) / gtm.rate)
}
