package chat

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

func handlePartyMessage(request MessageRequest, session *server.Session) {
	player := model.GetSroWorldInstance().PlayersByUniqueId[session.UserContext.UniqueID]

	p := network.EmptyPacket()
	p.MessageID = opcode.ChatResponse
	p.WriteByte(1)
	p.WriteByte(request.ChatType)
	p.WriteByte(request.ChatIndex)
	session.Conn.Write(p.ToBytes())

	for _, v := range player.Party.Members {
		p1 := network.EmptyPacket()
		p1.MessageID = opcode.ChatUpdate
		p1.WriteByte(request.ChatType)
		p1.WriteString(session.UserContext.CharName)
		p1.WriteString(request.Message)
		recievingPlayer := model.GetSroWorldInstance().PlayersByUniqueId[v]
		recievingPlayer.Session.Conn.Write(p1.ToBytes())
	}
}
