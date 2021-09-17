package character

import (
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/sirupsen/logrus"
)

type SelectEntityHandler struct {
	channel chan server.PacketChannelData
}

func InitSelectEntityHandler() {
	queue := server.PacketManagerInstance.GetQueue(opcode.EntitySelectRequest)
	handler := SelectEntityHandler{channel: queue}
	go handler.Handle()
}

func (h *SelectEntityHandler) Handle() {
	for {
		data := <-h.channel
		// TODO implement real logic
		uniqueId, err := data.ReadUInt32()
		if err != nil {
			logrus.Panicf("failed to read unique id")
		}

		world := service.GetWorldServiceInstance()
		object, err := world.GetObjectByUniqueId(uniqueId)

		p := network.EmptyPacket()
		p.MessageID = opcode.EntitySelectResponse
		if err != nil {
			logrus.Error(err)
			p.WriteByte(0)
			p.WriteByte(0)
			data.Conn.Write(p.ToBytes())
			return
		}
		p.WriteByte(1)
		p.WriteUInt32(uniqueId)
		if object.GetTypeInfo().IsNPCNpc() {
			// TODO: depends to talk options, so ignore for now
			p.WriteByte(0)
			return
		} else if object.GetTypeInfo().IsNPCMob() {
			p.WriteByte(1)
			p.WriteUInt32(0) // TODO: Monster HP
			p.WriteByte(1)
			p.WriteByte(5)

		} else if object.GetTypeInfo().IsPlayerCharacter() {
			p.WriteUInt32(0)
			p.WriteByte(0) // TODO: Trader Level
			p.WriteByte(0) // TODO: Hunter Level
			p.WriteByte(0) // TODO: Thief Level
			p.WriteByte(0)
		}
		data.Conn.Write(p.ToBytes())
	}
}
