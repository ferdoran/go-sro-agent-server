package party

import (
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

type PartyAgentCreateRequestHandler struct {
	channel chan server.PacketChannelData
}

func InitPartyAgentCreateRequestHandler() {
	handler := PartyAgentCreateRequestHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PartyCreateRequest)}
	go handler.Handle()
}

func (h *PartyAgentCreateRequestHandler) Handle() {
	// TODO: implement
	//for {
	//	data := <- h.channel
	//}
}
