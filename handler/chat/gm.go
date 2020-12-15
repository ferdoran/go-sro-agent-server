package chat

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"sync"
)

var adminCommandRegex = regexp.MustCompile("^\\.(?P<command>\\w+)\\s?(?P<args>.*)$")

type GmMessageHandler struct {
}

const (
	SetSpeed        string = "speed"
	CurrentPosition string = "curpos"
	JanganWest      string = "jgwest"
	Takla1          string = "takla1"
	Takla2          string = "takla2"
	LogLevel        string = "loglvl"
)

var gmMessageHandlerInstance *GmMessageHandler
var once sync.Once

func GetGmMessageHandlerInstance() *GmMessageHandler {
	once.Do(func() {
		gmMessageHandlerInstance = &GmMessageHandler{}
	})
	return gmMessageHandlerInstance
}

func (h *GmMessageHandler) HandleAdminMessage(request MessageRequest, session *server.Session) {
	if adminCommandRegex.MatchString(request.Message) {
		strComponents := adminCommandRegex.FindStringSubmatch(request.Message)
		command := strComponents[1]
		args := strComponents[2]
		logrus.Debugf("[GM] - gm command %v with args %v\n", command, args)
		switch command {
		case SetSpeed:
			world := model.GetSroWorldInstance()
			player := world.PlayersByUniqueId[session.UserContext.UniqueID]
			newSpeed, err := strconv.ParseFloat(args, 32)
			if err != nil {
				// TODO send message back
				logrus.Errorf("failed to parse %v as int\n", args)
				return
			}
			if newSpeed < 0 {
				logrus.Warnf("speed cannot be negative")
				return
			}

			player.SetRunSpeed(float32(newSpeed))
			p := network.EmptyPacket()
			p.MessageID = opcode.EntityUpdateMovementSpeed
			p.WriteUInt32(player.UniqueID)
			p.WriteFloat32(player.GetWalkSpeed()) // WalkSpeed
			p.WriteFloat32(player.GetRunSpeed())
			player.Session.Conn.Write(p.ToBytes())

			//spawnEngine := spawn.GetSpawnEngineInstance()
			//spawnEngine.UpdatedMovementSpeed(player)

			logrus.Infof("[GM] - updated %s's movement speed to %f\n", player.CharName, player.GetRunSpeed())
		case LogLevel:
			switch args {
			case "debug":
				logrus.SetLevel(logrus.DebugLevel)
			case "trace":
				logrus.SetLevel(logrus.TraceLevel)
			case "info":
				logrus.SetLevel(logrus.InfoLevel)
			}
		case CurrentPosition:
			world := model.GetSroWorldInstance()
			player := world.PlayersByUniqueId[session.UserContext.UniqueID]
			p := network.EmptyPacket()
			p.MessageID = opcode.ChatUpdate
			p.WriteByte(PM)
			p.WriteString("System")
			p.WriteString(fmt.Sprintf("Current Position: ( %f | %f | %f )", player.Position.X, player.Position.Y, player.Position.Z))
			player.Session.Conn.Write(p.ToBytes())
		case JanganWest:
			warpPlayer(
				session.UserContext.UniqueID,
				435,
				0,
				1745,
				24999)
		case Takla1:
			warpPlayer(
				session.UserContext.UniqueID,
				1374.843750,
				-28.878524,
				937.109375,
				25991)
		case Takla2:
			warpPlayer(
				session.UserContext.UniqueID,
				939.031250,
				-522.488159,
				992.015625,
				26246)
		}
	} else {
		// TODO: Change all players to local region
		players := model.GetSroWorldInstance().PlayersByUniqueId
		for _, v := range players {
			p := network.EmptyPacket()
			p.MessageID = opcode.ChatUpdate
			p.WriteByte(request.ChatType)
			p.WriteUInt32(session.UserContext.UniqueID)
			p.WriteString(request.Message)
			v.Session.Conn.Write(p.ToBytes())
		}

		p1 := network.EmptyPacket()
		p1.MessageID = opcode.ChatResponse
		p1.WriteByte(1)
		p1.WriteByte(request.ChatType)
		p1.WriteByte(request.ChatIndex)
		session.Conn.Write(p1.ToBytes())
	}
}

func warpPlayer(playerUniqueId uint32, x, y, z float32, regionId int16) {
	world := model.GetSroWorldInstance()
	player := world.PlayersByUniqueId[playerUniqueId]

	newPosition := model.Position{
		X:       x,
		Y:       y,
		Z:       z,
		Heading: 0,
		Region:  world.Regions[regionId],
	}
	player.SetPosition(newPosition)
	player.SendPositionUpdate()
}
