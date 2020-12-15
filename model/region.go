package model

import (
	"github.com/ferdoran/go-sro-agent-server/engine/geo/math"
	"github.com/ferdoran/go-sro-fileutils/navmesh"
	"github.com/ferdoran/go-sro-framework/utils"
	"github.com/g3n/engine/math32"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	RegionWidth  = 1920
	RegionHeight = 1920
)

type Region struct {
	ID                int16
	Cells             []*TerrainCell
	Objects           []*navmesh.Object
	HeightMap         navmesh.HeightMap
	TileMap           navmesh.TileMap
	InternalEdges     navmesh.NavMeshInternalEdges
	GlobalEdges       navmesh.NavMeshGlobalEdges
	RegionToWorld     *math32.Matrix4
	WorldToRegion     *math32.Matrix4
	GlobalEdgeFlags   map[byte]int
	InternalEdgeFlags map[byte]int
	Spawns            []Spawn
	VisibleObjects    map[uint32]ISRObject
	mutex             sync.Mutex
}

func NewRegionFromNavMeshData(id int16, data navmesh.NavMeshData) Region {
	r := Region{
		ID:             id,
		Objects:        data.ObjectList.Objects,
		HeightMap:      data.HeightMap,
		TileMap:        data.TileMap,
		InternalEdges:  data.NavMeshInternalEdges,
		GlobalEdges:    data.NavMeshGlobalEdges,
		Spawns:         make([]Spawn, 0),
		VisibleObjects: make(map[uint32]ISRObject),
	}

	r.Cells = make([]*TerrainCell, 0)
	for idx, cell := range data.NavigationCells.Cells {
		r.Cells = append(r.Cells, &TerrainCell{
			ID:               idx,
			Box2:             math32.NewBox2(&cell.Min, &cell.Max),
			Min:              &cell.Min,
			Max:              &cell.Max,
			LinkedNeighbours: make([]*TerrainCell, 0),
			RegionLinks:      make(map[int16][]*TerrainCell),
			RegionID:         r.ID,
			ObjCount:         cell.ObjCount,
			Objects:          cell.Objects,
		})
	}
	x, z := utils.Int16ToXAndZ(id)
	r.RegionToWorld = math32.NewMatrix4().Compose(
		math32.NewVector3(float32(x), r.HeightMap.Heights[0], float32(z)),
		math32.NewQuaternion(0, 0, 0, 0),
		math32.NewVector3(1, 1, 1),
	)
	r.WorldToRegion = r.RegionToWorld.Clone()
	err := r.WorldToRegion.GetInverse(r.RegionToWorld)
	if err != nil {
		logrus.Panic(err)
	}
	return r
}

func (r *Region) LinkInternalEdges() {
	//var skippedEdges []int
	for _, edge := range r.InternalEdges.InternalEdges {
		sourceCell := r.Cells[edge.AssocCell0]
		if edge.AssocCell1 == 0xFFFF {
			// destination cell unknown
			//skippedEdges = append(skippedEdges, edgeId)
			continue
		}
		destinationCell := r.Cells[edge.AssocCell1]
		sourceCell.LinkedNeighbours = append(sourceCell.LinkedNeighbours, destinationCell)
		destinationCell.LinkedNeighbours = append(destinationCell.LinkedNeighbours, sourceCell)
	}

	//logrus.Tracef("skipped %d edges for region %d\n", len(skippedEdges), r.ID)
}

func (r *Region) LinkGlobalEdges(regions map[int16]*Region) {
	for _, edge := range r.GlobalEdges.GlobalEdges {
		destRegion := regions[int16(edge.AssocRegion1)]
		if destRegion == nil {
			if r.ID == 25000 || r.ID == 25256 || edge.AssocRegion1 == 25000 || edge.AssocRegion1 == 25256 {
				logrus.Debugf("Failed to link %d and %d\n", r.ID, edge.AssocRegion1)
			}
			continue
		}
		if int(edge.AssocCell1) > len(destRegion.Cells) {
			logrus.Tracef("Destination cell does not exist")
		}
		destCell := destRegion.Cells[edge.AssocCell1]

		sourceCell := r.Cells[edge.AssocCell0]
		sourceCell.RegionLinks[destRegion.ID] = append(sourceCell.RegionLinks[destRegion.ID], destCell)
	}
	if r.ID == 25256 {
		c32 := r.Cells[32]
		for _, neighbour := range c32.RegionLinks[25000] {
			logrus.Tracef("R[%d] TerrainCell[%d] is linked with R[%d] TerrainCell[%d]", 25256, c32.ID, 25000, neighbour.ID)
		}
	}
}

func (r *Region) GetCellAtOffset(x, z float32) *TerrainCell {
	tile := r.GetTileAtOffset(int(x), int(z))
	return r.Cells[tile.CellID]
}

func (r *Region) GetTileAtOffset(x, z int) navmesh.Tile {
	xTile := x / 20
	zTile := z / 20

	return r.TileMap.Map[xTile+96*zTile]
}

