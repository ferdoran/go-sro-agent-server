package service

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	PartiesPerPage = 20
)

type PartyService struct {
	Parties            map[uint32]*model.Party
	FormedParties      map[uint32]*model.Party
	JoinRequests       map[uint32]*model.PartyJoinRequest
	partyNumberCounter uint32
	joinRequestCounter uint32
	mutex              sync.Mutex
}

var partyServiceInstance *PartyService
var partyServiceOnce sync.Once

func GetPartyServiceInstance() *PartyService {
	partyServiceOnce.Do(func() {
		partyServiceInstance = &PartyService{
			Parties:            make(map[uint32]*model.Party),
			FormedParties:      make(map[uint32]*model.Party),
			JoinRequests:       make(map[uint32]*model.PartyJoinRequest),
			partyNumberCounter: 1,
			joinRequestCounter: 1,
		}
	})

	return partyServiceInstance
}

func (p *PartyService) FormParty(partyFormRequest model.PartyFormRequest) uint32 {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	player, err := worldServiceInstance.GetPlayerByUniqueId(partyFormRequest.MasterUniqueID)

	if err != nil {
		logrus.Panic(errors.Wrap(err, "failed to form party"))
	}

	if player.HasParty() {
		logrus.Panicf("cannot form party because player %s is already in a party", player.GetName())
	}

	party := &model.Party{
		Number:            p.partyNumberCounter,
		MasterUniqueID:    partyFormRequest.MasterUniqueID,
		MasterJID:         partyFormRequest.MasterJID,
		MasterName:        partyFormRequest.MasterName,
		CountryType:       partyFormRequest.CountryType,
		MemberCount:       1,
		PartySettingsFlag: partyFormRequest.PartySettingsFlag,
		PurposeType:       partyFormRequest.PurposeType,
		LevelMin:          partyFormRequest.LevelMin,
		LevelMax:          partyFormRequest.LevelMax,
		Title:             partyFormRequest.Title,
		Members:           []uint32{partyFormRequest.MasterUniqueID},
		Mutex:             &sync.Mutex{},
	}
	p.partyNumberCounter++
	p.Parties[party.Number] = party
	p.FormedParties[party.Number] = party
	err = player.AddToParty(party)

	if err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to create party %d for player %s", party.Number, party.MasterName)))
		delete(p.Parties, party.Number)
	}

	return party.Number
}

func (p *PartyService) DeletePartyFromMatching(partyNumber, playerUniqueId uint32) {
	party, err := p.GetParty(partyNumber)

	if err != nil {
		logrus.Warnf("failed to delete party %d from matching: %s", partyNumber, err.Error())
		return
	}
	partyMaster, err := worldServiceInstance.GetPlayerByUniqueId(party.MasterUniqueID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to delete party"))
		return
	}

	if party.MasterUniqueID == playerUniqueId {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		delete(p.FormedParties, partyNumber)
		p := network.EmptyPacket()
		p.MessageID = opcode.PartyMatchingDeleteResponse
		p.WriteBool(true)
		p.WriteUInt32(partyNumber)
		partyMaster.GetSession().Conn.Write(p.ToBytes())
	} else {
		logrus.Warnf("failed to delete party %d from matching: player %d is not the owner", partyNumber, playerUniqueId)
	}
}

func (p *PartyService) GetParty(partyNumber uint32) (*model.Party, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if party, exists := p.Parties[partyNumber]; exists {
		return party, nil
	}

	return nil, errors.New(fmt.Sprintf("party %d does not exist", partyNumber))
}

func (p *PartyService) GetFormedParties(playerUniqueId uint32, page byte) []*model.Party {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	formedParties := make([]*model.Party, 0)

	for _, formedParty := range p.FormedParties {
		formedParties = append(formedParties, formedParty)
	}

	minIndex := (page - 1) * PartiesPerPage
	maxIndex := page * PartiesPerPage

	if len(formedParties) == 0 {
		return formedParties
	}

	if int(minIndex) >= len(formedParties) {
		minIndex = 0
	}

	if int(maxIndex) >= len(formedParties) {
		maxIndex = byte(len(formedParties) - 1)
	}

	player, err := worldServiceInstance.GetPlayerByUniqueId(playerUniqueId)
	if err != nil {
		logrus.Warnf("failed to add players aprty to formed parties")
		return formedParties
	}

	parties := formedParties[minIndex:maxIndex]
	if player.HasParty() {
		parties = append([]*model.Party{player.GetParty()}, parties...)
	}

	return parties[minIndex:maxIndex]
}

