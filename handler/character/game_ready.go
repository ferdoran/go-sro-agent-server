package character

import (
	"github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/agent-server/model"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

type GameReadyHandler struct{}

func (h *GameReadyHandler) Handle(data server.PacketChannelData) {
	world := model.GetSroWorldInstance()
	player := world.PlayersByUniqueId[data.UserContext.UniqueID]
	player.LifeState = model.Alive
	logrus.Debugf("Player %s's client is ready", player.CharName)
	// TODO tell all players around that character is not spawning anymore
}
