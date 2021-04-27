package exchange

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type ExchangeConfirmHandler struct {
	channel chan server.PacketChannelData
}

func InitExchangeConfirmHandler() {
	handler := ExchangeConfirmHandler{channel: server.PacketManagerInstance.GetQueue(opcode.ExchangeConfirmRequest)}
	go handler.Handle()
}

func (h *ExchangeConfirmHandler) Handle() {
	for {
		data := <-h.channel

		log.Debugln("Confirming exchange")
		//TODO: Compute result
		p := network.EmptyPacket()
		p.MessageID = opcode.ExchangeConfirmResponse
		p.WriteByte(1)               		 // Result
		data.Session.Conn.Write(p.ToBytes())
	}
}