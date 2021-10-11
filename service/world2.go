package service

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/config"
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

var ErrRegionNotExisting = errors.New("region does not exist")

type WorldService struct {
	visibleObjects         map[uint32]model.ISRObject
	playersByUniqueId      map[uint32]*model.Player
	playersByCharName      map[string]*model.Player
	npcsByUniqueId         map[uint32]*model.NPC
	pets                   map[uint32]model.ISRObject // TODO Move to own type+
	regions                map[int16]navmeshv2.RtNavmeshTerrain
	movingObjects          map[uint32]model.ICharacter
	movementData           map[uint32]*model.MovementData
	visibleObjectsByRegion map[int16][]model.ISRObject
	uniqueIdCounter        uint32
	spawns                 map[int16][]*model.SpawnArea
	mutex                  *sync.Mutex
	loader                 *navmesh.Loader
	loader2                *navmeshv2.Loader
	navmeshGobPath         string
}

var worldServiceInstance *WorldService
var worldServiceOnce sync.Once

func GetWorldServiceInstance() *WorldService {
	worldServiceOnce.Do(func() {
		worldServiceInstance = &WorldService{
			visibleObjects:         make(map[uint32]model.ISRObject),
			playersByUniqueId:      make(map[uint32]*model.Player),
			playersByCharName:      make(map[string]*model.Player),
			npcsByUniqueId:         make(map[uint32]*model.NPC),
			pets:                   make(map[uint32]model.ISRObject),
			regions:                make(map[int16]navmeshv2.RtNavmeshTerrain),
			movingObjects:          make(map[uint32]model.ICharacter),
			movementData:           make(map[uint32]*model.MovementData),
			visibleObjectsByRegion: make(map[int16][]model.ISRObject),
			spawns:                 make(map[int16][]*model.SpawnArea),
			uniqueIdCounter:        0,
			mutex:                  &sync.Mutex{},
			loader:                 navmesh.NewLoader(viper.GetString(config.AgentDataPath)),
			loader2:                navmeshv2.NewLoader(viper.GetString(config.AgentDataPath)),
			navmeshGobPath:         viper.GetString(config.AgentPrelinkedNavdataFile),
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
	w.visibleObjectsByRegion[o.GetNavmeshPosition().Region.ID] = append(w.visibleObjectsByRegion[o.GetNavmeshPosition().Region.ID], o)
}

func (w *WorldService) PlayerDisconnected(uid uint32, charName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, reg := range w.regions {
		for i, obj := range w.visibleObjectsByRegion[reg.Region.ID] {
			if obj.GetUniqueID() == uid {
				s := w.visibleObjectsByRegion[reg.Region.ID]
				s[i] = s[len(s)-1]
				s = s[:len(s)-1]
			}
			obj.GetKnownObjectList().RemoveObject(w.playersByCharName[charName])

		}
	}
	delete(w.visibleObjects, uid)
	delete(w.playersByUniqueId, uid)
	delete(w.playersByCharName, charName)
	delete(w.movingObjects, uid)
}

func (w *WorldService) LoadGameServerRegions(gameServerId int) map[int16]navmeshv2.RtNavmeshTerrain {
	log.Info("loading game server regions")
	gsRegions := model.GetRegionsForGameServer(gameServerId)

	_, err := os.Stat(w.navmeshGobPath)

	if os.IsNotExist(err) {
		log.Infof("prelinked navdata file does not exist. loading from data then")
		w.loader2.LoadNavMeshInfos()
		w.loader2.LoadTerrainMeshes(nil)
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
		utils.PrintProgress(i+1, numRegions)
		w.regions[reg] = w.loader2.RegionData[reg]
	}

	log.Infof("Skipped %d dungeons for %s\n", dungeonsCounter, continent)
	log.Infof("Finished loading regions for %s\n", continent)
}

func (w *WorldService) GetRegion(regionId int16) (navmeshv2.RtNavmeshTerrain, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if reg, exists := w.regions[regionId]; exists {
		return reg, nil
	}
	return navmeshv2.RtNavmeshTerrain{}, errors.New(fmt.Sprintf("region does not exist: %d", regionId))
}

func (w *WorldService) GetRegions() map[int16]navmeshv2.RtNavmeshTerrain {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.regions
}

func (w *WorldService) GetSpawnsForRegion(regionId int16) []*model.SpawnArea {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.spawns[regionId]
}

func (w *WorldService) LoadSpawnDataForContinent(continent string) {
	spawns := model.GetSpawnsForContinent(continent)
	log.Infof("initialising %d spawns", len(spawns))
	for _, spawn := range spawns {
		reg, exists := w.regions[spawn.RegionID]
		if !exists {
			log.Debugf("found spawn for non-existent region: %d", spawn.RegionID)
			continue
		}
		spawnPos := math32.NewVector3(spawn.X, spawn.Y, spawn.Z)
		cell, err := reg.ResolveCell(spawnPos)
		if err != nil {
			log.Error(err)
		}
		s := model.Spawn{
			Position: navmeshv2.RtNavmeshPosition{
				Cell:     &cell,
				Instance: nil,
				Region:   reg.Region,
				Offset:   spawnPos,
			},
			RefObjID:       spawn.RefObjID,
			NpcCodeName:    spawn.NpcCodeName,
			Radius:         spawn.Radius,
			GenerateRadius: spawn.GenerateRadius,
			MaxTotalCount:  spawn.MaxTotalCount,
			DelayTimeMin:   spawn.DelayTimeMin,
			DelayTimeMax:   spawn.DelayTimeMax,
		}
		w.spawns[reg.Region.ID] = append(w.spawns[reg.Region.ID], model.InitSpawnAreaFromSpawnNest(s))
	}
}

func (w *WorldService) RegisterMovingCharacter(char model.ICharacter, movementData model.MovementData) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.movingObjects[char.GetUniqueID()] == nil {
		w.movingObjects[char.GetUniqueID()] = char
		w.movementData[char.GetUniqueID()] = &movementData
		log.Debugf("registered moving character %d", char.GetUniqueID())
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
	delete(w.movementData, uid)
	log.Debugf("Removing moving object %d", uid)
}

func (w *WorldService) NewPosFromGlobalCoordinates(vPos *math32.Vector3) (navmeshv2.RtNavmeshPosition, error) {
	regionX := byte(vPos.X / navmeshv2.RegionWidth)
	regionZ := byte(vPos.Z / navmeshv2.RegionHeight)

	regionId := utils.XAndZToInt16(regionX, regionZ)
	region, err := w.GetRegion(regionId)
	if err != nil {
		return navmeshv2.RtNavmeshPosition{}, err
	}

	pX := vPos.X - float32(regionX)*float32(navmeshv2.RegionWidth)
	pZ := vPos.Z - float32(regionZ)*float32(navmeshv2.RegionHeight)
	pos := math32.NewVector3(pX, 0, pZ)
	cell, err := region.ResolveCell(pos)
	if err != nil {
		return navmeshv2.RtNavmeshPosition{}, err
	}

	return navmeshv2.RtNavmeshPosition{
		Cell:     &cell,
		Instance: nil,
		Region:   region.Region,
		Offset:   pos,
	}, nil
}

func (w *WorldService) GetVisibleObjectsForRegion(regionId int16) []model.ISRObject {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.visibleObjectsByRegion[regionId]
}

func (w *WorldService) GetKnownObjectsAroundObject(object model.ISRObject) map[uint32]model.ISRObject {
	// TODO switch to neighbour planes instead of neighbour regions
	regions := w.GetNeighbourRegions(object.GetNavmeshPosition().Region.ID)
	knownObjects := make(map[uint32]model.ISRObject)
	for _, reg := range regions {
		for _, otherObject := range w.GetVisibleObjectsForRegion(reg.Region.ID) {
			if object.GetUniqueID() == otherObject.GetUniqueID() {
				continue
			}

			if w.DistanceTo(object.GetNavmeshPosition(), otherObject.GetNavmeshPosition()) <= model.RegionHeight/3 {
				// TODO What about stealth / invisible characters?
				knownObjects[otherObject.GetUniqueID()] = otherObject
			}
		}

	}
	return knownObjects
}

func (w *WorldService) GetNeighbourPlanes(position navmeshv2.RtNavmeshPosition) {
	_, err := w.GetRegion(position.Region.ID)
	if err != nil {
		log.Panic(err)
	}
	// TODO implement

}

func (w *WorldService) GetNeighbourRegions(regionId int16) []navmeshv2.RtNavmeshTerrain {
	regions := make([]navmeshv2.RtNavmeshTerrain, 0)
	x, z := utils.Int16ToXAndZ(regionId)

	for x1 := x - 1; x1 <= x+1; x1++ {
		for z1 := z - 1; z1 <= z+1; z1++ {
			if reg, exists := w.regions[utils.XAndZToInt16(byte(x1), byte(z1))]; exists {
				regions = append(regions, reg)
			}
		}
	}
	return regions
}

func (w *WorldService) UpdatePosition(p model.ICharacter) bool {
	// TODO implement
	movementData := w.GetMovementData(p)
	if movementData == nil {
		log.Debugf("update position called without movementData")
		return true
	}
	log.Debugf("updating position for %d", p.GetUniqueID())

	currentTime := time.Now()
	movementSpeed := p.GetMovementSpeed()
	deltaTime := currentTime.Sub(movementData.UpdateTime)

	if deltaTime <= 0 {
		// Position was just updated
		return false
	}

	curPos := p.GetNavmeshPosition()
	curWorldVec := curPos.GetGlobalCoordinates()
	var walkVector *math32.Vector3
	nextPosIsTarget := false

	if movementData != nil && movementData.HasDestination {
		targetWorldVec := movementData.TargetPosition.GetGlobalCoordinates()
		walkVector = targetWorldVec.Clone().Sub(curWorldVec.Clone()).Normalize()
	} else {
		x := math32.Cos(math32.DegToRad(p.GetMovementData().DirectionAngle))
		z := math32.Sin(math32.DegToRad(p.GetMovementData().DirectionAngle))

		walkVector = math32.NewVector3(x, 0, z) // already normalized
	}
	nextPosVec := curWorldVec.Clone().Add(walkVector.MultiplyScalar(movementSpeed * float32(deltaTime.Seconds())))

	newPos, err := w.NewPosFromGlobalCoordinates(nextPosVec)

	if err != nil {
		log.Panic(errors.Wrap(err, "failed to calculate new player position"))
	}

	if movementData.HasDestination && w.DistanceToSquared(curPos, newPos) >= w.DistanceToSquared(curPos, movementData.TargetPosition) {
		newPos = movementData.TargetPosition
		nextPosIsTarget = true
	}
	currentTerrain, err := w.GetRegion(curPos.Region.ID)
	if err != nil {
		log.Panic(err)
	}
	curCell, err := currentTerrain.ResolveCell(curPos.Offset)
	if err != nil {
		log.Panic(err)
	}
	nextTerrain, err := w.GetRegion(newPos.Region.ID)
	if err != nil {
		log.Panic(err)
	}
	nextCell, err := nextTerrain.ResolveCell(newPos.Offset)
	if err != nil {
		log.Panic(err)
	}
	heading := math32.Atan2(walkVector.Z, walkVector.X)
	newPos.Heading = heading

	hasCollision, collision, inObject, objectPosition := navmeshv2.FindObjectCollisions(curPos, newPos, currentTerrain.Objects, nextTerrain.Objects)
	if hasCollision {
		newPos.Region = collision.Region
		newPos.Offset = collision.VectorGlobal
		if err != nil {
			log.Panic(err)
		}
		log.Debugf("object collision at %s with edge flag %d", newPos.String(), collision.ContactEdge.GetFlag())
		//p.SetNavmeshPosition(newPos)
		p.SendPositionUpdate()
		return true
	}

	if inObject && objectPosition != nil {
		newPos.Offset = objectPosition
		objectPosition = nil
	}

	if curCell.Index != nextCell.Index && !inObject {
		hasCollision, collision = navmeshv2.FindTerrainCollisions(curPos, newPos, curCell, nextCell)
		if hasCollision {
			newPos.Region = collision.Region
			newPos.Offset = collision.VectorLocal
			newReg, err := w.GetRegion(newPos.Region.ID)
			if err != nil {
				log.Panic(err)
			}
			log.Debugf("object collision at %s with edge flag %d", newPos.String(), collision.ContactEdge.GetFlag())
			newPos.Offset.Y = newReg.ResolveHeight(newPos.Offset)
			p.SetNavmeshPosition(newPos)
			p.SendPositionUpdate()
			return true
		}
	}

	// Check terrain collisions

	if !inObject {
		newPos.Offset.Y = nextTerrain.ResolveHeight(newPos.Offset)
	}
	p.SetNavmeshPosition(newPos)
	if curPos.Region.ID != newPos.Region.ID {
		w.RemoveVisibleObjectFromRegion(p, curPos.Region.ID)
		w.AddVisibleObjectToRegion(p, newPos.Region.ID)
	}

	if nextPosIsTarget {
		return true
	}

	movementData.UpdateTime = currentTime
	return false
	// TODO
}

func (w *WorldService) AddVisibleObjectToRegion(p model.ICharacter, regionId int16) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	alreadyInList := false
	for _, obj := range w.visibleObjectsByRegion[regionId] {
		if obj.GetUniqueID() == p.GetUniqueID() {
			alreadyInList = true
		}
	}
	if !alreadyInList {
		w.visibleObjectsByRegion[regionId] = append(w.visibleObjectsByRegion[regionId], p)
	}
}

func (w *WorldService) RemoveVisibleObjectFromRegion(p model.ICharacter, regionId int16) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for i, obj := range w.visibleObjectsByRegion[regionId] {
		if obj.GetUniqueID() == p.GetUniqueID() {
			s := w.visibleObjectsByRegion[regionId]
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}
}

func (w *WorldService) GetMovementData(p model.ICharacter) *model.MovementData {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.movementData[p.GetUniqueID()]
}

func (w *WorldService) DistanceToSquared(from, to navmeshv2.RtNavmeshPosition) float32 {
	vFromGlobal := from.GetGlobalCoordinates()
	vToGlobal := to.GetGlobalCoordinates()

	dx := vFromGlobal.X - vToGlobal.X
	dy := vFromGlobal.Y - vToGlobal.Y
	dz := vFromGlobal.Z - vToGlobal.Z

	return dx*dx + dy*dy + dz*dz
}

func (w *WorldService) DistanceTo(from, to navmeshv2.RtNavmeshPosition) float32 {
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
