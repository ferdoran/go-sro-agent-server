package model

type BaseStats struct {
	HP  int
	MP  int
	Str int
	Int int
}

type AttackStats struct {
	PhyAttackMin int
	PhyAttackMax int
	MagAttackMin int
	MagAttackMax int
	HitRate      int
	CriticalRate int
}

type DefenseStats struct {
	PhyDef    int
	MagDef    int
	BlockRate int
	ParryRate int
}

type BonusBaseStats struct {
	BaseStats
	HpPercent int
	MpPercent int
}

type BonusAttackStats struct {
	AttackStats
	PhyAttackPercent    int
	MagAttackPercent    int
	HitRatePercent      int
	CriticalRatePercent int
}

type BonusDefenseStats struct {
	DefenseStats
	PhyDefPercent    int
	MagDefPercent    int
	ParryRatePercent int
	BlockRatePercent int
}
