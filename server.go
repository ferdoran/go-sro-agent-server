package main

import (
	"bufio"
	"github.com/ferdoran/go-sro-agent-server/config"
	"github.com/ferdoran/go-sro-agent-server/handler/character"
	"github.com/ferdoran/go-sro-agent-server/handler/chat"
	"github.com/ferdoran/go-sro-agent-server/handler/inventory"
	"github.com/ferdoran/go-sro-agent-server/handler/lobby"
	"github.com/ferdoran/go-sro-agent-server/handler/logout"
	"github.com/ferdoran/go-sro-agent-server/handler/party"
	"github.com/ferdoran/go-sro-agent-server/handler/stall"
	"github.com/ferdoran/go-sro-agent-server/manager"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/logging"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	gwHandlers "github.com/ferdoran/go-sro-gateway-server/handler"
	log "github.com/sirupsen/logrus"
	"os"
)

type AgentServer struct {
	*server.Server
	Config                  config.AgentConfig
	failedLogins            map[string]int
	GatewaySession          *server.Session
	Tokens                  map[string]lobby.LoginTokenData
	UnhandledPacketsLogger  *log.Logger
	CreateLoginTokenHandler *lobby.CreateLoginTokenHandler
	//SpawnEngine             *spawn.SpawnEngine
}

func NewAgentServer(config config.AgentConfig) AgentServer {
	server := server.NewEngine(
		config.AgentServer.IP,
		config.AgentServer.Port,
		network.EncodingOptions{
			None:         false,
			Disabled:     false,
			Encryption:   true,
			EDC:          true,
			KeyExchange:  true,
			KeyChallenge: false,
		},
		config.Config,
	)

	server.ModuleID = config.AgentServer.ModuleID
	return AgentServer{
		Server:                 &server,
		Config:                 config,
		failedLogins:           make(map[string]int),
		Tokens:                 make(map[string]lobby.LoginTokenData),
		UnhandledPacketsLogger: logging.UnhandledPacketLogger(),
		//SpawnEngine:            spawn.GetSpawnEngineInstance(),
	}
}

func (a *AgentServer) Start() {
	go a.Server.Start()
	a.handlePackets()
}

func (a *AgentServer) handlePackets() {

	gwHandlers.NewPatchRequestHandler()
	gwHandlers.NewShardlistRequestHandler()
	lobby.NewAuthRequestHandler(a.Tokens)
	lobby.NewCharSelectionJoinRequestHandler()
	lobby.NewCharSelectionActionRequestHandler()
	lobby.NewCharSelectionRenameRequestHandler()
	character.NewGuideHandler()
	character.NewMovementHandler()
	chat.NewChatHandler()
	logout.NewLogoutHandler()
	stall.NewStallCreateHandler()
	stall.NewStallDestroyHandler()
	stall.NewStallLeaveHandler()
	stall.NewStallUpdateHandler()
	inventory.NewInventoryHandler()
	party.NewPartyMatchingFormHandler()
	party.NewPartyMatchingUpdateHandler()
	party.NewPartyMatchingDeleteHandler()
	party.NewPartyMatchingListHandler()
	party.NewPartyMatchingJoinRequestHandler()
	party.NewPartyMatchingPlayerJoinRequestHandler()
	party.NewPartyKickHandler()

	for {
		select {
		case closedSession := <-a.SessionClosed:
			if closedSession == a.GatewaySession {
				log.Println("Gateway connection closed")
				a.GatewaySession = nil
				continue
			} else if closedSession.UserContext.UniqueID != 0 {
				log.Debugf("player %s disconnected", closedSession.UserContext.CharName)
				//a.SpawnEngine.PlayerDisconnected(closedSession.UserContext.UniqueID, closedSession.UserContext.CharName)
				// TODO check what player disconnected and reset his auth tokens

				delete(a.Tokens, closedSession.UserContext.Username)
				//delete(a.AuthRequestHandler.Tokens, closedSession.UserContext.Username)
				//delete(a.CreateLoginTokenHandler.Tokens, closedSession.UserContext.Username)
				delete(a.Server.Sessions, closedSession.ID)
				world := model.GetSroWorldInstance()
				world.PlayerDisconnected(closedSession.UserContext.UniqueID, closedSession.UserContext.CharName)
			}

		case connectedBackend := <-a.BackendConnection:
			a.serverModuleConnected(connectedBackend)
		case data := <-a.Server.PacketChannel:
			if a.GatewaySession != nil && data.Session != nil && a.GatewaySession.ID == data.Session.ID {
				a.handleGatewayPackets(data)
				continue
			}

			switch data.MessageID {
			//case opcode.AuthRequest:
			//	handler := &lobby.AuthRequestHandler{Tokens: a.Tokens}
			//	handler.Handle(data)
			case opcode.StallTalkRequest:
				log.Debugf("Stall Talk Request not handled")

			default:
				a.UnhandledPacketsLogger.Printf("Unhandled packet %+v\n", data.Packet)
			}
		}
	}
}

func (a *AgentServer) handleGatewayPackets(data server.PacketChannelData) {

	switch data.MessageID {
	case opcode.GatewayLoginTokenRequest:
		handler := &lobby.CreateLoginTokenHandler{
			Session: a.GatewaySession,
			Tokens:  a.Tokens,
		}
		a.CreateLoginTokenHandler = handler
		handler.Handle(data.Packet)
	default:
		log.Printf("Unhandled gateway packet %+v\n", data.Packet)
	}
}

func (a *AgentServer) serverModuleConnected(data server.BackendConnectionData) {
	switch data.ModuleID {
	case a.Config.GatewayServer.ModuleID:
		a.GatewaySession = data.Session
		a.Server.Sessions[data.Session.ID] = nil
	}
}

func main() {
	logging.Init()
	reader := bufio.NewReader(os.Stdin)

	config.LoadConfig("config.json")
	log.Println("Starting server...")
	gw := NewAgentServer(config.GlobalConfig)

	for k, v := range model.RefItems {
		model.RefObjects[k] = v.RefObject
	}
	for k, v := range model.RefChars {
		model.RefObjects[k] = v.RefObject
	}

	world := model.InitSroWorldInstance(config.GlobalConfig.AgentServer.DataPath, config.GlobalConfig.AgentServer.NavmeshGOB)
	world.LoadGameServerRegions(1)

	model.RefItems = model.GetAllRefItems()
	model.RefChars = model.GetAllRefChars()

	manager.GetGameTimeManagerInstance().Start()

	klm := manager.GetKnownListManager()
	klm.Start()

	sm := manager.GetSpawnManagerInstance()
	sm.Start()

	model.GetSroWorldInstance().InitiallySpawnAllNpcs()

	gw.Start()
	log.Println("Press Enter to exit...")
	reader.ReadString('\n')
}
