package chat

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
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
