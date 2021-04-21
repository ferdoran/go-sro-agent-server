package chat

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func handleWhisperMessage(request MessageRequest, session *server.Session) {
	world := service.GetWorldServiceInstance()
	receivingPlayer, err := world.GetPlayerByCharName(request.Receiver)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to send private message"))
		return
	}
	p := network.EmptyPacket()
	p.MessageID = opcode.ChatResponse

	if receivingPlayer == nil {
		// Failed to send message
		logrus.Debugf("failed to send message from %s to %s\n", session.UserContext.CharName, request.Receiver)
		p.WriteByte(2)
		p.WriteUInt16(3)
		p.WriteByte(request.ChatType)
		p.WriteByte(request.ChatIndex)
		session.Conn.Write(p.ToBytes())
		return
	}

	p.WriteByte(1)
	p.WriteByte(request.ChatType)
	p.WriteByte(request.ChatIndex)
	session.Conn.Write(p.ToBytes())

	p2 := network.EmptyPacket()
	p2.MessageID = opcode.ChatUpdate
	p2.WriteByte(request.ChatType)
	p2.WriteString(session.UserContext.CharName)
	p2.WriteString(request.Message)
	receivingPlayer.Session.Conn.Write(p2.ToBytes())
}
