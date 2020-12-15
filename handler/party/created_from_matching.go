package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
)

func SendPartyCreatedFromMatchingResponse(party model.Party, joinedPlayer uint32) {
	p := network.EmptyPacket()
	p.MessageID = opcode.PartyCreatedFromMatchingResponse
	p.WriteByte(255)
	p.WriteUInt32(party.Number)
	p.WriteUInt32(1)
	p.WriteByte(5)
	p.WriteByte(party.MemberCount)
	for _, v := range party.Members {
		player := model.GetSroWorldInstance().PlayersByUniqueId[v]
		p.WriteByte(255)
		p.WriteUInt32(2)
		p.WriteString(player.CharName)
		p.WriteUInt32(1907)
		p.WriteByte(10)
		p.WriteByte(170)
		p.WriteUInt16(25000)
		p.WriteUInt16(1007)
		p.WriteUInt16(65534)
		p.WriteUInt16(1710)
		p.WriteUInt32(65537)
		p.WriteUInt16(0)
		p.WriteByte(4)
		p.WriteUInt32(290)
		p.WriteUInt32(0)
	}
	for _, v := range party.Members {
		if v != party.MasterUniqueID {
			player := model.GetSroWorldInstance().PlayersByUniqueId[v]
			player.Session.Conn.Write(p.ToBytes())
		}
	}
	// Send out a separate package only containing the pt master
	partyMaster := model.GetSroWorldInstance().PlayersByUniqueId[party.MasterUniqueID]
	p1 := network.EmptyPacket()
	p1.MessageID = opcode.PartyCreatedFromMatchingResponse
	p1.WriteByte(255)
	p1.WriteUInt32(party.Number)
	p1.WriteUInt32(1)
	p1.WriteByte(5)
	p1.WriteByte(party.MemberCount)
	p1.WriteByte(255)
	p1.WriteUInt32(2)
	p1.WriteString(partyMaster.CharName)
	p1.WriteUInt32(1907)
	p1.WriteByte(10)
	p1.WriteByte(170)
	p1.WriteUInt16(25000)
	p1.WriteUInt16(1007)
	p1.WriteUInt16(65534)
	p1.WriteUInt16(1710)
	p1.WriteUInt32(65537)
	p1.WriteUInt16(0)
	p1.WriteByte(4)
	p1.WriteUInt32(290)
	p1.WriteUInt32(0)
	partyMaster.Session.Conn.Write(p1.ToBytes())

	// Send out another package containing the new user
	newPartyMember := model.GetSroWorldInstance().PlayersByUniqueId[joinedPlayer]
	p2 := network.EmptyPacket()
	p2.MessageID = opcode.PartyUpdateResponse
	p2.WriteByte(2)
	p2.WriteByte(255)
	p2.WriteUInt32(party.Number)
	p2.WriteString(newPartyMember.CharName)
	p2.WriteUInt32(1907)
	p2.WriteByte(5)
	p2.WriteByte(170)
	p2.WriteUInt16(25000)
	p2.WriteUInt16(1007)
	p2.WriteUInt16(65534)
	p2.WriteUInt16(1710)
	p2.WriteUInt32(65537)
	p2.WriteUInt16(0)
	p2.WriteByte(4)
	p2.WriteUInt32(290)
	p2.WriteUInt32(65537)
	partyMaster.Session.Conn.Write(p2.ToBytes())
}
