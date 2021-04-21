package manager

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"sync"
	"time"
)

const (
	KnownListUpdateRate          = time.Millisecond * 100
	KnownListManagerInitialDelay = time.Second
)

type KnownListManager struct {
	Manager
}

var knownListManagerInstance *KnownListManager
var knownListManagerOnce sync.Once

func GetKnownListManager() *KnownListManager {
	knownListManagerOnce.Do(func() {
		knownListManagerInstance = &KnownListManager{}
		knownListManagerInstance.mutex = sync.Mutex{}
		knownListManagerInstance.Name = "KnownListManager"
		knownListManagerInstance.rate = KnownListUpdateRate
		knownListManagerInstance.initialDelay = KnownListManagerInitialDelay
		knownListManagerInstance.runnerFunc = knownListManagerInstance.updateKnownLists
	})

	return knownListManagerInstance
}

func (k *KnownListManager) updateKnownLists() {
	world := service.GetWorldServiceInstance()
	for k.started {
		select {
		case <-k.ticker.C:
			for _, region := range world.GetRegions() {
				for _, object := range region.GetVisibleObjects() {
					knownObjects := world.GetKnownObjectsAroundObject(region, object)
					knownObjectsList := object.GetKnownObjectList()

					// Remove unknown objects first
					for uid, unknownObj := range knownObjectsList.GetKnownObjects() {
						if knownObjects[uid] == nil {
							knownObjectsList.RemoveObject(unknownObj)
						}
					}

					// Add new objects
					for _, knownObj := range knownObjects {
						if !knownObjectsList.Knows(knownObj) {
							knownObjectsList.AddObject(knownObj)
						}
					}

				}
			}
		}
	}

	k.mutex.Lock()
	defer k.mutex.Unlock()
	k.started = false
}
