package stall

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

func updateStallName(data server.PacketChannelData) {
	stallName, err1 := data.ReadString()
	if err1 != nil {
		log.Panicln("Failed to read stall name")
	}

	p := network.EmptyPacket()
	p.MessageID = opcode.StallEntityNameResponse
	p.WriteUInt32(data.Session.UserContext.UniqueID)
	p.WriteString(stallName)
	data.Session.Conn.Write(p.ToBytes())
}
