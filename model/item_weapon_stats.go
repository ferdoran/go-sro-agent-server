package model

type WeaponStats struct {
	MinPhyLower float32
	MinPhyUpper float32
	MaxPhyLower float32
	MaxPhyUpper float32

	MinPhyReinforcementLower float32
	MinPhyReinforcementUpper float32
	MaxPhyReinforcementLower float32
	MaxPhyReinforcementUpper float32

	MinMagLower float32
	MinMagUpper float32
	MaxMagLower float32
	MaxMagUpper float32

	MinMagReinforcementLower float32
	MinMagReinforcementUpper float32
	MaxMagReinforcementLower float32
	MaxMagReinforcementUpper float32

	AttackRateLower float32
	AttackRateUpper float32

	CritLower int
	CritUpper int

	DurabilityLower float32
	DurabilityUpper float32
}
