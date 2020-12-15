package model

type Accessory struct {
	Equipment
	AccessoryStats
}

type AccessoryType int

const (
	Chinese AccessoryType = iota
	European
)

type AccessoryPiece int

const (
	Earring AccessoryPiece = iota
	Necklace
	Ring
)

func (a *Accessory) GetAccessoryType() AccessoryType {
	if a.TypeID1 == 3 && a.TypeID2 == 1 {
		switch a.TypeID3 {
		case 5:
			return Chinese
		case 12:
			return European
		default:
			return -1
		}
	}
	return -1
}

func (a *Accessory) GetAccessoryPiece() AccessoryPiece {
	if a.TypeID1 == 3 && a.TypeID2 == 1 && (a.TypeID3 == 5 || a.TypeID3 == 12) {
		switch a.TypeID4 {
		case 1:
			return Earring
		case 2:
			return Necklace
		case 3:
			return Ring
		default:
			return -1
		}
	}
	return -1
}
