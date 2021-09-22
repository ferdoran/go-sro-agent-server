package service

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/config"
	"github.com/ferdoran/go-sro-agent-server/engine/geo"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-agent-server/navmesh"
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

type WorldService struct {
	visibleObjects    map[uint32]model.ISRObject
	playersByUniqueId map[uint32]*model.Player
	playersByCharName map[string]*model.Player
	npcsByUniqueId    map[uint32]*model.NPC
	pets              map[uint32]model.ISRObject // TODO Move to own type+
	regions           map[int16]*model.Region
	movingObjects     map[uint32]model.ICharacter
	uniqueIdCounter   uint32
	mutex             *sync.Mutex
	loader            *navmesh.Loader
	loader2           *navmeshv2.Loader
	navmeshGobPath    string
}

var worldServiceInstance *WorldService
var worldServiceOnce sync.Once

func GetWorldServiceInstance() *WorldService {
	worldServiceOnce.Do(func() {
		worldServiceInstance = &WorldService{
			visibleObjects:    make(map[uint32]model.ISRObject),
			playersByUniqueId: make(map[uint32]*model.Player),
			playersByCharName: make(map[string]*model.Player),
			npcsByUniqueId:    make(map[uint32]*model.NPC),
			pets:              make(map[uint32]model.ISRObject),
			regions:           make(map[int16]*model.Region),
			movingObjects:     make(map[uint32]model.ICharacter),
			uniqueIdCounter:   0,
			mutex:             &sync.Mutex{},
			loader:            navmesh.NewLoader(viper.GetString(config.AgentDataPath)),
			loader2:           navmeshv2.NewLoader(viper.GetString(config.AgentDataPath)),
			navmeshGobPath:    viper.GetString(config.AgentPrelinkedNavdataFile),
		}
	})

	return worldServiceInstance
}

func (w *WorldService) GetObjectByUniqueId(objectUniqueId uint32) (model.ISRObject, error) {
	if player, exists := w.visibleObjects[objectUniqueId]; exists {
		return player, nil
	} else {
		return nil, errors.New(fmt.Sprintf("object with uniqueId %d does not exist", objectUniqueId))
	}
}

func (w *WorldService) GetPlayerByUniqueId(playerUniqueId uint32) (*model.Player, error) {
	if player, exists := w.playersByUniqueId[playerUniqueId]; exists {
		return player, nil
	} else {
		return nil, errors.New(fmt.Sprintf("player with uniqueId %d does not exist", playerUniqueId))
	}
}

func (w *WorldService) AddPlayer(p *model.Player) {
	// TODO probably do more checks
	// TODO add visible objects that player can see
	w.AddVisibleObject(p)
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.playersByUniqueId[p.GetUniqueID()] = p
	w.playersByCharName[p.CharName] = p
}

func (w *WorldService) AddVisibleObject(o model.ISRObject) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.visibleObjects[w.uniqueIdCounter] = o
	o.SetUniqueID(w.uniqueIdCounter)
	w.uniqueIdCounter++
	o.GetPosition().Region.AddVisibleObject(o)
}

func (w *WorldService) PlayerDisconnected(uid uint32, charName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, reg := range w.regions {
		for _, obj := range reg.GetVisibleObjects() {
			obj.GetKnownObjectList().RemoveObject(w.playersByCharName[charName])
		}
		reg.RemoveVisibleObject(w.visibleObjects[uid])
	}
	delete(w.visibleObjects, uid)
	delete(w.playersByUniqueId, uid)
	delete(w.playersByCharName, charName)
}

func (w *WorldService) LoadGameServerRegions(gameServerId int) map[int16]*model.Region {
	log.Info("loading game server regions")
	gsRegions := model.GetRegionsForGameServer(gameServerId)

	_, err := os.Stat(w.navmeshGobPath)

	if os.IsNotExist(err) {
		log.Infof("prelinked navdata file does not exist. loading from data then")
		w.loader.LoadNavMeshInfos()
		//w.loader2.LoadNavMeshInfos()
		w.loader.LoadNavMeshData()
		//w.loader2.LoadTerrainMeshes()
	} else {
		log.Infof("loading prelinked navdata file")
		w.loader.LoadPrecomputedNavmeshDataFromGOB(w.navmeshGobPath)
	}

	for _, region := range gsRegions {
		w.AddRegions(region.ContinentName, region.Regions...)
		w.LoadSpawnDataForContinent(region.ContinentName)
	}

	log.Info("finished loading game server regions")
	return w.regions
}

func (w *WorldService) AddRegions(continent string, regions ...int16) {
	log.Debugf("loading regions for %s", continent)
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
		navMeshData := w.loader.NavMeshData[fileName]
		region := model.NewRegionFromNavMeshData(reg, navMeshData)
		region.LinkInternalEdges()
		w.regions[reg] = &region
	}

	log.Debugln("linking global edges")
	for _, reg := range w.regions {
		reg.LinkGlobalEdges(w.regions)
		reg.CalculateObjectMatrices()

	}

	log.Infof("Skipped %d dungeons for %s\n", dungeonsCounter, continent)
	log.Infof("Finished loading regions for %s\n", continent)
}