func (r *Region) GetYAtOffset(x, z float32) float32 {
	tX := int(x / 20)
	tZ := int(z / 20)

	h1 := r.HeightMap.Heights[tX+97*tZ]
	h2 := r.HeightMap.Heights[tX+97*(tZ+1)]
	h3 := r.HeightMap.Heights[(tX+1)+97*tZ]
	h4 := r.HeightMap.Heights[(tX+1)+97*(tZ+1)]

	// h1--------h3
	// |   |      |
	// |   |      |
	// h5--+------h6
	// |   |      |
	// h1--------h3

	tileOffsetX := x - (20 * float32(tX))
	tileOffsetXLength := tileOffsetX / 20
	tileOffsetZ := z - (20 * float32(tZ))
	tileOffsetZLength := tileOffsetZ / 20

	h5 := h1 + (h2-h1)*tileOffsetZLength
	h6 := h3 + (h4-h3)*tileOffsetZLength
	yHeight := h5 + (h6-h5)*tileOffsetXLength

	return yHeight
}

func (r *Region) CanEnter(source, destination *TerrainCell) bool {
	// TODO Would be easier to ray cast and check all touched cells instead of expanding all neighbour nodes
	logrus.Tracef("checking if can enter TerrainCell [%d] from TerrainCell [%d]\n", destination.ID, source.ID)
	if source.ID == destination.ID {
		logrus.Tracef("TerrainCell [%d] and TerrainCell [%d] are same", source.ID, destination.ID)
		return true
	}

	for _, neighbour := range source.LinkedNeighbours {
		logrus.Tracef("TerrainCell[%d] linked to [%d]\n", destination.ID, neighbour.ID)
		if destination.ID == neighbour.ID {
			logrus.Tracef("TerrainCell[%d] and TerrainCell[%d] are neighbours ", source.ID, destination.ID)
			return true
		}
	}

	for _, neighbour := range destination.LinkedNeighbours {
		logrus.Tracef("TerrainCell[%d] linked to [%d]\n", source.ID, neighbour.ID)
		if source.ID == neighbour.ID {
			logrus.Tracef("TerrainCell[%d] and TerrainCell[%d] are neighbours ", source.ID, destination.ID)
			return true
		}
	}

	for _, rLink := range source.RegionLinks[destination.RegionID] {
		if rLink.ID == destination.ID {
			logrus.Tracef("R[%d] TerrainCell[%d] and R[%d] TerrainCell[%d] are neighbours ", source.RegionID, source.ID, destination.RegionID, destination.ID)
			return true
		}
	}

	logrus.Tracef("TerrainCell[%d] (%v, %v) and TerrainCell[%d] (%v, %v) are not linked\n", source.ID, source.Box2, source.Max, destination.ID, destination.Min, destination.Max)
	return false
}

func (r *Region) CalculateObjectMatrices() {
	internalEdgeFlags := make(map[byte]int)
	globalEdgeFlags := make(map[byte]int)
	for _, o := range r.Objects {
		o.Rotation = math.NewQuaternion(-o.Yaw, 0, 0)
		o.LocalToWorld = math32.NewMatrix4().Compose(o.Position, o.Rotation, math32.NewVector3(1, 1, 1))
		o.WorldToLocal = o.LocalToWorld.Clone()

		err := o.WorldToLocal.GetInverse(o.LocalToWorld)
		if err != nil {
			logrus.Panic(err)
		}

		for _, ge := range o.GlobalEdges {
			globalEdgeFlags[ge.Flag]++
		}

		for _, ie := range o.InternalEdges {
			internalEdgeFlags[ie.Flag]++
		}
	}

	r.GlobalEdgeFlags = globalEdgeFlags
	r.InternalEdgeFlags = internalEdgeFlags
}

func (r *Region) GetNeighbourRegions() []*Region {
	regions := make([]*Region, 0)
	x, z := utils.Int16ToXAndZ(r.ID)
	world := GetSroWorldInstance()

	for x1 := x - 1; x1 <= x+1; x1++ {
		for z1 := z - 1; z1 <= z+1; z1++ {
			if reg := world.Regions[utils.XAndZToInt16(byte(x), byte(z))]; reg != nil {
				regions = append(regions, reg)
			}
		}
	}
	return regions
}

func (r *Region) GetKnownObjectsAroundObject(object ISRObject) map[uint32]ISRObject {
	regions := r.GetNeighbourRegions()
	knownObjects := make(map[uint32]ISRObject)
	for _, reg := range regions {
		for _, otherObject := range reg.VisibleObjects {
			if object.GetUniqueID() == otherObject.GetUniqueID() {
				continue
			}

			if object.GetPosition().DistanceTo(otherObject.GetPosition()) <= 1000 {
				// TODO What about stealth / invisible characters?
				knownObjects[otherObject.GetUniqueID()] = otherObject
			}
		}

	}
	return knownObjects
}

func (r *Region) AddVisibleObject(object ISRObject) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, isPlayer := object.(IPlayer)

	if !isPlayer {
		logrus.Tracef("object %s is not of type IPlayer. it got type %s and %T", object.GetName(), object.GetType(), object)
	}
	if r.VisibleObjects[object.GetUniqueID()] == nil {
		r.VisibleObjects[object.GetUniqueID()] = object
	}
}

func (r *Region) RemoveVisibleObject(object ISRObject) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.VisibleObjects[object.GetUniqueID()] != nil {
		delete(r.VisibleObjects, object.GetUniqueID())
	}
}

func (p *Player) GetMovementSpeed() float32 {

	switch p.MotionState {
	case Walking:
		return p.WalkSpeed
	case Running:
		if p.BodyState == Berserk {
			return p.RunSpeed * 2
		}
		return p.RunSpeed
	default:
		return p.RunSpeed
	}
}
