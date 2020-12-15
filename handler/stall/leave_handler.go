package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

type StallLeaveHandler struct {
}

func NewStallLeaveHandler() server.PacketHandler {
	handler := StallLeaveHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.StallLeaveRequest, handler)
	return handler
}

func (s StallLeaveHandler) Handle(data server.PacketChannelData) {
	p := network.EmptyPacket()
	p.MessageID = opcode.StallLeaveResponse
	p.WriteByte(1)
	data.Session.Conn.Write(p.ToBytes())
}
