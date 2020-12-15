package party

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type PartyMatchingJoinRequestHandler struct {
}

func NewPartyMatchingJoinRequestHandler() server.PacketHandler {
	handler := PartyMatchingJoinRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PartyMatchingJoinRequest, handler)
	return handler
}

func (h PartyMatchingJoinRequestHandler) Handle(data server.PacketChannelData) {
	partyNumber, err := data.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read party number")
	}

	model.CurrentRequestID++
	requestId := model.CurrentRequestID
	handler := &PartyMatchingPlayerJoinRequestHandler{}
	handler.AskMaster(data, requestId, partyNumber)
}

func (h PartyMatchingJoinRequestHandler) SendJoinResponse(requestId uint32) {
	if joinRequest, ok := model.GetJoinRequest(requestId); ok {
		p := network.EmptyPacket()
		p.MessageID = opcode.PartyMatchingJoinResponse
		p.WriteByte(1)
		p.WriteUInt16(joinRequest.AcceptCode)
		requestingPlayer := model.GetSroWorldInstance().PlayersByUniqueId[joinRequest.PlayerUniqueID]
		requestingPlayer.Session.Conn.Write(p.ToBytes())
		requestingPlayer.Session.Conn.Write(p.ToBytes())

		if hasJoined, party := joinRequest.CleanupJoinRequest(); hasJoined {
			SendMemberCountResponse(joinRequest.PlayerUniqueID, party.MemberCount)
			SendPartyCreateResponse(party.MasterUniqueID)
			SendPartyCreatedFromMatchingResponse(party, joinRequest.PlayerUniqueID)
		}
	}
}
