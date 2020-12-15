package model

import (
	"github.com/sirupsen/logrus"
	"sync"
)

const BaseInventorySize byte = 45

const (
	SlotHelmet          byte = 0
	SlotChest           byte = 1
	SlotShoulder        byte = 2
	SlotGlove           byte = 3
	SlotPants           byte = 4
	SlotBoots           byte = 5
	SlotPrimaryWeapon   byte = 6
	SlotSecondaryWeapon byte = 7
	SlotExtra           byte = 8
	SlotEarring         byte = 9
	SlotNecklace        byte = 10
	SlotLeftRing        byte = 11
	SlotRightRing       byte = 12
)

const (
	EquipItem = iota
	UnequipItem
	MoveItem
)

type Inventory struct {
	Items map[byte]Item
	mutex sync.Mutex
}

func (i *Inventory) MoveItems(sourceSlot, targetSlot byte, player *Player) (bool, int) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if targetSlot <= SlotRightRing {
		// TODO implement other Slot Checks
		var canEquip bool
		sourceItem := i.Items[sourceSlot]
		switch targetSlot {
		case SlotPrimaryWeapon:
			canEquip = i.canEquipPrimaryWeapon(sourceItem, player)
		case SlotSecondaryWeapon:
			canEquip = i.canEquipSecondaryWeapon(sourceItem, player)
		case SlotHelmet:
			canEquip = i.canEquipHelmet(sourceItem, player)
		case SlotChest:
			canEquip = i.canEquipChest(sourceItem, player)
		case SlotShoulder:
			canEquip = i.canEquipShoulder(sourceItem, player)
		case SlotGlove:
			canEquip = i.canEquipGloves(sourceItem, player)
		case SlotPants:
			canEquip = i.canEquipPants(sourceItem, player)
		case SlotBoots:
			canEquip = i.canEquipBoots(sourceItem, player)
		case SlotEarring:
			canEquip = i.canEquipEarring(sourceItem, player)
		case SlotNecklace:
			canEquip = i.canEquipNecklace(sourceItem, player)
		case SlotLeftRing:
			canEquip = i.canEquipRing(sourceItem, player)
		case SlotRightRing:
			canEquip = i.canEquipRing(sourceItem, player)
		default:
			logrus.Infof("player [%s] is trying to equip items between slots %d and %d\n", player.CharName, sourceSlot, targetSlot)
		}

		if canEquip {
			i.swapItems(sourceSlot, targetSlot)
			return true, EquipItem
		} else {
			logrus.Debugf("player [%s] cannot equip item [%v]. Target slot contains [%v]\n", player.CharName, sourceItem, i.Items[targetSlot])
		}
	} else if sourceSlot <= SlotRightRing && i.Items[sourceSlot].IsEquipment() {
		if _, ok := player.Inventory.Items[targetSlot]; !ok {
			i.swapItems(sourceSlot, targetSlot)
			return true, UnequipItem
		}
	} else {
		// TODO any more checks required?
		logrus.Infof("player [%s] is moving items between slots %d and %d\n", player.CharName, sourceSlot, targetSlot)
		i.swapItems(sourceSlot, targetSlot)
		return true, MoveItem
	}

	return false, -1
}

func (i *Inventory) canEquipPrimaryWeapon(sourceItem Item, player *Player) bool {
	if (player.IsChinese() && sourceItem.IsChineseWeapon()) || (player.IsEuropean() && sourceItem.IsEuropeanWeapon()) {
		if sourceItem.RequiredLevelType1 == 1 && player.Level >= sourceItem.RequiredLevel1 {
			return true
		}
	}

	return false
}

func (i *Inventory) swapItems(sourceSlot, targetSlot byte) {
	tmp, ok := i.Items[targetSlot]
	i.Items[targetSlot] = i.Items[sourceSlot]

	if !ok {
		delete(i.Items, sourceSlot)
	} else {
		i.Items[sourceSlot] = tmp
	}
}

func (i *Inventory) canEquipArmorPart(item Item, targetSlot byte) bool {
	numItemsOfOtherType := 0
	itemsHaveSameArmorType := false
	if item.IsCHGarmentPart() || item.IsEUGarmentPart() {
		numItemsOfOtherType = i.equippedArmorOrProtectorParts()
		if equippedItem, ok := i.Items[targetSlot]; ok && (equippedItem.IsCHGarmentPart() || equippedItem.IsEUGarmentPart()) {
			itemsHaveSameArmorType = true
		}
	} else if item.IsCHArmorPart() || item.IsCHProtectorPart() || item.IsEUArmorPart() || item.IsEUProtectorPart() {
		numItemsOfOtherType = i.equippedGarmentParts()
		if equippedItem, ok := i.Items[targetSlot]; ok && (equippedItem.IsCHArmorPart() || equippedItem.IsCHProtectorPart() || equippedItem.IsEUArmorPart() || equippedItem.IsEUProtectorPart()) {
			itemsHaveSameArmorType = true
		}
	}
	switch numItemsOfOtherType {
	case 0:
		return true
	case 1:
		return itemsHaveSameArmorType
	default:
		return false
	}
}

func (i *Inventory) canEquipHelmet(item Item, player *Player) bool {
	if player.IsChinese() && item.IsCHHelmet() {
		return i.canEquipArmorPart(item, SlotHelmet) && i.fulfillsLevelRequirements(item, player)
	} else if player.IsEuropean() && item.IsEUHelmet() {
		if i.canEquipArmorPart(item, SlotHelmet) {
			return i.fulfillsSkillRequirements(item, player) && i.fulfillsLevelRequirements(item, player)
		}
	}
	return false
}

