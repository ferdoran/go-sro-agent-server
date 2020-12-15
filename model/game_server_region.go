package model

import "gitlab.ferdoran.de/game-dev/go-sro/framework/db"

type GameServerRegion struct {
	ContinentName string
	Regions       []int16
}

func (gsr *GameServerRegion) HasRegion(regionId int16) bool {
	for _, region := range gsr.Regions {
		if regionId == region {
			return true
		}
	}
	return false
}

const (
	SelectGameserverRegions string = "SELECT gr.Continent_Name, r.wRegionID  FROM `SRO_SHARD`.`GAME_SERVER_REGION` gr INNER JOIN `SRO_SHARD`.`REGION_REFERENCE` AS r ON r.ContinentName = gr.Continent_Name WHERE gr.Game_Server_ID=?"
)

func GetRegionsForGameServer(gameServerId int) []GameServerRegion {
	conn := db.OpenConnShard()
	defer conn.Close()

	queryHandle, err := conn.Query(SelectGameserverRegions, gameServerId)
	db.CheckError(err)

	regions := make(map[string][]int16)
	for queryHandle.Next() {
		var continent string
		var regionId int16
		err = queryHandle.Scan(&continent, &regionId)
		db.CheckError(err)

		regions[continent] = append(regions[continent], regionId)
	}

	var gsRegions []GameServerRegion
	for k, v := range regions {
		gsRegions = append(gsRegions, GameServerRegion{
			ContinentName: k,
			Regions:       v,
		})
	}
	return gsRegions
}
