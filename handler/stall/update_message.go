package stall

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

func updateMessage(data server.PacketChannelData) {
	message, err1 := data.ReadString()
	if err1 != nil {
		log.Panicln("Failed to read message")
	}

	p := network.EmptyPacket()
	p.MessageID = opcode.StallUpdateResponse
	p.WriteByte(1)
	p.WriteByte(StallMessage)
	p.WriteString(message)
	data.Session.Conn.Write(p.ToBytes())
}
