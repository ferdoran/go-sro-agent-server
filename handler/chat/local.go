package chat

import (
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

func handleAllMessage(request MessageRequest, session *server.Session) {
	// TODO: Change all players to local region
	players := model.GetSroWorldInstance().PlayersByUniqueId
	for _, v := range players {
		p := network.EmptyPacket()
		p.MessageID = opcode.ChatUpdate
		p.WriteByte(request.ChatType)
		p.WriteUInt32(session.UserContext.UniqueID)
		p.WriteString(request.Message)
		v.Session.Conn.Write(p.ToBytes())
	}

	p1 := network.EmptyPacket()
	p1.MessageID = opcode.ChatResponse
	p1.WriteByte(1)
	p1.WriteByte(request.ChatType)
	p1.WriteByte(request.ChatIndex)
	session.Conn.Write(p1.ToBytes())
}