package chat

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

func handleGlobalMessage(request MessageRequest, session *server.Session) {
	// TODO: Remove global from players inventory

	players := model.GetSroWorldInstance().PlayersByUniqueId
	for _, v := range players {
		p := network.EmptyPacket()
		p.MessageID = opcode.ChatUpdate
		p.WriteByte(request.ChatType)
		p.WriteString(session.UserContext.CharName)
		p.WriteString(request.Message)
		v.Session.Conn.Write(p.ToBytes())
	}
}
