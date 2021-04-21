package inventory

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/sirupsen/logrus"
)

type InventoryHandler struct {
	channel chan server.PacketChannelData
}

type ItemInventoryOperation byte

const (
	MoveItemOperation ItemInventoryOperation = iota
	_
	PutItemIntoStorageOperation
	TakeItemFromStorageOperation
	PutItemIntoExchangeWindowOperation
	TakeItemFromExchangeWindowOperation // 5
	_
	DropItemOperation
	BuyItemOperation
	SellItemOperation
	DropGoldOperation // 10
	_
	_
	_
	_
	_ // 15
	_
	_
	BuyItemFromMallOperation
	BuySpecialtyGoodsOperation
	SellSpecialtyGoodsOperation // 20
	_
	_
	EquipAvatarItemOperation
	UnequipAvatarItemOperation
	_ // 25
	PutItemIntoPetInventoryOperation
	TakeItemFromPetInventoryOperation
)

func InitInventoryHandler() {
	queue := server.PacketManagerInstance.GetQueue(opcode.ItemOperationRequest)
	handler := InventoryHandler{channel: queue}
	go handler.Handle()
}

func (h *InventoryHandler) Handle() {
	for {
		data := <-h.channel
		operationType, err := data.ReadByte()
		if err != nil {
			// FIXME not necessarily a fail but after a successful exchange a 0x7034 packet without payload is sent for each item that has been traded
			logrus.Errorf("failed to read inventory operation type\n")
		}

		switch ItemInventoryOperation(operationType) {
		case MoveItemOperation:
			moveItem(data)
		case DropGoldOperation:
			dropGold(data)
		}
	}
}

func dropGold(data server.PacketChannelData) {
	goldAmount, err := data.ReadUInt64()

	if err != nil {
		logrus.Errorf("failed to read gold amount")
	}

	logrus.Infof("Player [%s] is dropping %d gold\n", data.UserContext.CharName, goldAmount)
}

func moveItem(data server.PacketChannelData) {
	sourceSlot, err := data.ReadByte()
	if err != nil {
		logrus.Errorf("failed to read source slot\n")
	}

	targetSlot, err := data.ReadByte()
	if err != nil {
		logrus.Errorf("failed to read target slot\n")
	}

	amount, err := data.ReadUInt16()
	if err != nil {
		logrus.Errorf("failed to read unknown value\n")
	}

	if sourceSlot == targetSlot {
		return
	}

	world := service.GetWorldServiceInstance()
	player, err := world.GetPlayerByCharName(data.Session.UserContext.CharName)
	if err == nil {
		// player not online
		logrus.Tracef("Player %s is not online\n", player.CharName)
		return
	}

	movedItems, moveAction := player.Inventory.MoveItems(sourceSlot, targetSlot, player)

	if !movedItems {
		logrus.Tracef("Failed to move items %v %v\n", movedItems, moveAction)
		return
	}

	switch moveAction {
	case model.EquipItem:
		// Send equip item packet
		player.SendEquipItemPacket(player.Inventory.Items[targetSlot], targetSlot)
		player.SendStatsUpdate()
	case model.UnequipItem:
		// Send unequip item packet
		player.SendUnequipItemPacket(player.Inventory.Items[targetSlot], sourceSlot)
		player.SendStatsUpdate()
	}

	p := network.EmptyPacket()
	p.MessageID = opcode.ItemOperationResponse
	p.WriteBool(true) // result
	p.WriteByte(byte(MoveItemOperation))
	p.WriteByte(sourceSlot)
	p.WriteByte(targetSlot)
	p.WriteUInt16(amount)
	p.WriteByte(0) // TODO find out what this value stands for

	data.Conn.Write(p.ToBytes())

}
