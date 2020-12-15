package party

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"sync"
)

type PartyMatchingDeleteHandler struct {
}

func NewPartyMatchingDeleteHandler() server.PacketHandler {
	handler := PartyMatchingDeleteHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PartyMatchingDeleteRequest, handler)
	return handler
}

func (h PartyMatchingDeleteHandler) Handle(data server.PacketChannelData) {
	partyNumber, err := data.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read partyNumber")
	}

	party := model.Party{
		Number:    partyNumber,
		MasterJID: data.UserContext.UserID,
		Mutex:     &sync.Mutex{},
	}
	party.DeletePartyFromMatching(data.UserContext.UniqueID)

	p := network.EmptyPacket()
	p.MessageID = opcode.PartyMatchingDeleteResponse
	p.WriteByte(1)
	p.WriteUInt32(partyNumber)
	data.Session.Conn.Write(p.ToBytes())
}
