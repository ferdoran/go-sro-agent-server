package stall

import (
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type StallUpdateHandler struct {
	channel chan server.PacketChannelData
}

func InitStallUpdateHandler() {
	handler := StallUpdateHandler{channel: server.PacketManagerInstance.GetQueue(opcode.StallUpdateRequest)}
	go handler.Handle()
}

func (s *StallUpdateHandler) Handle() {
	for {
		data := <-s.channel
		updateType, err := data.ReadByte()
		if err != nil {
			log.Panicln("Failed to read update type")
		}
		switch updateType {
		case StallUpdateItem:
			updateItem(data)
		case StallAddItem:
			fallthrough
		case StallRemoveItem:
			addRemoveItem(data, updateType)
		case StallFleaMarketMode:
		case StallState:
			updateState(data)
		case StallMessage:
			updateMessage(data)
		case StallName:
			updateStallName(data)
		}
	}
}
