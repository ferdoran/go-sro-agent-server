package model

type TypeInfo struct {
	TypeID1 int
	TypeID2 int
	TypeID3 int
	TypeID4 int
}

func (ti TypeInfo) IsEquipment() bool {
	return ti.TypeID1 == 3 && ti.TypeID2 == 1
}

func (ti TypeInfo) IsItem() bool {
	return ti.TypeID1 == 3
}

func (ti TypeInfo) IsContainer() bool {
	return ti.TypeID1 == 3 && ti.TypeID2 == 2
}

func (ti TypeInfo) IsExpendable() bool {
	return ti.TypeID1 == 3 && ti.TypeID2 == 3
}

func (ti TypeInfo) IsWeapon() bool {
	return ti.IsEquipment() && ti.TypeID3 == 6
}

func (ti TypeInfo) IsChineseWeapon() bool {
	return ti.IsWeapon() && ti.TypeID4 > 1 && ti.TypeID4 < 7
}

func (ti TypeInfo) IsEuropeanWeapon() bool {
	return ti.IsWeapon() && ti.TypeID4 >= 7 && ti.TypeID4 < 16
}

func (ti TypeInfo) IsShield() bool {
	return ti.IsEquipment() && ti.TypeID3 == 4
}

func (ti TypeInfo) IsCHShield() bool {
	return ti.IsShield() && ti.TypeID4 == 1
}

func (ti TypeInfo) IsEUShield() bool {
	return ti.IsShield() && ti.TypeID4 == 2
}

func (ti TypeInfo) IsSword() bool {
	return ti.IsWeapon() && ti.TypeID4 == 2
}

func (ti TypeInfo) IsBow() bool {
	return ti.IsWeapon() && ti.TypeID4 == 6
}

func (ti TypeInfo) IsCrossbow() bool {
	return ti.IsWeapon() && ti.TypeID4 == 12
}

func (ti TypeInfo) IsBlade() bool {
	return ti.IsWeapon() && ti.TypeID4 == 2
}

func (ti TypeInfo) Is1HSword() bool {
	return ti.IsWeapon() && ti.TypeID4 == 2
}

func (ti TypeInfo) IsWarlockRod() bool {
	return ti.IsWeapon() && ti.TypeID4 == 10
}

func (ti TypeInfo) IsClericRod() bool {
	return ti.IsWeapon() && ti.TypeID4 == 15
}

func (ti TypeInfo) IsOneHandedWeapon() bool {
	return ti.IsSword() || ti.IsBlade() || ti.Is1HSword() || ti.IsWarlockRod() || ti.IsClericRod()
}

func (ti TypeInfo) IsCHArmorPart() bool {
	return ti.IsEquipment() && ti.TypeID3 == 3
}

func (ti TypeInfo) IsCHProtectorPart() bool {
	return ti.IsEquipment() && ti.TypeID3 == 2
}

func (ti TypeInfo) IsCHGarmentPart() bool {
	return ti.IsEquipment() && ti.TypeID3 == 1
}

func (ti TypeInfo) IsCHAccessory() bool {
	return ti.IsEquipment() && ti.TypeID3 == 5
}

func (ti TypeInfo) IsEUAccessory() bool {
	return ti.IsEquipment() && ti.TypeID3 == 12
}

func (ti TypeInfo) IsEUArmorPart() bool {
	return ti.IsEquipment() && ti.TypeID3 == 11
}

func (ti TypeInfo) IsEUProtectorPart() bool {
	return ti.IsEquipment() && ti.TypeID3 == 10
}

func (ti TypeInfo) IsEUGarmentPart() bool {
	return ti.IsEquipment() && ti.TypeID3 == 9
}

func (ti TypeInfo) IsCHHelmet() bool {
	return (ti.IsCHArmorPart() || ti.IsCHProtectorPart() || ti.IsCHGarmentPart()) && ti.TypeID4 == 1
}

func (ti TypeInfo) IsCHShoulder() bool {
	return (ti.IsCHArmorPart() || ti.IsCHProtectorPart() || ti.IsCHGarmentPart()) && ti.TypeID4 == 2
}

func (ti TypeInfo) IsCHChest() bool {
	return (ti.IsCHArmorPart() || ti.IsCHProtectorPart() || ti.IsCHGarmentPart()) && ti.TypeID4 == 3
}

func (ti TypeInfo) IsCHPant() bool {
	return (ti.IsCHArmorPart() || ti.IsCHProtectorPart() || ti.IsCHGarmentPart()) && ti.TypeID4 == 4
}

