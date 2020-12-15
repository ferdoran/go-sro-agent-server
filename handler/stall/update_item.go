package stall

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"sync"
)

func updateItem(data server.PacketChannelData) {
	stallSlot, err := data.ReadByte()
	if err != nil {
		log.Panicln("Failed to read stall slot")
	}

	stackCount, err2 := data.ReadUInt16()
	if err2 != nil {
		log.Panicln("Failed to read stack count")
	}

	price, err3 := data.ReadUInt64()
	if err3 != nil {
		log.Panicln("Failed to read price")
	}

	unkUShort0, err4 := data.ReadUInt16()
	if err4 != nil {
		log.Panicln("Failed to read unkUShort0")
	}

	stallEntry := model.StallEntry{
		StallSlot:  stallSlot,
		StackCount: stackCount,
		Price:      price,
		UnkUshort0: unkUShort0,
		Mutex:      &sync.Mutex{},
	}

	stallEntry = stallEntry.UpdateItem(data.UserContext.UniqueID)

	p := network.EmptyPacket()
	p.MessageID = opcode.StallUpdateResponse
	p.WriteByte(1)
	p.WriteByte(StallUpdateItem)
	p.WriteByte(stallEntry.StallSlot)
	p.WriteUInt16(stallEntry.StackCount)
	p.WriteUInt64(stallEntry.Price)
	p.WriteUInt16(0) // TODO: Handle errors correctly
	data.Session.Conn.Write(p.ToBytes())
}
