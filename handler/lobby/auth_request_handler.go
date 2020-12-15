package lobby

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

const (
	MacAddressLength          = 6
	OpcodeAuthResponse uint16 = 0xA103
)

type AuthRequestHandler struct {
	Tokens map[string]LoginTokenData
}

func NewAuthRequestHandler(tokens map[string]LoginTokenData) server.PacketHandler {
	handler := &AuthRequestHandler{Tokens: tokens}
	server.PacketManagerInstance.RegisterHandler(opcode.AuthRequest, handler)
	return handler
}

func (h *AuthRequestHandler) Handle(packet server.PacketChannelData) {
	log.Println("Agent auth")
	token, err := packet.ReadUInt32()
	if err != nil {
		log.Panicln("Failed to read token")
	}
	username, err := packet.ReadString()
	if err != nil {
		log.Panicln("Failed to read username")
	}
	password, err := packet.ReadString()
	if err != nil {
		log.Panicln("Failed to read password")
	}
	// TODO Use content id?
	_, err = packet.ReadByte()
	if err != nil {
		log.Panicln("Failed to read content id")
	}
	// TODO Use Mac Address?
	_, err = packet.ReadBytes(MacAddressLength)
	if err != nil {
		log.Panicln("Failed to read mac address")
	}

	loginData := h.Tokens[username]

	p := network.EmptyPacket()
	p.MessageID = OpcodeAuthResponse
	if loginData.Username != username || loginData.Password != password || loginData.Token != token {
		// Invalid login
		p.WriteByte(ResultFalse)

		// TODO Scenarios for other error codes
		//  3 = Agent Server not in service
		//  4 = Server is full
		//  5 = IP Limit

		p.WriteByte(4)
		packet.Session.Conn.Write(p.ToBytes())
	} else {
		packet.Session.UserContext = server.UserContext{
			UserID:   loginData.AccountID,
			ShardID:  loginData.ShardID,
			Username: loginData.Username,
		}
		p.WriteByte(ResultTrue)
	}
	delete(h.Tokens, username)
	packet.Session.Conn.Write(p.ToBytes())

}
