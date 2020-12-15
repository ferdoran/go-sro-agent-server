package stall

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"sync"
)

func addRemoveItem(data server.PacketChannelData, updateType byte) {
	stallSlot, err := data.ReadByte()
	if err != nil {
		log.Panicln("Failed to read stall slot")
	}

	var inventorySlot byte
	var stackCount uint16
	var fleaMarketNetworkTidGroup uint32
	var price uint64

	if updateType == StallAddItem {
		inventorySlot, stackCount, price, fleaMarketNetworkTidGroup = getExtraData(data)
	}

	unkUShort0, err5 := data.ReadUInt16()
	if err5 != nil {
		log.Panicln("Failed to read unkUShort0")
	}

	stallEntry := model.StallEntry{
		StallSlot:                 stallSlot,
		InventorySlot:             inventorySlot,
		StackCount:                stackCount,
		Price:                     price,
		FleaMarketnetworkTidGroup: fleaMarketNetworkTidGroup,
		UnkUshort0:                unkUShort0,
		Mutex:                     &sync.Mutex{},
	}

	var playerStall *model.Stall
	if updateType == StallAddItem {
		playerStall = stallEntry.AddItem(data.UserContext.UniqueID)
	} else if updateType == StallRemoveItem {
		playerStall = stallEntry.RemoveItem(data.UserContext.UniqueID)
	}

	playerInventoryItems := model.GetSroWorldInstance().PlayersByUniqueId[data.UserContext.UniqueID].Inventory.Items

	p := network.EmptyPacket()
	p.MessageID = opcode.StallUpdateResponse
	p.WriteByte(1)
	p.WriteByte(updateType)
	p.WriteUInt16(0) // TODO: Handle errors correctly
	if playerStall != nil && playerStall.Entries != nil {
		for _, v := range playerStall.Entries {
			item := playerInventoryItems[v.InventorySlot]
			p.WriteByte(v.StallSlot)
			model.WriteRentInfo(&p, item)
			p.WriteUInt32(item.GetRefObjectID())
			model.WriteInventoryItem(&p, item)
			p.WriteByte(v.InventorySlot)
			p.WriteUInt16(v.StackCount)
			p.WriteUInt64(v.Price)
		}
	}
	p.WriteByte(255) // End of data
	data.Session.Conn.Write(p.ToBytes())
}

func getExtraData(data server.PacketChannelData) (byte, uint16, uint64, uint32) {
	inventorySlot, err1 := data.ReadByte()
	if err1 != nil {
		log.Panicln("Failed to read inventorySlot")
	}

	stackCount, err2 := data.ReadUInt16()
	if err2 != nil {
		log.Panicln("Failed to read stack count of item")
	}

	price, err3 := data.ReadUInt64()
	if err3 != nil {
		log.Panicln("Failed to read price of item")
	}

	fleaMarketNetworkTidGroup, err4 := data.ReadUInt32()
	if err4 != nil {
		log.Panicln("Failed to read flea market network id")
	}

	return inventorySlot, stackCount, price, fleaMarketNetworkTidGroup
}
