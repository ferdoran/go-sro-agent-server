package service

import (
	"fmt"
	"github.com/ferdoran/go-sro-agent-server/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type EntitySelectService struct {
	EntityUniqueID uint32
	IsPlayerCharacter bool
	IsNPCNpc bool
	IsNPCMob bool
}

func (p *EntitySelectService) GetEntity(entitySelectRequest model.EntitySelectRequest) error {
	entityUniqueId := entitySelectRequest.EntityUniqueID;
	isPlayerCharacter := false
	isNPCNpc := false
	isNPCMob := false

	ws := GetWorldServiceInstance()

	visibleObject, err := ws.GetVisibleObjectByUniqueId(entityUniqueId)
	if err != nil {
		log.Debugln("Failed to find a visible object with the UniqueID:", entityUniqueId)
	} else {
		if visibleObject.GetTypeInfo().IsPlayerCharacter() {
			isPlayerCharacter = true
		}
		if visibleObject.GetTypeInfo().IsNPCNpc() {
			isNPCNpc = true
		}
		if visibleObject.GetTypeInfo().IsNPCMob() {
			isNPCMob = true
		}
	}

	if !isPlayerCharacter && !isNPCNpc && !isNPCMob {
		return errors.New(fmt.Sprintf("Given UniqueID %d does not belong to a player or an NPC/mob", entityUniqueId))
	} else {
		p.EntityUniqueID    = entityUniqueId
		p.IsPlayerCharacter = isPlayerCharacter
		p.IsNPCNpc          = isNPCNpc
		p.IsNPCMob          = isNPCMob
		log.Debugf("Visible obj:\n%#v\n", p)
		return nil
	}
}