package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

type StallLeaveHandler struct {
	channel chan server.PacketChannelData
}

func InitStallLeaveHandler() {
	handler := StallLeaveHandler{channel: server.PacketManagerInstance.GetQueue(opcode.StallLeaveRequest)}
	go handler.Handle()
}

func (s *StallLeaveHandler) Handle() {
	for {
		data := <-s.channel
		p := network.EmptyPacket()
		p.MessageID = opcode.StallLeaveResponse
		p.WriteByte(1)
		data.Session.Conn.Write(p.ToBytes())
	}
}
