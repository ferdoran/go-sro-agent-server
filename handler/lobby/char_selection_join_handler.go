package lobby

import (
	"github.com/ferdoran/go-sro-agent-server/engine/environment"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-framework/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

const (
	OpcodeCharSelectionJoinRequest  uint16 = 0x7001
	OpcodeCharSelectionJoinResponse uint16 = 0xB001

	// Errors
	JOIN_ERROR_CHARACTER_OVER_3           uint16 = 0x405
	JOIN_ERROR_CHARACTER_SELECTWEAPON     uint16 = 0x404
	JOIN_ERROR_CHARACTER_NAME_STRING      uint16 = 0x40C
	JOIN_SMERR_NOT_ALLOWED_CHARNAME       uint16 = 0x40D
	JOIN_MSG_ERROR_ID                     uint16 = 0x410
	JOIN_MSG_ERROR_OVERLAP                uint16 = 0x411
	JOIN_SMERR_MAX_USER_EXCEEDED          uint16 = 0x414
	JOIN_SMERR_INVALID_CHARGEN_INFO       uint16 = 0x403
	JOIN_SMERR_FAILED_TO_CREATE_CHARACTER uint16 = 0x406
	JOIN_SMERR_CANT_FIND_GAMESERVER       uint16 = 0x409
	JOIN_SMERR_CANT_ACCESS_PARENT_SERVER  uint16 = 0x40F
	JOIN_SMERR_FAILED_TO_CREATE_NEW_USER  uint16 = 0x412
	JOIN_SMERR_FAILED_TO_ENTERLOBBY       uint16 = 0x415

	//UIO_MSG_ERROR_SEVER_CONNECT = 402,407,408,40A,40B,40E,413,416,417,418,
)

type CharSelectionJoinRequestHandler struct {
}

func NewCharSelectionJoinRequestHandler() server.PacketHandler {
	handler := CharSelectionJoinRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.JoinLobbyRequest, handler)
	return handler
}

