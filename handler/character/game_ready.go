package character

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/sirupsen/logrus"
)

type GameReadyHandler struct{}

func (h *GameReadyHandler) Handle(data server.PacketChannelData) {
	world := model.GetSroWorldInstance()
	player := world.PlayersByUniqueId[data.UserContext.UniqueID]
	player.LifeState = model.Alive
	logrus.Debugf("Player %s's client is ready", player.CharName)
	// TODO tell all players around that character is not spawning anymore
}
