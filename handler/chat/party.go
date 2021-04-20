package chat

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func handlePartyMessage(request MessageRequest, session *server.Session) {
	world := service.GetWorldServiceInstance()
	player, err := world.GetPlayerByUniqueId(session.UserContext.UniqueID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to send party message"))
	}
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
		recievingPlayer, err := world.GetPlayerByUniqueId(v)
		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to send party message to party member "+string(v)))
			continue
		}
		recievingPlayer.Session.Conn.Write(p1.ToBytes())
	}
}