func (h CharSelectionJoinRequestHandler) Handle(data server.PacketChannelData) {
	log.Println("AGENT_CHARACTER_SELECTION_JOIN")

	charName, err := data.ReadString()
	if err != nil {
		log.Panicln("Failed to read charName")
	}

	log.Println(charName)
	// TODO: Check for any errors with error codes above
	var result byte

	// TODO: Remove this after stuff above is implemented
	result = ResultTrue

	// TODO: Check if player can join
	pJoinResponse := network.EmptyPacket()
	pJoinResponse.MessageID = OpcodeCharSelectionJoinResponse

	pJoinResponse.WriteByte(result)
	data.Conn.Write(pJoinResponse.ToBytes())

	player := h.LoadPlayerData(charName, data.Session)
	world := model.GetSroWorldInstance()
	world.AddPlayer(player)
	player.LifeState = model.Spawning
	player.BodyState = model.NoStatus
	player.MotionState = model.NoMotion
	data.UserContext.UniqueID = player.UniqueID
	data.UserContext.CharName = player.CharName
	data.UserContext.RefObjId = player.GetRefObjectID()

	// SERVER_AGENT_CHARACTER_DATA_BEGIN
	log.Debugf("Setting up player: %+v\n", player)
	pDataBegin := network.EmptyPacket()
	pDataBegin.MessageID = opcode.CharacterDataBegin

	data.Conn.Write(pDataBegin.ToBytes())

	pDataBody := network.EmptyPacket()
	pDataBody.MessageID = opcode.CharacterDataBody
	pDataBody.WriteUInt32(utils.ToSilkroadTime(time.Now()))
	// 1. Get Char Data
	WriteCharDataToPacket(&pDataBody, player)

	// 2. Get Inventory
	inv := model.GetCharacterInventory(int64(player.ID))
	WriteInventoryToPacket(&pDataBody, inv)
	// 3. Avatar Inventory
	pDataBody.WriteByte(5)
	pDataBody.WriteByte(0)
	// 4. Masteries
	pDataBody.WriteByte(0) // mastery begin

	WriteMasteryOrSkill(&pDataBody, 0x00000101, 0)
	WriteMasteryOrSkill(&pDataBody, 0x00000102, 0)
	WriteMasteryOrSkill(&pDataBody, 0x00000103, 0)
	WriteMasteryOrSkill(&pDataBody, 0x00000111, 0)
	WriteMasteryOrSkill(&pDataBody, 0x00000112, 0)
	WriteMasteryOrSkill(&pDataBody, 0x00000113, 0)
	WriteMasteryOrSkill(&pDataBody, 0x00000114, 0)
	pDataBody.WriteByte(2) // Set next master to 2 (2 seems to tell that the masteries section finished)
	pDataBody.WriteByte(0) // unk Byte
	// 5. Skills
	pDataBody.WriteByte(2) // nextSkill
	// 6. Quests
	pDataBody.WriteUInt16(1) // Completed Quest Count
	pDataBody.WriteUInt32(1) // Completed Quest
	//pDataBody.WriteUInt32(0)
	pDataBody.WriteByte(0) // Active Quest Count
	// 7. Collection Book
	pDataBody.WriteByte(0)   // unk Byte
	pDataBody.WriteUInt32(0) // CollectionBookStartedThemeCount

	pDataBody.WriteUInt32(player.GetUniqueID()) // UniqueID
	// 8. Position
	model.WritePosition(&pDataBody, player.GetPosition())
	// 9. Movement
	pDataBody.WriteByte(0)                                 // HasDestination
	pDataBody.WriteByte(1)                                 // Type
	pDataBody.WriteByte(0)                                 // Source
	pDataBody.WriteUInt16(uint16(player.Position.Heading)) // Angle
	// 10. State
	pDataBody.WriteByte(byte(model.Spawning))     // LifeState
	pDataBody.WriteByte(0)                        // unkByte
	pDataBody.WriteByte(byte(model.NoMotion))     // MotionState
	pDataBody.WriteByte(byte(player.BodyState))   // Status
	pDataBody.WriteFloat32(player.GetWalkSpeed()) // WalkSpeed
	pDataBody.WriteFloat32(player.GetRunSpeed())  // RunSpeed
	pDataBody.WriteFloat32(player.GetHwanSpeed()) // HwanSpeed
	pDataBody.WriteByte(0)                        // BuffCount
	pDataBody.WriteString(player.CharName)
	pDataBody.WriteUInt16(0)                 // JobName.Length
	pDataBody.WriteByte(0)                   // JobType
	pDataBody.WriteByte(1)                   // JobLevel
	pDataBody.WriteUInt32(0)                 // JobExp
	pDataBody.WriteUInt32(0)                 // JobContribution
	pDataBody.WriteUInt32(0)                 // JobReward
	pDataBody.WriteByte(0)                   // PVP State
	pDataBody.WriteByte(0)                   // TransportFlag
	pDataBody.WriteByte(0)                   // InCombat
	pDataBody.WriteByte(0xFF)                // PVP Flag
	pDataBody.WriteUInt64(0)                 // GuideFlag
	pDataBody.WriteUInt32(uint32(player.ID)) // JID
	if player.IsGm {
		pDataBody.WriteByte(1) // GMFlag
	} else {
		pDataBody.WriteByte(0) // GMFlag
	}
	pDataBody.WriteByte(0)            // ActivationFlag
	pDataBody.WriteByte(0)            // Hotkeys.Count
	pDataBody.WriteUInt16(0)          // AutoHPConfig
	pDataBody.WriteUInt16(0)          // AutoMPConfig
	pDataBody.WriteUInt16(0)          // AutoUniversalConfig
	pDataBody.WriteByte(0)            // AutoPotionDelay
	pDataBody.WriteByte(0)            // BlockedWhisperCount
	pDataBody.WriteUInt32(0x00010001) // unkUShort
	pDataBody.WriteByte(0)            // unkByte

	data.Conn.Write(pDataBody.ToBytes())

	// SERVER_AGENT_CHARACTER_DATA_END
	pDataEnd := network.EmptyPacket()
	pDataEnd.MessageID = opcode.CharacterDataEnd
	data.Conn.Write(pDataEnd.ToBytes())

	// Environment CELESTIAL
	celestialPosition := environment.GetCelestialPositionForPlayer(player)
	pCelestialPos := network.EmptyPacket()
	pCelestialPos.MessageID = opcode.CelestialPosition
	pCelestialPos.WriteUInt32(celestialPosition.CharUniqueID)
	pCelestialPos.WriteUInt16(celestialPosition.Moonphase)
	pCelestialPos.WriteByte(celestialPosition.Hour)
	pCelestialPos.WriteByte(celestialPosition.Minute)
	data.Conn.Write(pCelestialPos.ToBytes())

	// Environment WEATHER
	weatherType, intensity := environment.GetCurrentWeather()
	pWeather := network.EmptyPacket()
	pWeather.MessageID = opcode.WeatherUpdate
	pWeather.WriteByte(byte(weatherType))
	pWeather.WriteByte(intensity)
	data.Conn.Write(pWeather.ToBytes())

	player.SendStatsUpdate()
	player.SendMovementStateUpdate()

	pDataBody = network.EmptyPacket()
	pDataBody.MessageID = 0x3077
	pDataBody.WriteByte(0)
	pDataBody.WriteByte(0)
	data.Conn.Write(pDataBody.ToBytes())
	//log.Debugf(pDataBody.String())

	player.LifeState = model.Alive

	//spawnEngine := spawn.GetSpawnEngineInstance()
	//spawnEngine.PositionChanged(&player)
	log.Infof("Send Agent spawn data")
}