func (w *WorldService) GetRegion(regionId int16) (*model.Region, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if reg, exists := w.regions[regionId]; exists {
		return reg, nil
	}
	return nil, errors.New(fmt.Sprintf("region does not exist: %d", regionId))
}

func (w *WorldService) GetRegions() map[int16]*model.Region {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.regions
}

func (w *WorldService) LoadSpawnDataForContinent(continent string) {
	spawns := model.GetSpawnsForContinent(continent)
	log.Infof("initialising %d spawns", len(spawns))
	for _, spawn := range spawns {
		reg := w.regions[spawn.RegionID]
		if reg == nil {
			log.Debugf("found spawn for non-existent region: %d", spawn.RegionID)
			continue
		}
		s := model.Spawn{
			Position: model.Position{
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
			MaxTotalCount:  spawn.MaxTotalCount,
			DelayTimeMin:   spawn.DelayTimeMin,
			DelayTimeMax:   spawn.DelayTimeMax,
		}
		reg.Spawns = append(reg.Spawns, model.InitSpawnAreaFromSpawnNest(s))
	}
}

func (w *WorldService) RegisterMovingCharacter(char model.ICharacter) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.movingObjects[char.GetUniqueID()] == nil {
		w.movingObjects[char.GetUniqueID()] = char
	}
}

func (w *WorldService) GetMovingObjects() map[uint32]model.ICharacter {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	newMap := make(map[uint32]model.ICharacter)

	for k, v := range w.movingObjects {
		newMap[k] = v
	}
	return newMap
}

func (w *WorldService) RemoveMovingObject(uid uint32) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.movingObjects, uid)
}

func (w *WorldService) GetWorldCoordinatesForPosition(p model.Position) (float32, float32, float32) {
	rX, rZ := utils.Int16ToXAndZ(p.Region.ID)

	if p.X > model.RegionWidth {
		rX++
		p.X -= model.RegionWidth
	}

	if p.Z > model.RegionHeight {
		rZ++
		p.Z -= model.RegionHeight
	}

	regionId := utils.XAndZToInt16(byte(rX), byte(rZ))
	region, err := w.GetRegion(regionId)

	if err != nil {
		log.Panic(err)
	}

	p.Region = region
	x := (float32(rX) * model.RegionWidth) + p.X
	y := p.Region.GetYAtOffset(p.X, p.Z)
	z := (float32(rZ) * model.RegionHeight) + p.Z
	return x, y, z
}

func (w *WorldService) NewPosFromWorldCoordinates(x, z float32) (model.Position, error) {
	regionX := byte(x / model.RegionWidth)
	regionZ := byte(z / model.RegionHeight)

	regionId := utils.XAndZToInt16(regionX, regionZ)
	region, err := w.GetRegion(regionId)
	if err != nil {
		return model.Position{}, err
	}
	pX := x - float32(regionX)*float32(model.RegionWidth)
	pZ := z - float32(regionZ)*float32(model.RegionHeight)
	pY := region.GetYAtOffset(pX, pZ)

	return model.Position{
		X:       pX,
		Y:       pY,
		Z:       pZ,
		Heading: 0,
		Region:  region,
	}, nil
}

func (w *WorldService) GetKnownObjectsAroundObject(region *model.Region, object model.ISRObject) map[uint32]model.ISRObject {
	regions := w.GetNeighbourRegions(region)
	knownObjects := make(map[uint32]model.ISRObject)
	for _, reg := range regions {
		for _, otherObject := range reg.GetVisibleObjects() {
			if object.GetUniqueID() == otherObject.GetUniqueID() {
				continue
			}

			if w.DistanceTo(object.GetPosition(), otherObject.GetPosition()) <= model.RegionHeight/3 {
				// TODO What about stealth / invisible characters?
				knownObjects[otherObject.GetUniqueID()] = otherObject
			}
		}

	}
	return knownObjects
}

func (w *WorldService) GetNeighbourRegions(r *model.Region) []*model.Region {
	regions := make([]*model.Region, 0)
	x, z := utils.Int16ToXAndZ(r.ID)

	for x1 := x - 1; x1 <= x+1; x1++ {
		for z1 := z - 1; z1 <= z+1; z1++ {
			if reg, _ := w.regions[utils.XAndZToInt16(byte(x1), byte(z1))]; reg != nil {
				regions = append(regions, reg)
			}
		}
	}
	return regions
}