func (p *PartyService) GetFormedPartiesPageCount() byte {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	partyCount := len(p.FormedParties)
	if partyCount%PartiesPerPage != 0 {
		return 1 + byte(partyCount/PartiesPerPage)
	} else {
		return byte(partyCount / PartiesPerPage)
	}

}

func (p *PartyService) JoinFormedParty(partyNumber, playerUniqueId uint32) {
	party, err := p.GetParty(partyNumber)

	if err != nil {
		logrus.Error(errors.Wrap(err, "cannot join formed party"))
		return
	}

	player, err := worldServiceInstance.GetPlayerByUniqueId(playerUniqueId)

	if err != nil {
		logrus.Error(errors.Wrap(err, "cannot join formed party"))
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	joinRequest := &model.PartyJoinRequest{
		PartyNumber:    partyNumber,
		RequestID:      p.joinRequestCounter,
		PlayerJID:      player.ID,
		PlayerUniqueID: player.GetUniqueID(),
		AcceptResult:   false,
		Mutex:          &sync.Mutex{},
	}
	p.JoinRequests[joinRequest.RequestID] = joinRequest
	p.joinRequestCounter++

	partyMaster, err := worldServiceInstance.GetPlayerByUniqueId(party.MasterUniqueID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "cannot send join request to party master"))
		delete(p.JoinRequests, joinRequest.RequestID)
		return
	}

	packet := network.EmptyPacket()
	packet.MessageID = opcode.PartyMatchingPlayerJoinResponse
	packet.WriteUInt32(joinRequest.RequestID) // RequestID
	packet.WriteUInt32(uint32(player.ID))     // PlayerJID
	packet.WriteUInt32(party.Number)
	packet.WriteUInt32(0)                 // PlayerMasteryPrimaryID
	packet.WriteUInt32(0)                 // PlayerMasterySecondaryID
	packet.WriteByte(4)                   // unkByte01
	packet.WriteByte(255)                 // unkByte01
	packet.WriteUInt32(uint32(player.ID)) // PlayerJID_x2
	packet.WriteString(player.CharName)
	packet.WriteUInt32(1907)  // char model id
	packet.WriteByte(5)       // level
	packet.WriteByte(170)     // current hp/mp
	packet.WriteUInt16(25000) // regionId
	packet.WriteUInt16(1007)  // x
	packet.WriteUInt16(65534) // y
	packet.WriteUInt16(1710)  // z
	packet.WriteUInt32(65537) // unknown
	packet.WriteUInt16(0)     // guild name length
	packet.WriteByte(4)       // unknown
	packet.WriteUInt32(257)   // primary skill tree
	packet.WriteUInt32(290)   // secondary skill tree
	partyMaster.Session.Conn.Write(packet.ToBytes())
}

