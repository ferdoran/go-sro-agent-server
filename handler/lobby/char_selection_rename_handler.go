package lobby

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
)

const (
	OpcodeCharSelectionRenameRequest  uint16 = 0x7450
	OpcodeCharSelectionRenameResponse uint16 = 0xB450

	// Actions
	CharSelectionCharacterRename uint8 = 1
	CharSelectionGuildRename     uint8 = 2
	CharSelectionGuildNameCheck  uint8 = 3

	// Errors
	RENAME_MSG_ERROR_ID                      uint8 = 6
	RENAME_SMERR_NOT_ALLOWED_CHARNAME        uint8 = 7
	RENAME_MSG_GUILDERR_SAME_GUILDNAME_EXIST uint8 = 6
	RENAME_MSG_GUILD_NOT_CREATE              uint8 = 7
)

type CharSelectionRenameRequestHandler struct {
	channel chan server.PacketChannelData
}

func InitCharSelectionRenameRequestHandler() {
	handler := CharSelectionRenameRequestHandler{channel: server.PacketManagerInstance.GetQueue(opcode.LobbyRenameRequest)}
	go handler.Handle()
}

func (h *CharSelectionRenameRequestHandler) Handle() {
	for {
		packet := <-h.channel
		log.Println("AGENT_CHARACTER_SELECTION_RENAME")

		var result byte
		action, err := packet.ReadByte()
		if err != nil {
			log.Panicln("Failed to read action")
		}

		if action == CharSelectionCharacterRename || action == CharSelectionGuildRename {
			currentGuildName, err := packet.ReadString()
			if err != nil {
				log.Panicln("Failed to read currentGuildName")
			}

			newGuildName, err := packet.ReadString()
			if err != nil {
				log.Panicln("Failed to read newGuildName")
			}

			log.Println(currentGuildName)
			log.Println(newGuildName)

			// TODO: Change the char / guild name
		} else if action == CharSelectionGuildNameCheck {
			guildName, err := packet.ReadString()
			if err != nil {
				log.Panicln("Failed to read guildName")
			}

			log.Println(guildName)
			// TODO: Check if guild name exists
		}

		// TODO: Remove this after stuff above is implemented
		result = ResultTrue

		p := network.EmptyPacket()
		p.MessageID = OpcodeCharSelectionRenameResponse

		p.WriteByte(result)

		packet.Session.Conn.Write(p.ToBytes())
	}
}
