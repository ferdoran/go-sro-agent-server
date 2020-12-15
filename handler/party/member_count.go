package party

import (
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
)

func SendMemberCountResponse(requestingPlayerUniqueId uint32, memberCount byte) {
	p := network.EmptyPacket()
	p.MessageID = opcode.PartyMemberCountResponse
	p.WriteByte(1)
	p.WriteUInt32(uint32(memberCount))
	requestingPlayer := model.GetSroWorldInstance().PlayersByUniqueId[requestingPlayerUniqueId]
	requestingPlayer.Session.Conn.Write(p.ToBytes())
}