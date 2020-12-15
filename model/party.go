package model

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type Party struct {
	Number            uint32
	MasterUniqueID    uint32
	MasterJID         uint32
	MasterName        string
	CountryType       byte
	MemberCount       byte
	PartySettingsFlag PartySetting
	PurposeType       PartyPurpose
	LevelMin          byte
	LevelMax          byte
	Title             string
	Members           []uint32
	Mutex             *sync.Mutex
}

type JoinRequest struct {
	PartyNumber    uint32
	RequestID      uint32
	PlayerJID      int
	PlayerUniqueID uint32
	AcceptCode     uint16
	Mutex          *sync.Mutex
}

var CurrentPartyNumber uint32 = 0
var Parties []Party

var CurrentRequestID uint32 = 0
var JoinRequests map[uint32]*JoinRequest = make(map[uint32]*JoinRequest)

func (p *Party) FormParty(uniqueId uint32) uint32 {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if player, ok := GetSroWorldInstance().PlayersByUniqueId[uniqueId]; ok {
		if player.Party.MasterName != "" {
			CurrentPartyNumber++
			player.Party.Number = CurrentPartyNumber
			Parties = append(Parties, player.Party)
			for _, v := range player.Party.Members {
				p, _ := GetSroWorldInstance().PlayersByUniqueId[v]
				p.Party = player.Party
			}
		} else {
			CurrentPartyNumber++
			var members []uint32
			members = append(members, uniqueId)
			p.Members = members
			p.Number = CurrentPartyNumber
			p.MemberCount = byte(len(p.Members))
			Parties = append(Parties, *p)
			player.Party = *p
		}
	} else {
		log.Panicln("Player:", uniqueId, " not found!")
	}

	return CurrentPartyNumber
}

func (p *Party) UpdateParty(uniqueId uint32) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if player, ok := GetSroWorldInstance().PlayersByUniqueId[uniqueId]; ok {
		if player.Party.MasterJID == p.MasterJID {
			var parties []Party
			for _, v := range Parties {
				if v.Number != p.Number {
					parties = append(parties, v)
				} else {
					parties = append(parties, *p)
				}
			}
			Parties = parties
			player.Party = *p
		} else {
			log.Panicln("Only master can update party!")
		}
	} else {
		log.Panicln("Player:", uniqueId, " not found!")
	}
}

func (p *Party) DeletePartyFromMatching(uniqueId uint32) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if player, ok := GetSroWorldInstance().PlayersByUniqueId[uniqueId]; ok {
		if player.Party.MasterJID == p.MasterJID {
			var parties []Party
			for _, v := range Parties {
				if v.Number != p.Number {
					parties = append(parties, v)
				}
			}
			Parties = parties
			player.Party.Number = 0
		} else {
			log.Panicln("Only master can delete party entry!")
		}
	} else {
		log.Panicln("Player:", uniqueId, " not found!")
	}
}

func (j *JoinRequest) PutJoinRequest() {
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	JoinRequests[j.RequestID] = j
}

func (j *JoinRequest) CleanupJoinRequest() (bool, Party) {
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	hasJoined := false
	var party Party

	if j.AcceptCode == 1 {
		if _, ok := GetSroWorldInstance().PlayersByUniqueId[j.PlayerUniqueID]; ok {
			var pt Party
			for _, v := range Parties {
				if v.Number == j.PartyNumber {
					pt = v
					break
				}
			}
			pt.Members = append(pt.Members, j.PlayerUniqueID)
			pt.MemberCount++
			for _, v := range pt.Members {
				p, _ := GetSroWorldInstance().PlayersByUniqueId[v]
				p.Party = pt
			}
			party = pt
			hasJoined = true;
		}
	}

	delete(JoinRequests, j.RequestID)

	return hasJoined, party
}

func GetJoinRequest(requestId uint32) (*JoinRequest, bool) {
	joinRequest, ok := JoinRequests[requestId]
	return joinRequest, ok
}
