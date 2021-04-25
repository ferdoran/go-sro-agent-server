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
	"github.com/ferdoran/go-sro-agent-server/handler/entity_select"
	"github.com/ferdoran/go-sro-agent-server/manager"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/boot"
	"github.com/ferdoran/go-sro-framework/logging"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	gwHandlers "github.com/ferdoran/go-sro-gateway-server/handler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type AgentServer struct {
	*server.Server
	failedLogins            map[string]int
	GatewaySession          *server.Session
	Tokens                  map[string]lobby.LoginTokenData
	UnhandledPacketsLogger  *log.Logger
	CreateLoginTokenHandler *lobby.CreateLoginTokenHandler
	backendModules          map[string]string
	//SpawnEngine             *spawn.SpawnEngine
}

func NewAgentServer() AgentServer {
	serv := server.NewEngine(
		viper.GetString(config.AgentHost),
		viper.GetInt(config.AgentPort),
		network.EncodingOptions{
			None:         false,
			Disabled:     false,
			Encryption:   true,
			EDC:          true,
			KeyExchange:  true,
			KeyChallenge: false,
		},
	)

	serv.ModuleID = viper.GetString(config.AgentModuleId)
	backendModules := make(map[string]string)
	backendModules[viper.GetString(config.GatewayModuleId)] = viper.GetString(config.GatewaySecret)
	agentServer := AgentServer{
		Server:                 &serv,
		failedLogins:           make(map[string]int),
		Tokens:                 make(map[string]lobby.LoginTokenData),
		UnhandledPacketsLogger: logging.UnhandledPacketLogger(),
		backendModules:         backendModules,
	}

	server.InitBackendConnectionHandler(agentServer.BackendConnection, agentServer.backendModules)
	lobby.InitAuthRequestHandler(agentServer.Tokens)
	boot.RegisterComponent("packethandler", gwHandlers.InitPatchRequestHandler, 2)
	boot.RegisterComponent("packethandler", gwHandlers.InitShardlistRequestHandler, 2)
	boot.RegisterComponent("packethandler", lobby.InitCharSelectionJoinRequestHandler, 2)
	boot.RegisterComponent("packethandler", lobby.InitCharSelectionActionRequestHandler, 2)
	boot.RegisterComponent("packethandler", lobby.InitCharSelectionRenameRequestHandler, 2)
	boot.RegisterComponent("packethandler", character.InitGuideHandler, 2)
	boot.RegisterComponent("packethandler", character.InitMovementHandler, 2)
	boot.RegisterComponent("packethandler", chat.InitChatHandler, 2)
	boot.RegisterComponent("packethandler", logout.InitLogoutHandler, 2)
	boot.RegisterComponent("packethandler", stall.InitStallCreateHandler, 2)
	boot.RegisterComponent("packethandler", stall.InitStallDestroyHandler, 2)
	boot.RegisterComponent("packethandler", stall.InitStallLeaveHandler, 2)
	boot.RegisterComponent("packethandler", stall.InitStallUpdateHandler, 2)
	boot.RegisterComponent("packethandler", inventory.InitInventoryHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyMatchingFormHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyMatchingUpdateHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyMatchingDeleteHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyMatchingListHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyMatchingJoinRequestHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyMatchingPlayerJoinRequestHandler, 2)
	boot.RegisterComponent("packethandler", party.InitPartyKickHandler, 2)
	boot.RegisterComponent("packethandler", entity_select.InitEntitySelectHandler, 2)

	return agentServer
}

func (a *AgentServer) Start() {
	go a.Server.Start()

	a.handlePackets()
}

func (a *AgentServer) handlePackets() {

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
				world := service.GetWorldServiceInstance()
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
	case viper.GetString(config.GatewayModuleId):
		a.GatewaySession = data.Session
		a.Server.Sessions[data.Session.ID] = nil
	}
}

func main() {
	config.Initialize()
	logging.Init()
	boot.SetPhases("gamedata", "services", "packethandler", "network")
	reader := bufio.NewReader(os.Stdin)

	loadGameData := func() {
		for k, v := range model.RefItems {
			model.RefObjects[k] = v.RefObject
		}
		for k, v := range model.RefChars {
			model.RefObjects[k] = v.RefObject
		}
		model.RefItems = model.GetAllRefItems()
		model.RefChars = model.GetAllRefChars()
	}

	boot.RegisterComponent("gamedata", loadGameData, 1)

	setupWorld := func() {
		world := service.GetWorldServiceInstance()
		world.LoadGameServerRegions(1)
	}

	boot.RegisterComponent("gamedata", setupWorld, 1)

	startServices := func() {
		manager.GetGameTimeManagerInstance().Start()
		manager.GetKnownListManager().Start()
		manager.GetSpawnManagerInstance().Start()
		manager.GetRespawnManagerInstance().Start()
	}

	boot.RegisterComponent("services", startServices, 1)

	log.Println("starting agent server...")

	agentServer := NewAgentServer()
	startAgentServer := func() {
		agentServer.Start()
	}

	boot.RegisterComponent("network", startAgentServer, 1)
	boot.Boot()

	log.Println("Press Enter to exit...")
	reader.ReadString('\n')
}
