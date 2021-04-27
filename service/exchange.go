package service

import (
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-agent-server/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ExchangeUpdateType byte
const (
	Item ExchangeUpdateType = 1
	Gold ExchangeUpdateType = 2
)

type ExchangeService struct {
	Exchanges map[uint32]*model.ExchangeRequest
	mutex     sync.Mutex
}

var exchangeServiceInstance *ExchangeService
var exchangeServiceOnce sync.Once

func GetExchangeServiceInstance() *ExchangeService {
	exchangeServiceOnce.Do(func() {
		exchangeServiceInstance = &ExchangeService{
			Exchanges: make(map[uint32]*model.ExchangeRequest),
		}
	})
	return exchangeServiceInstance
}

func (p *ExchangeService) AskStartExchangeRequest(requestingPlayerUniqueID, requestedPlayerUniqueID uint32) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	exchange := &model.ExchangeRequest {
		RequestingUniqueID: requestingPlayerUniqueID,
		RequestedUniqueID:  requestedPlayerUniqueID,
		IsStarted:          false,
		Mutex:              &sync.Mutex{},
	}

	p.Exchanges[requestingPlayerUniqueID] = exchange
	p.Exchanges[requestedPlayerUniqueID]  = exchange

	player, err := worldServiceInstance.GetPlayerByUniqueId(requestedPlayerUniqueID)
	if err != nil {
		log.Debugln(err)
		log.Warnf("Requested player not found: %x\nRequesting player was: %x\n", requestedPlayerUniqueID, requestingPlayerUniqueID)
	} else {
		p := network.EmptyPacket()
		p.MessageID = opcode.PlayerInvitationResponse
		p.WriteByte(1)
		p.WriteUInt32(requestingPlayerUniqueID)
		player.Session.Conn.Write(p.ToBytes())
	}
}

func (p *ExchangeService) AnswerStartExchangeRequest(session *server.Session, requestedPlayerUniqueID uint32) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if exchange, exists := p.Exchanges[requestedPlayerUniqueID]; exists {
		player, err := worldServiceInstance.GetPlayerByUniqueId(exchange.RequestingUniqueID)
		if err != nil {
			log.Debugln(err)
			log.Warnf("Requested player not found: %x\nRequesting player was: %x\n", requestedPlayerUniqueID, exchange.RequestingUniqueID)
		} else {
			p := network.EmptyPacket()
			p.MessageID = opcode.ExchangeStartedResponse
			p.WriteUInt32(exchange.RequestingUniqueID)
			session.Conn.Write(p.ToBytes())

			p1 := network.EmptyPacket()
			p1.MessageID = opcode.ExchangeStartResponse
			p1.WriteByte(1)
			p1.WriteUInt32(exchange.RequestedUniqueID)
			player.Session.Conn.Write(p1.ToBytes())
		}
	} else {
		log.Warnf("Requested exchange entry not found: %d\n", requestedPlayerUniqueID)
	}
}