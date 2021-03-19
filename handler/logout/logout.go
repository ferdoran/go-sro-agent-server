package logout

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/sirupsen/logrus"
	"time"
)

type LogoutHandler struct {
	logoutRequestChannel chan server.PacketChannelData
	logoutCancelChannel  chan server.PacketChannelData
}

func InitLogoutHandler() {
	handler := LogoutHandler{
		logoutRequestChannel: server.PacketManagerInstance.GetQueue(opcode.AgentLogoutRequest),
		logoutCancelChannel:  server.PacketManagerInstance.GetQueue(opcode.AgentLogoutCancelRequest),
	}
	go handler.Handle()
}

const (
	LogoutModeExit byte = iota + 1
	LogoutModeRestart
	LogoutErrorInBattleState uint16 = iota + 0x801
	LogoutErrorInTeleportState
	CountdownTime byte = 5
)

func (h LogoutHandler) Handle() {
	for {
		select {
		case data := <-h.logoutRequestChannel:
			doLogout(data)
		case data := <-h.logoutCancelChannel:
			doCancelLogout(data)
		}
	}
}

func doLogout(data server.PacketChannelData) {
	logoutMode, err := data.ReadByte()
	if err != nil {
		logrus.Panicln("Failed to read logoutMode")
	}

	logrus.Debug("Logging out char")
	p := network.EmptyPacket()
	p.MessageID = opcode.AgentLogoutResponse
	// TODO: Evaluate the result based on current Player state
	var result byte = 1
	p.WriteByte(result)
	if result == 1 {
		p.WriteByte(CountdownTime) // Countdown in seconds
		p.WriteByte(logoutMode)
	} else if result == 2 {
		// TODO: Determine Error code
		p.WriteUInt16(LogoutErrorInBattleState)
	}
	data.Conn.Write(p.ToBytes())

	time.AfterFunc(time.Second*time.Duration(CountdownTime), func() {
		p1 := network.EmptyPacket()
		p1.MessageID = opcode.AgentLogoutSuccess
		data.Conn.Write(p1.ToBytes())
		// FIXME
		//spawnEngine := spawn.GetSpawnEngineInstance()
		//spawnEngine.PlayerDisconnected(data.UserContext.UniqueID, data.UserContext.CharName)
	})
}

func doCancelLogout(data server.PacketChannelData) {
	p := network.EmptyPacket()
	p.MessageID = opcode.AgentLogoutCancelResponse
	// TODO: Evaluate the result based on current Player state
	var result byte = 1
	p.WriteByte(result)
	if result == 2 {
		// TODO: Determine Error code
		p.WriteUInt16(LogoutErrorInBattleState)
	}
	data.Conn.Write(p.ToBytes())
}
