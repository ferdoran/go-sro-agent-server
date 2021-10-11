package manager

import (
	"github.com/ferdoran/go-sro-agent-server/config"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/spf13/viper"
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
		gameTimeManagerInstance.rate = time.Second / time.Duration(viper.GetInt(config.GameTimeTicksPerSecond))
		gameTimeManagerInstance.runnerFunc = gameTimeManagerInstance.moveObjects
	})
	return gameTimeManagerInstance
}

func (gtm *GameTimeManager) moveObjects() {
	world := service.GetWorldServiceInstance()
	for gtm.started {
		select {
		case <-gtm.ticker.C:
			for _, obj := range world.GetMovingObjects() {
				if world.UpdatePosition(obj) {
					world.RemoveMovingObject(obj.GetUniqueID())
				}

				// update known object list when position has changed
				knownObjects := world.GetKnownObjectsAroundObject(obj)
				knownObjectsList := obj.GetKnownObjectList()

				// Remove unknown objects first
				// TODO: shouldn't it be removed from the unknownObj too?
				for uid, unknownObj := range knownObjectsList.GetKnownObjects() {
					if knownObjects[uid] == nil {
						knownObjectsList.RemoveObject(unknownObj)
						unknownObj.GetKnownObjectList().RemoveObject(obj)
					}
				}

				// Add new objects
				// TODO: shouldn't it be added to the new objects known list too?
				for _, knownObj := range knownObjects {
					if !knownObjectsList.Knows(knownObj) {
						knownObjectsList.AddObject(knownObj)
						knownObj.GetKnownObjectList().AddObject(obj)
					}
				}
			}
		}
	}
}

func (gtm *GameTimeManager) GetCurrentTick() int64 {
	return int64(time.Now().Sub(gtm.referenceTime) / gtm.rate)
}
