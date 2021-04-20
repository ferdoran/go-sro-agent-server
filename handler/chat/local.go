package chat

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

func handleAllMessage(request MessageRequest, session *server.Session) {
	p := network.EmptyPacket()
	p.MessageID = opcode.ChatUpdate
	p.WriteByte(request.ChatType)
	p.WriteUInt32(session.UserContext.UniqueID)
	p.WriteString(request.Message)
	// TODO: Change all players to local region
	service.GetWorldServiceInstance().BroadcastRaw(p.ToBytes())

	p1 := network.EmptyPacket()
	p1.MessageID = opcode.ChatResponse
	p1.WriteByte(1)
	p1.WriteByte(request.ChatType)
	p1.WriteByte(request.ChatIndex)
	session.Conn.Write(p1.ToBytes())
}