func (i *Inventory) canEquipShoulder(item Item, player *Player) bool {
	if player.IsChinese() && item.IsCHShoulder() {
		return i.canEquipArmorPart(item, SlotShoulder) && i.fulfillsLevelRequirements(item, player)
	} else if player.IsEuropean() && item.IsEUShoulder() {
		if i.canEquipArmorPart(item, SlotShoulder) {
			return i.fulfillsSkillRequirements(item, player) && i.fulfillsLevelRequirements(item, player)
		}
	}
	return false
}

func (i *Inventory) canEquipChest(item Item, player *Player) bool {
	if player.IsChinese() && item.IsCHChest() {
		return i.canEquipArmorPart(item, SlotChest) && i.fulfillsLevelRequirements(item, player)
	} else if player.IsEuropean() && item.IsEUChest() {
		if i.canEquipArmorPart(item, SlotChest) {
			return i.fulfillsSkillRequirements(item, player) && i.fulfillsLevelRequirements(item, player)
		}
	}
	return false
}

func (i *Inventory) canEquipPants(item Item, player *Player) bool {
	if player.IsChinese() && item.IsCHPant() {
		return i.canEquipArmorPart(item, SlotPants) && i.fulfillsLevelRequirements(item, player)
	} else if player.IsEuropean() && item.IsEUPant() {
		if i.canEquipArmorPart(item, SlotPants) {
			return i.fulfillsSkillRequirements(item, player) && i.fulfillsLevelRequirements(item, player)
		}
	}
	return false
}

func (i *Inventory) canEquipGloves(item Item, player *Player) bool {
	if player.IsChinese() && item.IsCHGlove() {
		return i.canEquipArmorPart(item, SlotGlove) && i.fulfillsLevelRequirements(item, player)
	} else if player.IsEuropean() && item.IsEUGlove() {
		if i.canEquipArmorPart(item, SlotGlove) {
			return i.fulfillsSkillRequirements(item, player) && i.fulfillsLevelRequirements(item, player)
		}
	}
	return false
}

func (i *Inventory) canEquipBoots(item Item, player *Player) bool {
	if player.IsChinese() && item.IsCHBoots() {
		return i.canEquipArmorPart(item, SlotBoots) && i.fulfillsLevelRequirements(item, player)
	} else if player.IsEuropean() && item.IsEUBoots() {
		if i.canEquipArmorPart(item, SlotBoots) {
			return i.fulfillsSkillRequirements(item, player) && i.fulfillsLevelRequirements(item, player)
		}
	}
	return false
}

func (i *Inventory) canEquipNecklace(item Item, player *Player) bool {
	if (player.IsChinese() && item.IsCHNecklace()) || (player.IsEuropean() && item.IsEUNecklace()) {
		return i.fulfillsLevelRequirements(item, player)
	}
	return false
}

func (i *Inventory) canEquipEarring(item Item, player *Player) bool {
	if (player.IsChinese() && item.IsCHEarring()) || (player.IsEuropean() && item.IsEUEarring()) {
		return i.fulfillsLevelRequirements(item, player)
	}
	return false
}

func (i *Inventory) canEquipRing(item Item, player *Player) bool {
	if (player.IsChinese() && item.IsCHRing()) || (player.IsEuropean() && item.IsEURing()) {
		return i.fulfillsLevelRequirements(item, player)
	}
	return false
}

func (i *Inventory) equippedGarmentParts() int {
	numGarmentItems := 0
	for j := byte(0); j < SlotBoots; j++ {
		if item, ok := i.Items[j]; ok && (item.IsCHGarmentPart() || item.IsEUGarmentPart()) {
			numGarmentItems++
		}
	}
	return numGarmentItems
}

func (i *Inventory) equippedArmorOrProtectorParts() int {
	numArmorOrProtectorParts := 0
	for j := byte(0); j < SlotBoots; j++ {
		if item, ok := i.Items[j]; ok && (item.IsCHArmorPart() || item.IsEUArmorPart() || item.IsCHProtectorPart() || item.IsEUProtectorPart()) {
			numArmorOrProtectorParts++
		}
	}
	return numArmorOrProtectorParts
}

func (i *Inventory) fulfillsLevelRequirements(item Item, player *Player) bool {
	// TODO: is this all that has to be checked?
	return item.IsEquipment() && item.RequiredLevelType1 == 1 && item.RequiredLevel1 <= player.Level
}

func (i *Inventory) fulfillsSkillRequirements(item Item, player *Player) bool {
	// TODO: Check level and skill requirements
	return false
}

func (i *Inventory) canEquipSecondaryWeapon(item Item, player *Player) bool {
	// TODO any further checks needed? (e.g. siege weapons)
	primWeapon, ok := i.Items[SlotPrimaryWeapon]
	if player.IsChinese() {
		if ok && primWeapon.IsOneHandedWeapon() && item.IsCHShield() && i.fulfillsLevelRequirements(item, player) {
			return true
		} else if !ok && item.IsCHShield() && i.fulfillsLevelRequirements(item, player) {
			return true
		} else if !ok && primWeapon.IsBow() && item.IsArrow() {
			return true
		}

	} else if player.IsEuropean() {
		if ok && primWeapon.IsOneHandedWeapon() && item.IsEUShield() && i.fulfillsLevelRequirements(item, player) {
			return true
		} else if !ok && item.IsEUShield() && i.fulfillsLevelRequirements(item, player) {
			return true
		} else if !ok && primWeapon.IsCrossbow() && item.IsBolt() {
			return true
		}
	}
	return false
}
