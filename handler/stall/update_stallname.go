package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
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