func (ti TypeInfo) IsCHGlove() bool {
	return (ti.IsCHArmorPart() || ti.IsCHProtectorPart() || ti.IsCHGarmentPart()) && ti.TypeID4 == 5
}

func (ti TypeInfo) IsCHBoots() bool {
	return (ti.IsCHArmorPart() || ti.IsCHProtectorPart() || ti.IsCHGarmentPart()) && ti.TypeID4 == 6
}

func (ti TypeInfo) IsEUHelmet() bool {
	return (ti.IsEUArmorPart() || ti.IsEUProtectorPart() || ti.IsEUGarmentPart()) && ti.TypeID4 == 1
}

func (ti TypeInfo) IsEUShoulder() bool {
	return (ti.IsEUArmorPart() || ti.IsEUProtectorPart() || ti.IsEUGarmentPart()) && ti.TypeID4 == 2
}

func (ti TypeInfo) IsEUChest() bool {
	return (ti.IsEUArmorPart() || ti.IsEUProtectorPart() || ti.IsEUGarmentPart()) && ti.TypeID4 == 3
}

func (ti TypeInfo) IsEUPant() bool {
	return (ti.IsEUArmorPart() || ti.IsEUProtectorPart() || ti.IsEUGarmentPart()) && ti.TypeID4 == 4
}

func (ti TypeInfo) IsEUGlove() bool {
	return (ti.IsEUArmorPart() || ti.IsEUProtectorPart() || ti.IsEUGarmentPart()) && ti.TypeID4 == 5
}

func (ti TypeInfo) IsEUBoots() bool {
	return (ti.IsEUArmorPart() || ti.IsEUProtectorPart() || ti.IsEUGarmentPart()) && ti.TypeID4 == 6
}

func (ti TypeInfo) IsCHRing() bool {
	return ti.IsCHAccessory() && ti.TypeID4 == 3
}

func (ti TypeInfo) IsEURing() bool {
	return ti.IsEUAccessory() && ti.TypeID4 == 3
}

func (ti TypeInfo) IsCHEarring() bool {
	return ti.IsCHAccessory() && ti.TypeID4 == 1
}

func (ti TypeInfo) IsEUEarring() bool {
	return ti.IsEUAccessory() && ti.TypeID4 == 1
}

func (ti TypeInfo) IsCHNecklace() bool {
	return ti.IsCHAccessory() && ti.TypeID4 == 2
}

func (ti TypeInfo) IsEUNecklace() bool {
	return ti.IsEUAccessory() && ti.TypeID4 == 2
}

func (ti TypeInfo) IsArrow() bool {
	return ti.IsExpendable() && ti.TypeID3 == 4 && ti.TypeID4 == 1
}

func (ti TypeInfo) IsBolt() bool {
	return ti.IsExpendable() && ti.TypeID3 == 4 && ti.TypeID4 == 2
}

func (ti TypeInfo) IsCharacter() bool {
	return ti.TypeID1 == 1
}

func (ti TypeInfo) IsPlayerCharacter() bool {
	return ti.IsCharacter() && ti.TypeID2 == 1
}

func (ti TypeInfo) IsNPC() bool {
	return ti.IsCharacter() && ti.TypeID2 == 2
}

func (ti TypeInfo) IsNPCMob() bool {
	return ti.IsNPC() && ti.TypeID3 == 1
}

func (ti TypeInfo) IsNPCNpc() bool {
	return ti.IsNPC() && ti.TypeID3 == 2
}

func (ti TypeInfo) IsCOS() bool {
	return ti.IsNPC() && ti.TypeID3 == 3
}

func (ti TypeInfo) IsSiegeObject() bool {
	return ti.IsNPC() && ti.TypeID3 == 4
}

func (ti TypeInfo) IsSiegeStruct() bool {
	return ti.IsNPC() && ti.TypeID3 == 5
}

func (ti TypeInfo) IsGold() bool {
	return ti.IsExpendable() && ti.TypeID3 == 5 && ti.TypeID4 == 0
}

func (ti TypeInfo) IsTradeItem() bool {
	return ti.IsExpendable() && ti.TypeID3 == 8
}

func (ti TypeInfo) IsQuestItem() bool {
	return ti.IsExpendable() && ti.TypeID3 == 9
}

func (ti TypeInfo) IsStructure() bool {
	return ti.TypeID1 == 4
}
