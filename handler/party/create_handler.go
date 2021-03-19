package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
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

func SendPartyCreateResponse(ptMasterUniqueId uint32) {
	p := network.EmptyPacket()
	p.MessageID = opcode.PartyCreateResponse
	p.WriteByte(1)
	p.WriteUInt32(1)
	ptMaster := model.GetSroWorldInstance().PlayersByUniqueId[ptMasterUniqueId]
	ptMaster.Session.Conn.Write(p.ToBytes())
}
