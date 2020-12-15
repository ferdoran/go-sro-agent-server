package model

type Armour struct {
	Equipment
	ArmorStats
	DurabilityLower float32
	DurabilityUpper float32
}

type ArmorType int

const (
	Garment ArmorType = iota
	Protector
	Armor
	Robe
	LightArmor
	HeavyArmor
)

type ArmorPiece int

const (
	Head ArmorPiece = iota
	Shoulder
	Body
	Legs
	Arms
	Foot
)

func (a *Armour) GetArmorType() ArmorType {
	if a.TypeID1 == 3 && a.TypeID2 == 1 {
		switch a.TypeID3 {
		case 1:
			return Garment
		case 2:
			return Protector
		case 3:
			return Armor
		case 9:
			return Robe
		case 10:
			return LightArmor
		case 11:
			return HeavyArmor
		default:
			return -1
		}
	}
	return -1
}

func (a *Armour) GetArmorPiece() ArmorPiece {
	if a.TypeID1 == 3 && a.TypeID2 == 1 && (a.TypeID3 == 1 || a.TypeID3 == 2 || a.TypeID3 == 3 || a.TypeID3 == 9 || a.TypeID3 == 10 || a.TypeID3 == 11) {
		switch a.TypeID4 {
		case 1:
			return Head
		case 2:
			return Shoulder
		case 3:
			return Body
		case 4:
			return Legs
		case 5:
			return Arms
		case 6:
			return Foot
		default:
			return -1
		}
	}
	return -1
}
