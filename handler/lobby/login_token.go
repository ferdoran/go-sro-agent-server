package lobby

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

type CreateLoginTokenHandler struct {
	Session *server.Session
	Tokens  map[string]LoginTokenData
}

type LoginTokenData struct {
	AccountID uint32
	Username  string
	Password  string
	Token     uint32
	ShardID   uint16
}

func (h *CreateLoginTokenHandler) Handle(packet network.Packet) {
	accountId, err := packet.ReadUInt32()
	if err != nil {
		log.Panic("Failed to read account id")
	}
	username, err := packet.ReadString()
	if err != nil {
		log.Panic("Failed to read username")
	}
	password, err := packet.ReadString()
	if err != nil {
		log.Panic("Failed to read password")
	}
	token, err := packet.ReadUInt32()
	if err != nil {
		log.Panic("Failed to read token")
	}
	shardId, err := packet.ReadUInt16()
	if err != nil {
		log.Panic("Failed to read shard id")
	}

	// TODO Remove pw from log entry
	log.Debugf("login token request (%v, %v, %v, %v, %v)\n", accountId, username, password, token, shardId)
	data := LoginTokenData{AccountID: accountId, Username: username, Password: password, Token: token, ShardID: shardId}

	h.Tokens[username] = data
}
