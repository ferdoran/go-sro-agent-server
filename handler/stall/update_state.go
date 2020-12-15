package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

func updateState(data server.PacketChannelData) {
	isOpen, err := data.ReadBool()
	if err != nil {
		log.Panicln("Failed to read isOpen")
	}

	stallNetworkResult, err1 := data.ReadUInt16()
	if err1 != nil {
		log.Panicln("Failed to read stallNetworkResult")
	}

	p := network.EmptyPacket()
	p.MessageID = opcode.StallUpdateResponse
	p.WriteByte(1)
	p.WriteByte(StallState)
	p.WriteBool(isOpen)
	p.WriteUInt16(stallNetworkResult)
	data.Session.Conn.Write(p.ToBytes())
}
