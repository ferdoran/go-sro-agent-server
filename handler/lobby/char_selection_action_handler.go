package lobby

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	DeletionTimeInDays                       = 7
	OpcodeCharSelectionActionRequest  uint16 = 0x7007
	OpcodeCharSelectionActionResponse uint16 = 0xB007

	// Actions
	CharActionCreate    uint8 = 1
	CharActionList      uint8 = 2
	CharActionDelete    uint8 = 3
	CharActionCheckName uint8 = 4
	CharActionRestore   uint8 = 5

	// Errors
	ACTION_MSG_ERROR_CHARACTER_OVER_3       uint16 = 0x405
	ACTION_MSG_ERROR_CHARACTER_SELECTWEAPON uint16 = 0x404
	ACTION_MSG_ERROR_CHARACTER_NAME_STRING  uint16 = 0x40C
	ACTION_SMERR_NOT_ALLOWED_CHARNAME       uint16 = 0x40D
	ACTION_MSG_ERROR_ID                     uint16 = 0x410
	ACTION_MSG_ERROR_OVERLAP                uint16 = 0x411
	ACTION_SMERR_MAX_USER_EXCEEDED          uint16 = 0x414
	ACTION_SMERR_INVALID_CHARGEN_INFO       uint16 = 0x403
	ACTION_SMERR_FAILED_TO_CREATE_CHARACTER uint16 = 0x406
	ACTION_SMERR_CANT_FIND_GAMESERVER       uint16 = 0x409
	ACTION_SMERR_CANT_ACCESS_PARENT_SERVER  uint16 = 0x40F
	ACTION_SMERR_FAILED_TO_CREATE_NEW_USER  uint16 = 0x412
	ACTION_SMERR_FAILED_TO_ENTERLOBBY       uint16 = 0x415

	// Memberclass
	Member uint8 = 1
	Master uint8 = 2

	ResultTrue  uint8 = 0x01
	ResultFalse uint8 = 0x02
)

type CharSelectionActionRequestHandler struct {
}

func NewCharSelectionActionRequestHandler() server.PacketHandler {
	handler := CharSelectionActionRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.LobbyActionRequest, handler)
	return handler
}

