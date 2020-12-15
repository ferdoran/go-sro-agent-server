package model

const VarianceStatMask = ^uint64(0) >> (64 - 5)

type WeaponVariance struct {
	Durability   int
	PhyReinforce int
	MagReinforce int
	HitRate      int
	PhyAttack    int
	MagAttack    int
	CriticalRate int
}

type ArmorVariance struct {
	Durability   int
	PhyReinforce int
	MagReinforce int
	PhyDefense   int
	MagDefense   int
	ParryRate    int
}

type ShieldVariance struct {
	Durability   int
	PhyReinforce int
	MagReinforce int
	PhyDefense   int
	MagDefense   int
	BlockRate    int
}

type AccessoryVariance struct {
	PhyAbsorb int
	MagAbsorb int
}

func (wv WeaponVariance) ToVariance() (variance uint64) {
	variance |= uint64(wv.Durability)
	variance <<= 5
	variance |= uint64(wv.PhyReinforce)
	variance <<= 5
	variance |= uint64(wv.MagReinforce)
	variance <<= 5
	variance |= uint64(wv.HitRate)
	variance <<= 5
	variance |= uint64(wv.PhyAttack)
	variance <<= 5
	variance |= uint64(wv.MagAttack)
	variance <<= 5
	variance |= uint64(wv.CriticalRate)
	return
}

func (av ArmorVariance) ToVariance() (variance uint64) {
	variance |= uint64(av.Durability)
	variance <<= 5
	variance |= uint64(av.PhyReinforce)
	variance <<= 5
	variance |= uint64(av.MagReinforce)
	variance <<= 5
	variance |= uint64(av.PhyDefense)
	variance <<= 5
	variance |= uint64(av.MagDefense)
	variance <<= 5
	variance |= uint64(av.ParryRate)
	return
}

func (sv ShieldVariance) ToVariance() (variance uint64) {
	variance |= uint64(sv.Durability)
	variance <<= 5
	variance |= uint64(sv.PhyReinforce)
	variance <<= 5
	variance |= uint64(sv.MagReinforce)
	variance <<= 5
	variance |= uint64(sv.BlockRate)
	variance <<= 5
	variance |= uint64(sv.PhyDefense)
	variance <<= 5
	variance |= uint64(sv.MagDefense)
	return
}

func (av AccessoryVariance) ToVariance() (variance uint64) {
	variance |= uint64(av.PhyAbsorb)
	variance <<= 5
	variance |= uint64(av.MagAbsorb)
	return
}

func WeaponStatsFromVariance(variance uint64) (weaponVariance WeaponVariance) {
	weaponVariance.CriticalRate = int(variance & VarianceStatMask)
	variance >>= 5
	weaponVariance.MagAttack = int(variance & VarianceStatMask)
	variance >>= 5
	weaponVariance.PhyAttack = int(variance & VarianceStatMask)
	variance >>= 5
	weaponVariance.HitRate = int(variance & VarianceStatMask)
	variance >>= 5
	weaponVariance.MagReinforce = int(variance & VarianceStatMask)
	variance >>= 5
	weaponVariance.PhyReinforce = int(variance & VarianceStatMask)
	variance >>= 5
	weaponVariance.Durability = int(variance & VarianceStatMask)
	return
}

func ArmorStatsFromVariance(variance uint64) (armorVariance ArmorVariance) {
	armorVariance.ParryRate = int(variance & VarianceStatMask)
	variance >>= 5
	armorVariance.MagDefense = int(variance & VarianceStatMask)
	variance >>= 5
	armorVariance.PhyDefense = int(variance & VarianceStatMask)
	variance >>= 5
	armorVariance.MagReinforce = int(variance & VarianceStatMask)
	variance >>= 5
	armorVariance.PhyReinforce = int(variance & VarianceStatMask)
	variance >>= 5
	armorVariance.Durability = int(variance & VarianceStatMask)
	return
}

func ShieldStatsFromVariance(variance uint64) (shieldVariance ShieldVariance) {
	shieldVariance.MagDefense = int(variance & VarianceStatMask)
	variance >>= 5
	shieldVariance.PhyDefense = int(variance & VarianceStatMask)
	variance >>= 5
	shieldVariance.BlockRate = int(variance & VarianceStatMask)
	variance >>= 5
	shieldVariance.MagReinforce = int(variance & VarianceStatMask)
	variance >>= 5
	shieldVariance.PhyReinforce = int(variance & VarianceStatMask)
	variance >>= 5
	shieldVariance.Durability = int(variance & VarianceStatMask)
	return
}

func AccessoryStatsFromVariance(variance uint64) (accessoryVariance AccessoryVariance) {
	accessoryVariance.MagAbsorb = int(variance & VarianceStatMask)
	variance >>= 5
	accessoryVariance.PhyAbsorb = int(variance & VarianceStatMask)
	return
}