func (p *PartyService) AnswerJoinRequest(requestId, joiningPlayerJID, acceptingPlayerJID uint32, acceptCode bool) {
	joinRequest, err := p.GetJoinReqeust(requestId)

	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to accept join request"))
		return
	}

	party, err := p.GetParty(joinRequest.PartyNumber)

	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to accept join request"))
		return
	}

	if acceptingPlayerJID != party.MasterJID {
		// only party master can accept
		logrus.Error(errors.New("only party master can accept party join requests"))
		return
	}

	if joiningPlayerJID != uint32(joinRequest.PlayerJID) {
		logrus.Error(errors.New("joining player id must match player id for join request"))
		return
	}

	joinRequest.AcceptResult = acceptCode

	requestingPlayer, err := worldServiceInstance.GetPlayerByUniqueId(joinRequest.PlayerUniqueID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to send join request response to requesting player"))
		return
	}

	if joinRequest.AcceptResult {
		// add player to party
		party.Members = append(party.Members, requestingPlayer.GetUniqueID())
		party.MemberCount++

		// inform player
		partyJoinResponse := network.EmptyPacket()
		partyJoinResponse.MessageID = opcode.PartyMatchingJoinResponse
		partyJoinResponse.WriteByte(1)
		partyJoinResponse.WriteUInt16(1) // AcceptResult is somehow sent as uint16
		requestingPlayer.Session.Conn.Write(partyJoinResponse.ToBytes())

		// send member count to newly joined player
		partyMemberCountResponse := network.EmptyPacket()
		partyMemberCountResponse.MessageID = opcode.PartyMemberCountResponse
		partyMemberCountResponse.WriteByte(1)
		partyMemberCountResponse.WriteUInt32(uint32(party.MemberCount))
		requestingPlayer.Session.Conn.Write(partyMemberCountResponse.ToBytes())

		if party.MemberCount == 2 {
			// new party was created
			if ptMaster, err := worldServiceInstance.GetPlayerByUniqueId(party.MasterUniqueID); err != nil {
				logrus.Error(errors.Wrap(err, "failed to inform party master about created party"))
			} else {
				partyCreatedResponse := network.EmptyPacket()
				partyCreatedResponse.MessageID = opcode.PartyCreateResponse
				partyCreatedResponse.WriteByte(1)
				partyCreatedResponse.WriteUInt32(1)
				ptMaster.Session.Conn.Write(partyCreatedResponse.ToBytes())

				p1 := network.EmptyPacket()
				p1.MessageID = opcode.PartyCreatedFromMatchingResponse
				p1.WriteByte(0xFF) // splitter
				p1.WriteUInt32(party.Number)
				p1.WriteUInt32(uint32(ptMaster.ID))            // MasterJID
				p1.WriteByte(party.PartySettingsFlag.ToByte()) // type?
				p1.WriteByte(party.MemberCount)
				p1.WriteByte(255)                   // splitter
				p1.WriteUInt32(uint32(ptMaster.ID)) // MemberJID
				p1.WriteString(ptMaster.GetName())
				p1.WriteUInt32(1907)               // char model id
				p1.WriteByte(byte(ptMaster.Level)) // level
				p1.WriteByte(170)
				p1.WriteUInt16(uint16(ptMaster.GetNavmeshPosition().Region.ID))         // regionId
				p1.WriteUInt16(uint16(ptMaster.GetNavmeshPosition().Offset.X))          // x
				p1.WriteUInt16(uint16(ptMaster.GetNavmeshPosition().Offset.Y + 0xFFFF)) // y
				p1.WriteUInt16(uint16(ptMaster.GetNavmeshPosition().Offset.Z))          // z
				p1.WriteUInt32(65537)                                                   // unknown
				p1.WriteUInt16(0)                                                       // guild name length
				p1.WriteByte(4)                                                         // unknown
				p1.WriteUInt32(290)                                                     // primary skill tree
				p1.WriteUInt32(0)                                                       // secondary skill tree
				ptMaster.Session.Conn.Write(p1.ToBytes())

				p2 := network.EmptyPacket()
				p2.MessageID = opcode.PartyUpdateResponse
				p2.WriteByte(2)   // party action type | 1 = close pt | 2 = player joined | 3 = left pt or kicked | 6 = player info | 9 = new pt master
				p2.WriteByte(255) // splitter
				p2.WriteUInt32(uint32(requestingPlayer.ID))
				p2.WriteString(requestingPlayer.CharName)
				p2.WriteUInt32(1907)                                                            // Char model id
				p2.WriteByte(5)                                                                 // level
				p2.WriteByte(170)                                                               // HP / MP
				p2.WriteUInt16(uint16(requestingPlayer.GetNavmeshPosition().Region.ID))         // regionId
				p2.WriteUInt16(uint16(requestingPlayer.GetNavmeshPosition().Offset.X))          // x
				p2.WriteUInt16(uint16(requestingPlayer.GetNavmeshPosition().Offset.Y + 0xFFFF)) // y
				p2.WriteUInt16(uint16(requestingPlayer.GetNavmeshPosition().Offset.Z))          // z
				p2.WriteUInt32(65537)                                                           // unknown
				p2.WriteUInt16(0)                                                               // guild name length
				p2.WriteByte(4)                                                                 // unknown
				p2.WriteUInt32(290)                                                             // primary skill tree
				p2.WriteUInt32(65537)                                                           // secondary skill tree
				ptMaster.Session.Conn.Write(p2.ToBytes())
			}
		}
	} else {
		// inform player
		partyJoinResponse := network.EmptyPacket()
		partyJoinResponse.MessageID = opcode.PartyMatchingJoinResponse
		partyJoinResponse.WriteByte(1)
		partyJoinResponse.WriteUInt16(0) // AcceptResult is somehow sent as uint16
		requestingPlayer.Session.Conn.Write(partyJoinResponse.ToBytes())
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.JoinRequests, joinRequest.RequestID)
}

