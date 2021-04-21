package character

import (
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/ferdoran/go-sro-agent-server/service"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GameReadyHandler struct{}

func (h *GameReadyHandler) Handle(data server.PacketChannelData) {
	world := service.GetWorldServiceInstance()
	player, err := world.GetPlayerByUniqueId(data.UserContext.UniqueID)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to update player life state"))
		return
	}
	player.LifeState = model.Alive
	logrus.Debugf("Player %s's client is ready", player.CharName)
	// TODO tell all players around that character is not spawning anymore
}
