package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type PartyMatchingListHandler struct {
}

func NewPartyMatchingListHandler() server.PacketHandler {
	handler := PartyMatchingListHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PartyMatchingListRequest, handler)
	return handler
}

func (h PartyMatchingListHandler) Handle(data server.PacketChannelData) {
	pageIndex, err := data.ReadByte()
	if err != nil {
		log.Panicln("Failed to read page index")
	}

	partyCount := len(model.Parties)
	player := model.GetSroWorldInstance().PlayersByUniqueId[data.UserContext.UniqueID]

	var isPartyMember byte
	if player.Party.Number != 0 {
		isPartyMember = 2
	} else {
		if partyCount > 0 {
			isPartyMember = 1
		} else {
			isPartyMember = 0
		}
	}

	p := network.EmptyPacket()
	p.MessageID = opcode.PartyMatchingListResponse
	p.WriteByte(1)
	p.WriteByte(1) // pageCount TODO: calculate
	p.WriteByte(pageIndex)
	p.WriteByte(isPartyMember)
	for _, v := range model.Parties {
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
		if player.Party.Number == v.Number {
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
	}
	data.Session.Conn.Write(p.ToBytes())
}
