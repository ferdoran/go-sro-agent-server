package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/fileutils/navmesh"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/utils"
	"strings"
	"sync"
)

type SroWorld struct {
	VisibleObjects    map[uint32]ISRObject
	PlayersByUniqueId map[uint32]*Player
	PlayersByCharName map[string]*Player
	NPCsByUniqueId    map[uint32]*NPC
	Pets              map[uint32]*ISRObject // TODO Move to own type+
	Regions           map[int16]*Region
	MovingObjects     map[uint32]ICharacter
	uniqueIdCounter   uint32
	mutex             *sync.Mutex
	Loader            *navmesh.Loader
	NavmeshGobPath    string
}

var sroWorldInstance *SroWorld
var sroWorldOnce sync.Once

func InitSroWorldInstance(dataPath, navmeshGobPath string) *SroWorld {
	sroWorldOnce.Do(func() {
		sroWorldInstance = &SroWorld{
			VisibleObjects:    make(map[uint32]ISRObject),
			PlayersByUniqueId: make(map[uint32]*Player),
			PlayersByCharName: make(map[string]*Player),
			Pets:              make(map[uint32]*ISRObject),
			MovingObjects:     make(map[uint32]ICharacter),
			uniqueIdCounter:   100_000,
			mutex:             &sync.Mutex{},
			Regions:           make(map[int16]*Region),
			Loader:            navmesh.NewLoader(dataPath),
			NavmeshGobPath:    navmeshGobPath,
		}
	})
	return sroWorldInstance
}

func GetSroWorldInstance() *SroWorld {
	return sroWorldInstance
}

func (w *SroWorld) AddPlayer(p *Player) {
	// TODO probably do more checks
	// TODO add visible objects that player can see
	w.AddVisibleObject(p)
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.PlayersByUniqueId[p.GetUniqueID()] = p
	w.PlayersByCharName[p.CharName] = p
}

func (w *SroWorld) AddVisibleObject(o ISRObject) {
	w.mutex.Lock()
	w.VisibleObjects[w.uniqueIdCounter] = o
	o.SetUniqueID(w.uniqueIdCounter)
	w.uniqueIdCounter++
	o.GetPosition().Region.AddVisibleObject(o)
	w.mutex.Unlock()
}

func (w *SroWorld) PlayerDisconnected(uid uint32, charName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, reg := range w.Regions {
		for _, obj := range reg.VisibleObjects {
			obj.GetKnownObjectList().RemoveObject(w.PlayersByCharName[charName])
		}
		reg.RemoveVisibleObject(w.VisibleObjects[uid])
	}
	delete(w.VisibleObjects, uid)
	delete(w.PlayersByUniqueId, uid)
	delete(w.PlayersByCharName, charName)
}

func (w *SroWorld) InitiallySpawnAllNpcs() {
	// TODO
	log.Info("spawning NPCs")
	for _, region := range w.Regions {
		for _, spawn := range region.Spawns {
			if strings.Contains(spawn.NpcCodeName, "FORTRESS") {
				continue
			}
			npc := &NPC{
				Type:  "NPC",
				Mutex: &sync.Mutex{},
			}
			npc.Position = spawn.Position
			npc.KnownObjectList = NewKnownObjectList(npc)
			npc.Name = spawn.NpcCodeName
			npc.RefObjectID = uint32(spawn.RefObjID)
			npc.TypeInfo = RefChars[npc.RefObjectID].TypeInfo
			w.AddVisibleObject(npc)
		}
	}
	log.Info("finished spawning NPCs")
}

func (w *SroWorld) LoadGameServerRegions(gameServerId int) map[int16]*Region {
	gsRegions := GetRegionsForGameServer(gameServerId)
	w.Loader.LoadNavMeshInfos()
	w.Loader.LoadPrecomputedNavmeshDataFromGOB(w.NavmeshGobPath)
	//w.Loader.LoadNavMeshData()
	//w.Loader.SaveNavmeshDataAsGOB(w.NavmeshGobPath)
	//w.Loader.SaveNavmeshDataAsJSON()
	for _, region := range gsRegions {
		w.AddRegions(region.ContinentName, region.Regions...)
		w.LoadSpawnDataForContinent(region.ContinentName)
	}
	GetSroWorldInstance().Regions = w.Regions

	return w.Regions
}

func (w *SroWorld) AddRegions(continent string, regions ...int16) {
	utils.PrintSection("Loading regions for " + continent)
	numRegions := len(regions)
	dungeonsCounter := 0
	for i, reg := range regions {
		if reg < 0 {
			// TODO Load dungeon file
			dungeonsCounter++
			continue
		}
		x, z := utils.Int16ToXAndZ(reg)
		log.Tracef("Loading region %02x%02x, X|Z (%d|%d)", z, x, x, z)
		fileName := fmt.Sprintf("nv_%02x%02x.nvm", z, x)
		utils.PrintProgress(i+1, numRegions)
		navMeshData := w.Loader.NavMeshData[fileName]
		region := NewRegionFromNavMeshData(reg, navMeshData)
		region.LinkInternalEdges()
		w.Regions[reg] = &region
	}

	log.Debugln("Linking global edges")
	for _, reg := range w.Regions {
		reg.LinkGlobalEdges(w.Regions)
		reg.CalculateObjectMatrices()

	}

	log.Infof("Skipped %d dungeons for %s\n", dungeonsCounter, continent)
	log.Infof("Finished loading regions for %s\n", continent)
}

func (w *SroWorld) LoadSpawnDataForContinent(continent string) {
	spawns := GetSpawnsForContinent(continent)
	for _, spawn := range spawns {
		reg := w.Regions[spawn.RegionID]
		if reg == nil {
			log.Debugf("found spawn for non-existent region: %d", spawn.RegionID)
			continue
		}
		s := Spawn{
			Position: Position{
				X:       spawn.X,
				Y:       spawn.Y,
				Z:       spawn.Z,
				Heading: float32(spawn.Heading),
				Region:  reg,
			},
			RefObjID:       spawn.RefObjID,
			NpcCodeName:    spawn.NpcCodeName,
			Radius:         spawn.Radius,
			GenerateRadius: spawn.GenerateRadius,
		}
		reg.Spawns = append(reg.Spawns, s)
	}
}

func (w *SroWorld) RegisterMovingCharacter(char ICharacter) {

	if w.MovingObjects[char.GetUniqueID()] == nil {
		w.mutex.Lock()
		defer w.mutex.Unlock()
		w.MovingObjects[char.GetUniqueID()] = char
	}
}

func (w *SroWorld) GetMovingObjects() map[uint32]ICharacter {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.MovingObjects
}

func (w *SroWorld) RemoveMovingObject(uid uint32) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.MovingObjects, uid)
}
