package chat

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
)

func handleStallMessage(request MessageRequest, session *server.Session) {
	p := network.EmptyPacket()
	p.MessageID = opcode.ChatResponse
	p.WriteByte(1)
	p.WriteByte(request.ChatType)
	p.WriteByte(request.ChatIndex)
	session.Conn.Write(p.ToBytes())
}
