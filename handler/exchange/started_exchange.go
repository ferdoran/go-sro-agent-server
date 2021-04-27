package exchange

import (
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-agent-server/service"
	log "github.com/sirupsen/logrus"
)

type ExchangeStartedHandler struct {
	channel chan server.PacketChannelData
}

func InitExchangeStartedHandler() {
	handler := ExchangeStartedHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PlayerInvitationRequest)}
	go handler.Handle()
}

func (h *ExchangeStartedHandler) Handle() {
	for {
		data := <-h.channel

		response, err := data.ReadByte()
		if err != nil {
			log.Panicln("Failed to read response")
		}

		unknown, err := data.ReadByte()	// Unknown, maybe type of the invitation?
		if err != nil {
			log.Panicln("Failed to read second byte",unknown)
		}

		if response == 1 {
			exchangeService := service.GetExchangeServiceInstance()
			exchangeService.AnswerStartExchangeRequest(data.Session, data.Session.UserContext.UniqueID)
		}
	}
}