func (h *CharSelectionJoinRequestHandler) LoadPlayerData(charName string, session *server.Session) *model.Player {
	char := model.GetCharacterByName(charName)
	world := model.GetSroWorldInstance()
	angle := rand.Int() % 0xFFFF
	player := &model.Player{
		Session:   session,
		ID:        char.ID,
		CharName:  charName,
		Scale:     char.Scale,
		Inventory: h.LoadPlayerInventory(int64(char.ID)),
		BaseStats: model.BaseStats{
			HP:  char.HP,
			MP:  char.MP,
			Str: char.Str,
			Int: char.Int,
		},
		// TODO calculate below values
		//BaseAttackStats:   model.AttackStats{},
		//BaseDefenseStats:  model.DefenseStats{},
		//BonusStats:        model.BonusBaseStats{},
		//BonusAttackStats:  model.BonusAttackStats{},
		//BonusDefenseStats: model.BonusDefenseStats{},
		//PhyAbsorbPercent:  0,
		//MagAbsorbPercent:  0,
		//PhyBalancePercent: 0,
		//MagBalancePercent: 0,
		SkillPoints:    int(char.SkillPoints),
		StatPoints:     char.StatPoints,
		ExpOffset:      uint64(char.Exp),
		SkillExpOffset: uint(char.SkillExp),
		Level:          char.Level,
		MaxLevel:       char.MaxLevel,
		Mutex:          sync.Mutex{},
		IsGm:           char.IsGm,
	}
	player.Type = model.PlayerType
	player.WalkSpeed = 16
	player.RunSpeed = 50
	player.HwanSpeed = 100
	player.RefObjectID = uint32(char.RefObjID)
	player.UniqueID = 0

	region, err := world.GetRegion(char.Region)

	if err != nil {
		log.Panic(err)
	}
	player.Position = model.Position{
		X:       char.PosX,
		Y:       char.PosY,
		Z:       char.PosZ,
		Heading: float32(angle),
		Region:  region,
	}
	player.Name = player.CharName
	player.TypeInfo = model.RefChars[player.GetRefObjectID()].TypeInfo
	player.CharKnownObjectList = model.NewCharKnownObjectList(player)
	player.KnownObjectList = player.CharKnownObjectList

	return player
}

func (h *CharSelectionJoinRequestHandler) LoadPlayerInventory(charId int64) model.Inventory {
	inventoryDao := model.GetCharacterInventory(charId)
	inventory := model.Inventory{
		Items: make(map[byte]model.Item),
	}

	for slot, item := range inventoryDao.Items {
		inventory.Items[slot] = item
	}

	return inventory
}

func WriteCharDataToPacket(p *network.Packet, data *model.Player) {
	p.WriteUInt32(data.GetRefObjectID())
	p.WriteByte(byte(data.Scale))
	p.WriteByte(byte(data.Level))
	p.WriteByte(byte(data.MaxLevel))
	p.WriteUInt64(data.ExpOffset)
	p.WriteUInt32(uint32(data.SkillExpOffset))
	p.WriteUInt64(1_000_000_000)           // Gold
	p.WriteUInt32(100_000)                 // Skill Points
	p.WriteUInt16(uint16(data.StatPoints)) // Stat Points
	p.WriteByte(0)                         // Berserker Points
	p.WriteUInt32(0)
	p.WriteUInt32(uint32(data.BaseStats.HP))
	p.WriteUInt32(uint32(data.BaseStats.MP))
	p.WriteByte(1)   // Auto Inverst Exp
	p.WriteByte(0)   // Daily PK
	p.WriteUInt16(0) // Total PK
	p.WriteUInt32(0) // PK Penalty Point
	p.WriteByte(0)   // Berserk Level
	p.WriteByte(0)   // Free PVP
}

func WriteInventoryToPacket(p *network.Packet, data model.Inventory) {
	p.WriteByte(45) // Inventory Size
	//p.WriteByte(0)  // Item Count
	//return
	p.WriteByte(byte(len(data.Items))) // Item Count

	for slot, item := range data.Items {
		p.WriteByte(slot)

		model.WriteRentInfo(p, item)

		p.WriteUInt32(item.GetRefObjectID())
		model.WriteInventoryItem(p, item)
	}
}

func WriteMasteryOrSkill(p *network.Packet, masteryId uint32, masteryLvl byte) {
	p.WriteByte(1)
	p.WriteUInt32(masteryId)
	p.WriteByte(masteryLvl)
}
