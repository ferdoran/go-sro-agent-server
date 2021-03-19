package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type PartyKickHandler struct {
	channel chan server.PacketChannelData
}

func InitPartyKickHandler() {
	handler := PartyKickHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PartyKickRequest)}
	go handler.Handle()
}

func (h *PartyKickHandler) Handle() {
	for {
		data := <-h.channel
		uniqueId, err := data.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read unique id")
		}

		partyMaster, ok := model.GetSroWorldInstance().PlayersByUniqueId[data.UserContext.UniqueID]
		if !ok {
			log.Panicln("Party master not found")
		}
		playerToKick, ok1 := model.GetSroWorldInstance().PlayersByUniqueId[uniqueId]
		if !ok1 {
			log.Panicln("Player to kick not found")
		}

		if partyMaster.Party.Number == playerToKick.Party.Number && partyMaster.UniqueID == partyMaster.Party.MasterUniqueID {
			for _, v := range partyMaster.Party.Members {
				var members []uint32
				if v != playerToKick.UniqueID {
					members = append(members, v)
				}
			}
			partyMaster.Party.MemberCount -= 1
			partyMaster.Party.UpdateParty(partyMaster.UniqueID)
			for _, v := range partyMaster.Party.Members {
				p, _ := model.GetSroWorldInstance().PlayersByUniqueId[v]
				p.Party = partyMaster.Party
			}
		}

		p := network.EmptyPacket()
		p.MessageID = opcode.PartyUpdateResponse
		p.WriteByte(1)
		p.WriteUInt16(176)
		for _, v := range partyMaster.Party.Members {
			player, _ := model.GetSroWorldInstance().PlayersByUniqueId[v]
			player.Session.Conn.Write(p.ToBytes())
		}
	}
}
