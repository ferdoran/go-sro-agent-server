package party

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type PartyMatchingPlayerJoinRequestHandler struct {
	channel chan server.PacketChannelData
}

func InitPartyMatchingPlayerJoinRequestHandler() {
	handler := PartyMatchingPlayerJoinRequestHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PartyMatchingPlayerJoinRequest)}
	go handler.Handle()
}

func (h *PartyMatchingPlayerJoinRequestHandler) Handle() {
	partyService := service.GetPartyServiceInstance()
	for {
		data := <-h.channel
		requestId, err := data.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read request id")
		}

		playerJid, err1 := data.ReadUInt32()
		if err1 != nil {
			log.Panicln("Failed to read player jid")
		}

		acceptRequest, err2 := data.ReadBool()
		if err2 != nil {
			log.Panicln("Failed to read acceptRequest")
		}

		partyService.AnswerJoinRequest(requestId, playerJid, data.UserContext.UserID, acceptRequest)
	}
}
