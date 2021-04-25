package entity_select

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-agent-server/model"
	log "github.com/sirupsen/logrus"
)

type EntitySelectHandler struct {
	channel chan server.PacketChannelData
}

func InitEntitySelectHandler() {
	handler := EntitySelectHandler{channel: server.PacketManagerInstance.GetQueue(opcode.EntitySelectRequest)}
	go handler.Handle()
}

func (h *EntitySelectHandler) Handle() {
	data := <-h.channel

	entityUniqueId, err := data.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read entityUniqueId")
	}

	entity := model.EntitySelectRequest {
		EntityUniqueID: entityUniqueId,
	}

	entitySelectService := service.EntitySelectService{}
	entitySelectErr := entitySelectService.GetEntity(entity)
	if entitySelectErr != nil {
		log.Panicln(entitySelectErr)
	} else {
		p := network.EmptyPacket()
		p.MessageID = opcode.EntitySelectResponse
		p.WriteByte(1)               	// Result
		p.WriteUInt32(entityUniqueId)	// UniqueID from the request
		if entitySelectService.IsPlayerCharacter {
			p.WriteByte(1)               	// 
			p.WriteByte(5)               	// 
			p.WriteByte(4)               	// 
		} else if entitySelectService.IsNPCNpc {
			p.WriteByte(1)   				// Blacksmith JG hardcoded now
			p.WriteByte(4)   				// 
			p.WriteByte(1)   				// 
			p.WriteByte(2)   				// 
			p.WriteByte(4)   				// 
			p.WriteByte(20)  				// 
			p.WriteByte(1)   				// 
			p.WriteUInt16(0) 				// 
		} else if entitySelectService.IsNPCMob {
			p.WriteByte(1)    				// Mangyang & Weasel values
			p.WriteUInt32(36) 				// 
			p.WriteByte(1)    				// 
			p.WriteByte(5)    				// 
		}
		data.Session.Conn.Write(p.ToBytes())
	}
}