func (w *WorldService) UpdatePosition(p model.ICharacter) bool {
	// TODO implement
	if p.GetMovementData() == nil {
		return true
	}
	currentTime := time.Now()
	movementSpeed := p.GetMovementSpeed()
	deltaTime := currentTime.Sub(p.GetMovementData().UpdateTime)

	if deltaTime <= 0 {
		// Position was just updated
		return false
	}

	curPos := p.GetPosition()
	curWorldX, _, curWorldZ := w.GetWorldCoordinatesForPosition(curPos)
	curWorldVec := math32.NewVector3(curWorldX, 0, curWorldZ)
	var walkVector *math32.Vector3
	nextPosIsTarget := false

	if p.GetMovementData().HasDestination {
		targetWorldX, _, targetWorldZ := w.GetWorldCoordinatesForPosition(p.GetMovementData().TargetPosition)
		targetWorldVec := math32.NewVector3(targetWorldX, 0, targetWorldZ)
		walkVector = targetWorldVec.Clone().Sub(curWorldVec.Clone()).Normalize()
	} else {
		x := math32.Cos(math32.DegToRad(p.GetMovementData().DirectionAngle))
		z := math32.Sin(math32.DegToRad(p.GetMovementData().DirectionAngle))

		walkVector = math32.NewVector3(x, 0, z) // already normalized
	}
	nextPosVec := curWorldVec.Clone().Add(walkVector.MultiplyScalar(movementSpeed * float32(deltaTime.Seconds())))

	newPos, err := w.NewPosFromWorldCoordinates(nextPosVec.X, nextPosVec.Z)

	if err != nil {
		log.Panic(errors.Wrap(err, "failed to calculate new player position"))
	}

	if p.GetMovementData().HasDestination && w.DistanceToSquared(curPos, newPos) >= w.DistanceToSquared(curPos, p.GetMovementData().TargetPosition) {
		newPos = p.GetMovementData().TargetPosition
		nextPosIsTarget = true
	}

	curCell := curPos.Region.GetCellAtOffset(curPos.X, curPos.Z)
	newCell := newPos.Region.GetCellAtOffset(newPos.X, newPos.Z)
	heading := math32.Atan2(walkVector.Z, walkVector.X)
	newPos.Heading = heading

	if curPos.Region.ID != newPos.Region.ID {
		log.Tracef("new position is in new region (%d) -> (%d)\n", curPos.Region.ID, newPos.Region.ID)
		if !curPos.Region.CanEnter(curCell, newCell) {
			p.StopMovement()
			log.Tracef("Cell collision between R(%d)[%d] and R(%d)[%d]\n", curCell.RegionID, curCell.ID, newCell.RegionID, newCell.ID)
			return true
		}
	}
	hasCollision, _, inObj, objPos := geo.FindCollisions(
		math32.NewVector3(curPos.X, curPos.Y, curPos.Z),
		math32.NewVector3(newPos.X, newPos.Y, newPos.Z),
		curPos.Region.ID,
		newPos.Region.ID,
		curPos.Region.Objects,
		newPos.Region.Objects)
	if hasCollision {
		p.StopMovement()
		p.SendPositionUpdate()
		return true
	}

	if inObj && objPos != nil && !geo.IsNextPositionTooHigh(curWorldVec, nextPosVec) {
		newPos.Y = objPos.Y
		log.Tracef("Changing position to obj position: %v", newPos)
		objPos = nil
	}

	if curCell.ID != newCell.ID && !inObj {
		log.Tracef("cell %d has %d objects\n", curCell.ID, curCell.ObjCount)
		if !p.GetPosition().Region.CanEnter(curCell, newCell) {
			p.StopMovement()
			log.Debugf("Cell collision between R(%d)[%d] and R(%d)[%d]\n", curCell.RegionID, curCell.ID, newCell.RegionID, newCell.ID)
			return true
		}
	}
	log.Tracef("setting new position to %v\n", newPos)
	if diff := math32.Abs(p.GetPosition().Y - newPos.Y); diff > 10 {
		log.Tracef("y-pos difference greater 10: %v\n", diff)
	}
	p.SetPosition(newPos)
	if curPos.Region != nil && newPos.Region != nil && curPos.Region.ID != newPos.Region.ID {
		curPos.Region.RemoveVisibleObject(p)
		newPos.Region.AddVisibleObject(p)
	}
	if nextPosIsTarget {
		p.StopMovement()
		return true
	}
	p.GetMovementData().UpdateTime = currentTime
	return false
	// TODO
}

func (w *WorldService) DistanceToSquared(from, to model.Position) float32 {
	x1, y1, z1 := w.GetWorldCoordinatesForPosition(from)
	x2, y2, z2 := w.GetWorldCoordinatesForPosition(to)

	dx := x1 - x2
	dy := y1 - y2
	dz := z1 - z2

	return dx*dx + dy*dy + dz*dz
}

func (w *WorldService) DistanceTo(from, to model.Position) float32 {
	return math32.Sqrt(w.DistanceToSquared(from, to))
}

func (w *WorldService) GetPlayerByCharName(charName string) (*model.Player, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if player, exists := w.playersByCharName[charName]; exists {
		return player, nil
	}

	return nil, errors.New(fmt.Sprintf("player with char name %s does not exist", charName))
}

func (w *WorldService) Broadcast(packet network.Packet) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, player := range w.playersByUniqueId {
		player.GetSession().Conn.Write(packet.ToBytes())
	}
}

func (w *WorldService) BroadcastRaw(message []byte) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, player := range w.playersByUniqueId {
		player.GetSession().Conn.Write(message)
	}
}
