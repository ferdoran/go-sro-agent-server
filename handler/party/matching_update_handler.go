package party

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"sync"
)

type PartyMatchingUpdateHandler struct{}

func NewPartyMatchingUpdateHandler() server.PacketHandler {
	handler := PartyMatchingUpdateHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PartyMatchingUpdateRequest, handler)
	return handler
}

func (h PartyMatchingUpdateHandler) Handle(data server.PacketChannelData) {
	partyNumber, err := data.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read")
	}

	_, err1 := data.ReadUInt32() // TODO: what is this? Possible Placeholder for partynumber
	if err1 != nil {
		log.Panicln("Failed to read")
	}

	partySetting, err2 := data.ReadByte()
	if err2 != nil {
		log.Panicln("Failed to read partySetting")
	}

	purposeType, err3 := data.ReadByte()
	if err3 != nil {
		log.Panicln("Failed to read purposeType")
	}

	levelMin, err4 := data.ReadByte()
	if err4 != nil {
		log.Panicln("Failed to read levelMin")
	}

	levelMax, err5 := data.ReadByte()
	if err5 != nil {
		log.Panicln("Failed to read levelMax")
	}

	title, err6 := data.ReadString()
	if err6 != nil {
		log.Panicln("Failed to read title")
	}

	party := model.Party{
		MasterJID:         data.UserContext.UserID,
		MasterName:        data.UserContext.CharName,
		CountryType:       0, // TODO: Figure out
		PartySettingsFlag: model.PartySetting(partySetting),
		PurposeType:       model.PartyPurpose(purposeType),
		LevelMin:          levelMin,
		LevelMax:          levelMax,
		Title:             title,
		Mutex:             &sync.Mutex{},
	}

	party.UpdateParty(data.UserContext.UniqueID)

	p := network.EmptyPacket()
	p.MessageID = opcode.PartyMatchingUpdateResponse
	p.WriteByte(1)
	p.WriteUInt32(partyNumber)
	p.WriteUInt32(0)
	p.WriteByte(partySetting)
	p.WriteByte(purposeType)
	p.WriteByte(levelMin)
	p.WriteByte(levelMax)
	p.WriteString(title)
	data.Session.Conn.Write(p.ToBytes())
}
