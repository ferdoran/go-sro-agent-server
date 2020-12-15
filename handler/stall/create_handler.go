package stall

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

const (
	StallAvatarDefault        = 0
	StallAvatarMangyang       = 3847
	StallAvatarBigEyeGhost    = 3848
	StallAvatarEarthGhost     = 3849
	StallAvatarSpecialMonster = 3850
)

type StallCreateHandler struct {
}

func NewStallCreateHandler() server.PacketHandler {
	handler := StallCreateHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.StallCreateRequest, handler)
	return handler
}

func (s StallCreateHandler) Handle(data server.PacketChannelData) {
	stallName, _ := data.ReadString()
	/*if err != nil {
		log.Panicln("Failed to read stall name")
	}*/

	// TODO: Check stall name for validity?
	log.Println(stallName)
	p := network.EmptyPacket()
	p.MessageID = opcode.StallCreateResponse
	p.WriteByte(1)
	data.Session.Conn.Write(p.ToBytes())

	p1 := network.EmptyPacket()
	p1.MessageID = opcode.StallEntityCreateResponse
	p1.WriteUInt32(data.Session.UserContext.UniqueID)
	p1.WriteString(stallName)

	// TODO: Check which stall avatar is activated
	p1.WriteUInt32(StallAvatarDefault)
	data.Session.Conn.Write(p1.ToBytes())
}
