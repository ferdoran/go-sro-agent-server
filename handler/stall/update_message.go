package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
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