func (p *PartyService) GetJoinReqeust(requestId uint32) (*model.PartyJoinRequest, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if joinRequest, exists := p.JoinRequests[requestId]; exists {
		return joinRequest, nil
	} else {
		return nil, errors.New(fmt.Sprintf("party join request %d does not exist", requestId))
	}
}

func (p *PartyService) UpdateParty(changedParty model.Party) {
	party, err := p.GetParty(changedParty.Number)

	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to update party"))
		return
	}

	if party.MasterUniqueID != changedParty.MasterUniqueID {
		logrus.Warnf("failed to update party: party masters must match")
		return
	}

	ptMaster, err := worldServiceInstance.GetPlayerByUniqueId(party.MasterUniqueID)

	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to update party"))
		return
	}

	party.LevelMin = changedParty.LevelMin
	party.LevelMax = changedParty.LevelMax
	party.Title = changedParty.Title

	partyMatchingUpdateResponse := network.EmptyPacket()
	partyMatchingUpdateResponse.MessageID = opcode.PartyMatchingUpdateResponse
	partyMatchingUpdateResponse.WriteByte(1)
	partyMatchingUpdateResponse.WriteUInt32(party.Number)
	partyMatchingUpdateResponse.WriteUInt32(2)
	partyMatchingUpdateResponse.WriteByte(party.PartySettingsFlag.ToByte())
	partyMatchingUpdateResponse.WriteByte(party.PurposeType.ToByte())
	partyMatchingUpdateResponse.WriteByte(party.LevelMin)
	partyMatchingUpdateResponse.WriteByte(party.LevelMax)
	partyMatchingUpdateResponse.WriteString(party.Title)
	ptMaster.Session.Conn.Write(partyMatchingUpdateResponse.ToBytes())
}

func (p *PartyService) KickPlayer(playerToKickUniqueId, requestingPlayerUniqueId uint32) {
	player, err := worldServiceInstance.GetPlayerByUniqueId(playerToKickUniqueId)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to kick player"))
		return
	}

	ptMaster, err := worldServiceInstance.GetPlayerByUniqueId(player.GetParty().MasterUniqueID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to kick player: could not find pt master"))
		return
	}

	if ptMaster.GetUniqueID() != requestingPlayerUniqueId {
		logrus.Error(errors.Wrap(err, "failed to kick player: requesting player must be master"))
		return
	}

	if !player.HasParty() {
		logrus.Warn("player must be in party to be kicked")
		return
	}

	if !ptMaster.HasParty() {
		logrus.Warn("party master does not have a party")
		return
	}

	party := player.GetParty()

	partyLeftPacket := network.EmptyPacket()
	partyLeftPacket.MessageID = opcode.PartyUpdateResponse
	partyLeftPacket.WriteByte(3) // party left
	partyLeftPacket.WriteUInt32(uint32(player.ID))
	partyLeftPacket.WriteByte(4) // leave type

	for index, memberUniqueID := range party.Members {
		member, _ := worldServiceInstance.GetPlayerByUniqueId(memberUniqueID)
		member.GetSession().Conn.Write(partyLeftPacket.ToBytes())
		if memberUniqueID == player.UniqueID {
			party.Members = append(party.Members[:index], party.Members[index+1:]...)
		}
	}
}
