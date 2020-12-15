package model

type Shield struct {
	Equipment
	ShieldStats
}

type ShieldType int

const (
	ChineseShield ShieldType = iota
	EuropeanShield
)

func (s *Shield) GetShieldType() ShieldType {
	if s.TypeID1 == 3 && s.TypeID2 == 1 && s.TypeID3 == 4 {
		switch s.TypeID4 {
		case 1:
			return ChineseShield
		case 2:
			return EuropeanShield
		default:
			return -1
		}
	}
	return -1
}
