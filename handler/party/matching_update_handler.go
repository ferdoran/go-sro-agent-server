package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
	"sync"
)

type PartyMatchingUpdateHandler struct {
	channel chan server.PacketChannelData
}

func InitPartyMatchingUpdateHandler() {
	handler := PartyMatchingUpdateHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PartyMatchingUpdateRequest)}
	go handler.Handle()
}

func (h *PartyMatchingUpdateHandler) Handle() {
	partyService := service.GetPartyServiceInstance()
	for {
		data := <-h.channel
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
			Number:            partyNumber,
			MasterJID:         data.UserContext.UserID,
			MasterName:        data.UserContext.CharName,
			MasterUniqueID:    data.UserContext.UniqueID,
			CountryType:       0, // TODO: Figure out
			PartySettingsFlag: model.PartySetting(partySetting),
			PurposeType:       model.PartyPurpose(purposeType),
			LevelMin:          levelMin,
			LevelMax:          levelMax,
			Title:             title,
			Mutex:             &sync.Mutex{},
		}

		partyService.UpdateParty(party)

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
}
