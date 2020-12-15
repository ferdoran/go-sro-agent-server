package model

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/utils"
	"time"
)

func WriteRentInfo(p *network.Packet, item Item) {
	// TODO find out what the rent type is
	rentType := uint32(0)
	p.WriteUInt32(rentType) // TODO RentType

	switch rentType {
	case 1:
		p.WriteUInt16(1)                                                         // RentInfo.CanDelete
		p.WriteUInt32(utils.ToSilkroadTime(time.Now()))                          // RentInfo.PeriodBeginTime
		p.WriteUInt32(utils.ToSilkroadTime(time.Now().Add(time.Hour * 24 * 30))) // RentInfo.PeriodEndTime
	case 2:
		p.WriteUInt16(1)                                               // RentInfo.CanDelete
		p.WriteUInt16(1)                                               // RentInfo.CanRecharge
		p.WriteUInt32(utils.ToSilkroadTime(time.Now().Add(time.Hour))) // RentInfo.MeterRateTime
	case 3:
		p.WriteUInt16(1)                                                         // RentInfo.CanDelete
		p.WriteUInt16(1)                                                         // RentInfo.CanRecharge
		p.WriteUInt32(utils.ToSilkroadTime(time.Now()))                          // RentInfo.PeriodBeginTime
		p.WriteUInt32(utils.ToSilkroadTime(time.Now().Add(time.Hour * 24 * 30))) // RentInfo.PeriodEndTime
		p.WriteUInt32(utils.ToSilkroadTime(time.Now().Add(time.Second * 2)))     // RentInfo.PackingTime
	}
}

func WriteInventoryItem(p *network.Packet, item Item) {
	if item.IsEquipment() {
		WriteEquipmentItem(p, item)
	} else if item.IsContainer() {
		WriteContainerItem(p, item)
	} else if item.IsExpendable() {
		WriteExpendableItem(p, item)
	}
}

func WriteEquipmentItem(p *network.Packet, item Item) {
	numMagParams := 0
	p.WriteByte(0) // OptLevel / Plus
	p.WriteUInt64(item.Variance)
	p.WriteUInt32(10) // Item Data - Probably durability
	p.WriteByte(0)    // Number of MagParams / Blues

	// TODO write mag params
	for i := 0; i < numMagParams; i++ {
		p.WriteUInt32(0) // magParam.Type
		p.WriteUInt32(0) // magParam.Value
	}

	// TODO write bind options
	numBindOptions := 0
	p.WriteByte(1) // bindingOptionType - 1 = Socket
	p.WriteByte(0) // bindingOptionCount

	for i := 0; i < numBindOptions; i++ {
		p.WriteByte(0)   // BindingOption.Slot
		p.WriteUInt32(0) // BindingOption.ID
		p.WriteUInt32(0) // BindingOption.nParam1
	}

	p.WriteByte(2) // bindingOptionType - 2 = Advanced Elixir
	p.WriteByte(0) // bindingOptionCount
	for i := 0; i < numBindOptions; i++ {
		p.WriteByte(0)   // BindingOption.Slot
		p.WriteUInt32(0) // BindingOption.ID
		p.WriteUInt32(0) // BindingOption.OptValue / Plus
	}

}

func WriteContainerItem(p *network.Packet, item Item) {
	if item.TypeID3 == 1 {
		// ITEM_COS_P
		p.WriteByte(0) // TODO COS State
		p.WriteUInt32(item.GetRefObjectID())
		p.WriteString(item.Name)

		if item.TypeID4 == 2 {
			p.WriteUInt32(3600) // TODO Seconds to Rent End Time
		}
		p.WriteByte(0) // TODO Unknown

	} else if item.TypeID3 == 2 {
		// ITEM_ETC_TRANS_MONSTER
		p.WriteUInt32(item.GetRefObjectID())
	} else if item.TypeID3 == 3 {
		// MAGIC_CUBE
		p.WriteUInt32(0) // TODO Amount of elixirs inside the cube
	}
}

func WriteExpendableItem(p *network.Packet, item Item) {
	p.WriteUInt16(1) // TODO Stack Size

	if item.TypeID3 == 11 {
		if item.TypeID4 == 1 || item.TypeID4 == 2 {
			// MAGICSTONE, ATTRSTONE
			p.WriteByte(0) // TODO AttributeAssimilationProbability
		}
	} else if item.TypeID3 == 14 || item.TypeID4 == 2 {
		// ITEM_MALL_GACHA_CARD_WIN
		// ITEM_MALL_GACHA_CARD_LOSE
		numMagParams := 0
		p.WriteByte(0) // TODO item.MagParamCount

		for i := 0; i < numMagParams; i++ {
			p.WriteUInt32(0) // TODO magParam.Type
			p.WriteUInt32(0) // TODO magParam.Value
		}
	}
}
