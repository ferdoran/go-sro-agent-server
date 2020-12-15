package party

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

type PartyKickHandler struct {
}

func NewPartyKickHandler() server.PacketHandler {
	handler := PartyKickHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PartyKickRequest, handler)
	return handler
}

func (h PartyKickHandler) Handle(data server.PacketChannelData) {
	uniqueId, err := data.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read unique id")
	}

	partyMaster, ok   := model.GetSroWorldInstance().PlayersByUniqueId[data.UserContext.UniqueID];
    if !ok {
        log.Panicln("Party master not found")
    }
	playerToKick, ok1 := model.GetSroWorldInstance().PlayersByUniqueId[uniqueId];
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