package exchange

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type ExchangeCancelHandler struct {
	channel chan server.PacketChannelData
}

func InitExchangeCancelHandler() {
	handler := ExchangeCancelHandler{channel: server.PacketManagerInstance.GetQueue(opcode.ExchangeCancelRequest)}
	go handler.Handle()
}

func (h *ExchangeCancelHandler) Handle() {
	for {
		data := <-h.channel

		log.Debugln("Cancelling exchange")
		//TODO: Compute result
		p := network.EmptyPacket()
		p.MessageID = opcode.ExchangeCancelResponse
		p.WriteByte(1)               		 // Result
		data.Session.Conn.Write(p.ToBytes())
	}
}