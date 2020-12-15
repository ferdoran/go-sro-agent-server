package model

type SpawnNest struct {
	ID                    int
	HiveID                int
	TacticsID             int
	Position              Position
	Radius                int
	GenerateRadius        int
	ChampionGenPercentage int
	DelayTimeMin          int
	DelayTimeMax          int
	MaxTotalCount         int
	HasFlag               bool
	CanRespawn            bool
	Type                  byte
}
