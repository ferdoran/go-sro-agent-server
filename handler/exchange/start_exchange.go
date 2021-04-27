package exchange

import (
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-agent-server/model"
	log "github.com/sirupsen/logrus"
)

type ExchangeStartHandler struct {
	channel chan server.PacketChannelData
}

func InitExchangeStartHandler() {
	handler := ExchangeStartHandler{channel: server.PacketManagerInstance.GetQueue(opcode.ExchangeStartRequest)}
	go handler.Handle()
}

func (h *ExchangeStartHandler) Handle() {
	for {
		data := <-h.channel

		entityUniqueId, err := data.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read entityUniqueId")
		}

		entity := model.EntitySelectRequest {
			EntityUniqueID: entityUniqueId,
		}
		entitySelectService := service.EntitySelectService{}
		entitySelectErr     := entitySelectService.GetEntity(entity)

		if entitySelectErr != nil {
			log.Warnln(entitySelectErr)
		} else {
			if entitySelectService.IsPlayerCharacter {
				exchangeService := service.GetExchangeServiceInstance()
				exchangeService.AskStartExchangeRequest(data.Session.UserContext.UniqueID, entitySelectService.EntityUniqueID)
			} else {
				log.Warnf("Exchange started with non player entity: %#v\n", entitySelectService)
			}
		}
	}
}