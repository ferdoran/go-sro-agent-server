package model

type SpawnHive struct {
	ID                     int
	KeepMonsterCountType   byte
	OverwriteMaxTotalCount int
	MonsterCountPerPC      int
	SpawnIncreaseRate      int
	MaxIncreaseRate        int
	Flag                   byte
	GameWorldID            int
	HatchObjType           byte
	DescString128          string
}
