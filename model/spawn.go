package model

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"github.com/ferdoran/go-sro-framework/db"
)

type Spawn struct {
	NestID         int
	Position       navmeshv2.RtNavmeshPosition
	RefObjID       uint64
	NpcCodeName    string
	Radius         int
	GenerateRadius int
	MaxTotalCount  int
	DelayTimeMin   int
	DelayTimeMax   int
}

type SpawnDB struct {
	NestID         int
	X              float32
	Y              float32
	Z              float32
	Heading        int
	RegionID       int16
	RefObjID       uint64
	NpcCodeName    string
	Radius         int
	GenerateRadius int
	MaxTotalCount  int
	DelayTimeMin   int
	DelayTimeMax   int
}

const SelectSpawnsForContinent string = `SELECT n.dwNestID, r.wRegionID, n.fLocalPosX, n.fLocalPosY, n.fLocalPosZ, n.wInitialDir, t.ObjID, c.CodeName, n.nRadius, n.nGenerateRadius, n.dwMaxTotalCount, n.dwDelayTimeMin, n.dwDelayTimeMax
FROM SRO_SHARD.REGION_REFERENCE r, SRO_SHARD.SPAWN_REF_NESTS n, SRO_SHARD.SPAWN_REF_TACTICS t, SRO_SHARD.CHAR_REF_DATA c
WHERE r.ContinentName=?
AND r.wRegionID = n.nRegionDBID
AND n.dwTacticsID = t.TacticsID
AND t.ObjID = c.RefObjID`

func GetSpawnsForContinent(continent string) []SpawnDB {
	conn := db.OpenConnShard()
	defer conn.Close()

	queryHandle, err := conn.Query(SelectSpawnsForContinent, continent)
	db.CheckError(err)

	spawns := make([]SpawnDB, 0)

	for queryHandle.Next() {
		var regionId int16
		var x, y, z float32
		var refObjId uint64
		var codeName string
		var nestId, heading, radius, generateRadius, maxTotalCount, delayTimeMin, delayTimeMax int

		err = queryHandle.Scan(&nestId, &regionId, &x, &y, &z, &heading, &refObjId, &codeName, &radius, &generateRadius, &maxTotalCount, &delayTimeMin, &delayTimeMax)
		db.CheckError(err)

		spawn := SpawnDB{
			NestID:         nestId,
			RegionID:       regionId,
			X:              x,
			Y:              y,
			Z:              z,
			Heading:        heading,
			RefObjID:       refObjId,
			NpcCodeName:    codeName,
			Radius:         radius,
			GenerateRadius: generateRadius,
			MaxTotalCount:  maxTotalCount,
			DelayTimeMin:   delayTimeMin,
			DelayTimeMax:   delayTimeMax,
		}

		spawns = append(spawns, spawn)
	}
	return spawns
}
