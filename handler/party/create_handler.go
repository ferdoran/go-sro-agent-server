package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

type PartyAgentCreateRequestHandler struct {
}

func (h *PartyAgentCreateRequestHandler) Handle(data server.PacketChannelData) {

}

func SendPartyCreateResponse(ptMasterUniqueId uint32) {
	p := network.EmptyPacket()
	p.MessageID = opcode.PartyCreateResponse
	p.WriteByte(1)
	p.WriteUInt32(1)
	ptMaster := model.GetSroWorldInstance().PlayersByUniqueId[ptMasterUniqueId]
	ptMaster.Session.Conn.Write(p.ToBytes())
}
