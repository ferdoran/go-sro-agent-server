package model

type Weapon struct {
	Equipment
	WeaponStats
	DurabilityLower float32
	DurabilityUpper float32
}

type WeaponType int

const (
	Sword WeaponType = iota
	Blade
	Spear
	Glaive
	Bow
	TwoHandedSword
	OneHandedSword
	DualAxe
	WarlockRod
	Staff
	Crossbow
	Dagger
	Harp
	ClericRod
	Siege
)

func (w *Weapon) GetWeaponType() WeaponType {
	if w.TypeID1 == 3 && w.TypeID2 == 1 && w.TypeID3 == 6 {
		switch w.TypeID4 {
		case 2:
			return Sword
		case 3:
			return Blade
		case 4:
			return Spear
		case 5:
			return Glaive
		case 6:
			return Bow
		case 7:
			return OneHandedSword
		case 8:
			return TwoHandedSword
		case 9:
			return DualAxe
		case 10:
			return WarlockRod
		case 11:
			return Staff
		case 12:
			return Crossbow
		case 13:
			return Dagger
		case 14:
			return Harp
		case 15:
			return ClericRod
		case 16:
			return Siege
		default:
			return -1
		}
	}
	return -1
}
