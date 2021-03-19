package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

type StallDestroyHandler struct {
	channel chan server.PacketChannelData
}

func InitStallDestroyHandler() {
	handler := StallDestroyHandler{channel: server.PacketManagerInstance.GetQueue(opcode.StallDestroyRequest)}
	go handler.Handle()
}

func (s *StallDestroyHandler) Handle() {
	for {
		data := <-s.channel
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
}
