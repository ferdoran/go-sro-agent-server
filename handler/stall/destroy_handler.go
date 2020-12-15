package stall

import (
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

type StallDestroyHandler struct {
}

func NewStallDestroyHandler() server.PacketHandler {
	handler := StallDestroyHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.StallDestroyRequest, handler)
	return handler
}

func (s StallDestroyHandler) Handle(data server.PacketChannelData) {
	p := network.EmptyPacket()
	p.MessageID = opcode.StallDestroyResponse
	p.WriteByte(1)
	data.Session.Conn.Write(p.ToBytes())

	p2 := network.EmptyPacket()
	p2.MessageID = opcode.StallEntityDestroyResponse
	p2.WriteUInt32(data.Session.UserContext.UniqueID)
	p2.WriteUInt16(0) // Error Code. TODO: Possible Error Code 15383 - Closing
	data.Session.Conn.Write(p2.ToBytes())
}
