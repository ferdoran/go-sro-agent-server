package model

import "github.com/ferdoran/go-sro-framework/db"

type Spawn struct {
	Position       Position
	RefObjID       uint64
	NpcCodeName    string
	Radius         int
	GenerateRadius int
}

type SpawnDB struct {
	X              float32
	Y              float32
	Z              float32
	Heading        int
	RegionID       int16
	RefObjID       uint64
	NpcCodeName    string
	Radius         int
	GenerateRadius int
}

const (
	SelectSpawnsForContinent string = `SELECT r.wRegionID, n.fLocalPosX, n.fLocalPosY, n.fLocalPosZ, n.wInitialDir, t.ObjID, c.CodeName, n.nRadius, n.nGenerateRadius 
FROM SRO_SHARD.REGION_REFERENCE r, SRO_SHARD.SPAWN_REF_NESTS n, SRO_SHARD.SPAWN_REF_TACTICS t, SRO_SHARD.CHAR_REF_DATA c
WHERE r.ContinentName=?
AND r.wRegionID = n.nRegionDBID
AND n.dwTacticsID = t.TacticsID
AND t.ObjID = c.RefObjID`
)

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
		var heading, radius, generateRadius int

		err = queryHandle.Scan(&regionId, &x, &y, &z, &heading, &refObjId, &codeName, &radius, &generateRadius)
		db.CheckError(err)

		spawn := SpawnDB{
			RegionID:       regionId,
			X:              x,
			Y:              y,
			Z:              z,
			Heading:        heading,
			RefObjID:       refObjId,
			NpcCodeName:    codeName,
			Radius:         radius,
			GenerateRadius: generateRadius,
		}

		spawns = append(spawns, spawn)
	}
	return spawns
}
