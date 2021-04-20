package party

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type PartyMatchingListHandler struct {
	channel chan server.PacketChannelData
}

func InitPartyMatchingListHandler() {
	handler := PartyMatchingListHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PartyMatchingListRequest)}
	go handler.Handle()
}

func (h *PartyMatchingListHandler) Handle() {
	partyService := service.GetPartyServiceInstance()
	for {
		data := <-h.channel
		pageIndex, err := data.ReadByte()
		if err != nil {
			log.Panicln("Failed to read page index")
		}

		parties := partyService.GetFormedParties(data.UserContext.UniqueID, pageIndex)
		pageCount := partyService.GetFormedPartiesPageCount()

		partyCount := len(parties)

		p := network.EmptyPacket()
		p.MessageID = opcode.PartyMatchingListResponse
		p.WriteByte(1)
		p.WriteByte(pageCount) // pageCount TODO: calculate
		p.WriteByte(pageIndex)
		p.WriteByte(byte(partyCount))
		for _, v := range parties {
			p.WriteUInt32(v.Number)
			p.WriteUInt32(v.MasterJID)
			p.WriteString(v.MasterName)
			p.WriteByte(v.CountryType)
			p.WriteByte(v.MemberCount)
			p.WriteByte(v.PartySettingsFlag.ToByte())
			p.WriteByte(v.PurposeType.ToByte())
			p.WriteByte(v.LevelMin)
			p.WriteByte(v.LevelMax)
			p.WriteString(v.Title)
		}
		data.Session.Conn.Write(p.ToBytes())
	}
}
