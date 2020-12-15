package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
)

func SendMemberCountResponse(requestingPlayerUniqueId uint32, memberCount byte) {
	p := network.EmptyPacket()
	p.MessageID = opcode.PartyMemberCountResponse
	p.WriteByte(1)
	p.WriteUInt32(uint32(memberCount))
	requestingPlayer := model.GetSroWorldInstance().PlayersByUniqueId[requestingPlayerUniqueId]
	requestingPlayer.Session.Conn.Write(p.ToBytes())
}
