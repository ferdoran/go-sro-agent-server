package party

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"sync"
)

type PartyMatchingPlayerJoinRequestHandler struct {
}

func NewPartyMatchingPlayerJoinRequestHandler() server.PacketHandler {
	handler := PartyMatchingPlayerJoinRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PartyMatchingPlayerJoinRequest, handler)
	return handler
}

func (h PartyMatchingPlayerJoinRequestHandler) Handle(data server.PacketChannelData) {
	requestId, err := data.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read request id")
	}

	_, err1 := data.ReadUInt32()
	if err1 != nil {
		log.Panicln("Failed to read player jid")
	}

	acceptCode, err2 := data.ReadByte()
	if err2 != nil {
		log.Panicln("Failed to read acceptCode")
	}

	if joinRequest, ok := model.GetJoinRequest(requestId); ok {
		joinRequest.AcceptCode = uint16(acceptCode)
		joinRequest.PutJoinRequest()
	}

	joinRequestHandler := &PartyMatchingJoinRequestHandler{}
	joinRequestHandler.SendJoinResponse(requestId)
}

func (h PartyMatchingPlayerJoinRequestHandler) AskMaster(data server.PacketChannelData, requestId, partyNumber uint32) {
	player := model.GetSroWorldInstance().PlayersByUniqueId[data.UserContext.UniqueID]

	var pt model.Party
	for _, v := range model.Parties {
		if v.Number == partyNumber {
			pt = v
			break
		}
	}

	joinRequest := model.JoinRequest{
		PartyNumber:    partyNumber,
		RequestID:      requestId,
		PlayerJID:      player.ID,
		PlayerUniqueID: player.UniqueID,
		AcceptCode:     0,
		Mutex:          &sync.Mutex{},
	}

	joinRequest.PutJoinRequest()

	ptMaster := model.GetSroWorldInstance().PlayersByUniqueId[pt.MasterUniqueID]
	p := network.EmptyPacket()
	p.MessageID = opcode.PartyMatchingPlayerJoinResponse
	p.WriteUInt32(joinRequest.RequestID) // RequestID
	p.WriteUInt32(uint32(player.ID))     // PlayerJID
	p.WriteUInt32(pt.Number)
	p.WriteUInt32(0)                 // PlayerMasteryPrimaryID
	p.WriteUInt32(0)                 // PlayerMasterySecondaryID
	p.WriteByte(4)                   // unkByte01
	p.WriteByte(255)                 // unkByte01
	p.WriteUInt32(uint32(player.ID)) // PlayerJID_x2
	p.WriteString(player.CharName)
	p.WriteUInt32(1907)  //
	p.WriteByte(5)       //
	p.WriteByte(170)     //
	p.WriteUInt16(25000) //
	p.WriteUInt16(1007)  //
	p.WriteUInt16(65534) //
	p.WriteUInt16(1710)  //
	p.WriteUInt32(65537) //
	p.WriteUInt16(0)     //
	p.WriteByte(4)       //
	p.WriteUInt32(257)   //
	p.WriteUInt32(290)   //
	ptMaster.Session.Conn.Write(p.ToBytes())
}
