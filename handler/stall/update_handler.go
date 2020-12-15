package stall

import (
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type StallUpdateHandler struct {
}

func NewStallUpdateHandler() server.PacketHandler {
	handler := StallUpdateHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.StallUpdateRequest, handler)
	return handler
}

func (s StallUpdateHandler) Handle(data server.PacketChannelData) {
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