func (h CharSelectionActionRequestHandler) Handle(packet server.PacketChannelData) {
	log.Println("AGENT_CHARACTER_SELECTION_ACTION")

	var result byte

	// TODO: Remove this after stuff above is implemented
	result = ResultTrue

	action, err := packet.ReadByte()
	if err != nil {
		log.Panicln("Failed to read action")
	}

	p := network.EmptyPacket()
	p.MessageID = OpcodeCharSelectionActionResponse

	if action == CharActionCreate {
		username, err := packet.ReadString()
		if err != nil {
			log.Panicln("Failed to read username")
		}

		charModel, err := packet.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read charModel")
		}

		charScale, err := packet.ReadByte()
		if err != nil {
			log.Panicln("Failed to read scale")
		}

		equippedChest, err := packet.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read equippedChest")
		}

		equippedPants, err := packet.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read equippedPants")
		}

		equippedBoots, err := packet.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read equippedBoots")
		}

		equippedWeapon, err := packet.ReadUInt32()
		if err != nil {
			log.Panicln("Failed to read equippedWeapon")
		}

		// TODO: Spawn coords for JG
		char := model.Char{
			RefObjID:   int(charModel),
			User:       int(packet.Session.UserContext.UserID),
			Shard:      int(packet.Session.UserContext.ShardID),
			Name:       username,
			Scale:      charScale,
			Level:      1,
			Exp:        0,
			SkillExp:   0,
			Str:        20,
			Int:        20,
			StatPoints: 0,
			HP:         200,
			MP:         200,
			PosX:       950,
			PosY:       40,
			PosZ:       1091,
			Region:     0x62A8,
			IsDeleting: false,
		}

		if res, _ := model.CreateChar(char, equippedWeapon, equippedChest, equippedBoots, equippedPants); res {
			result = ResultTrue
		} else {
			result = ResultFalse
		}

		p.WriteByte(action)
		p.WriteByte(result)

	} else if action == CharActionDelete || action == CharActionCheckName || action == CharActionRestore {
		username, err := packet.ReadString()
		if err != nil {
			log.Panicln("Failed to read username")
		}

		if action == CharActionDelete {
			// TODO: Do actually delete after 3 days or something.
			if res := model.MarkCharIsDeletion(1, username); res {
				result = ResultTrue
			} else {
				result = ResultFalse
			}
		} else if action == CharActionCheckName {
			if res := model.DoesCharNameExist(username); res {
				result = ResultFalse
				p.WriteUInt16(0x40D)
			} else {
				result = ResultTrue
			}
		} else if action == CharActionRestore {
			if res := model.MarkCharIsDeletion(0, username); res {
				result = ResultTrue
			} else {
				result = ResultFalse
			}
		}
	}

	if result == ResultTrue && action == CharActionList {
		p.WriteByte(action)
		p.WriteByte(result)

		// TODO: Get real user id
		chars := model.GetCharactersByUserId(int(packet.Session.UserContext.UserID))
		p.WriteByte(byte(len(chars)))

		for _, char := range chars {
			// Char data
			p.WriteUInt32(uint32(char.RefObjID))
			p.WriteString(char.Name)
			p.WriteByte(byte(char.Scale))
			p.WriteByte(byte(char.Level))
			p.WriteUInt64(uint64(char.Exp))
			p.WriteUInt16(uint16(char.Str))
			p.WriteUInt16(uint16(char.Int))
			p.WriteUInt16(uint16(char.StatPoints))
			p.WriteUInt32(uint32(char.HP))
			p.WriteUInt32(uint32(char.MP))
			// Is deleting
			if char.IsDeleting {
				p.WriteByte(ResultTrue)
				now := time.Now().UTC()
				diff := char.Utime.Add(DeletionTimeInDays * 24 * time.Hour).Sub(now)
				p.WriteUInt32(uint32(int(diff.Minutes())))
			} else {
				p.WriteByte(0)
			}
			// Guild
			// TODO: Add a real check
			// 	0 - No Guild
			// 	1 - Member
			// 	2 - Master
			p.WriteByte(0)
			p.WriteByte(0)
			// Academy
			// TODO: Add a real check
			// 	0 - No Academy
			// 	1 - Member
			// 	2 - Master
			p.WriteByte(0)
			// Items
			inv := model.GetCharacterInventory(int64(char.ID))
			p.WriteByte(byte(len(inv.Items)))
			for _, item := range inv.Items {
				// TODO: Only send certain items (armor + weapon)
				p.WriteUInt32(item.GetRefObjectID())
				// TODO: Add a real check for item + value
				p.WriteByte(0) // Plus

			}

			// Avatar
			// TODO: Add a real check
			p.WriteByte(0)
		}
		// TODO: Get the chars
		/*
			1   byte    characterCount
			foreach(character)
			{
				4   uint    character.RefObjID
				2   ushort  character.Name.Length
				*   string  character.Name
				1   byte    character.Scale
				1   byte    character.CurLevel
				8   ulong   character.ExpOffset
				2   ushort  character.Strength
				2   ushort  character.Intelligence
				2   ushort  character.StatPoint
				4   uint    character.CurHP
				4   uint    character.CurMP
				1   bool    isDeleting
				if(isDeleting)
				{
					4   uint    character.DeleteTime	//in Minutes
				}

				1   byte    guildMemberClass
				1   bool    isGuildRenameRequired
				if(isGuildRenameRequired)
				{
					2   ushort  CurGuildName.Length
					*   string  CurGuildName
				}
				1   byte    academyMemberClass

				1   byte    itemCount
				foreach(item)
				{
					4   uint    item.RefItemID
					1   byte    item.Plus
				}

				1   byte    avatarItemCount
				foreach(avatarItem)
				{
					4   uint    avatarItem.RefItemID
					1   byte    avatarItem.Plus
				}
			}
		*/
	} else {
		p.WriteByte(action)
		p.WriteByte(result)
	}

	packet.Session.Conn.Write(p.ToBytes())
}
