package exchange

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type ExchangeApproveHandler struct {
	channel chan server.PacketChannelData
}

func InitExchangeApproveHandler() {
	handler := ExchangeApproveHandler{channel: server.PacketManagerInstance.GetQueue(opcode.ExchangeApproveRequest)}
	go handler.Handle()
}

func (h *ExchangeApproveHandler) Handle() {
	for {
		data := <-h.channel

		log.Debugln("Approving exchange")
		//TODO: Compute result
		p := network.EmptyPacket()
		p.MessageID = opcode.ExchangeApproveResponse
		p.WriteByte(1)               		 // Result
		data.Session.Conn.Write(p.ToBytes())
	}
}