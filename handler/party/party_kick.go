package party

import (
	"github.com/ferdoran/go-sro-agent-server/service"
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
	partyService := service.GetPartyServiceInstance()
	for {
		data := <-h.channel
		uniqueId, err := data.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read unique id")
		}

		partyService.KickPlayer(uniqueId, data.UserContext.UniqueID)
	}
}
