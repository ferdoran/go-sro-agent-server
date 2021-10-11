package model

import (
	"github.com/ferdoran/go-sro-agent-server/navmeshv2"
	"time"
)

type MovementData struct {
	StartTime      time.Time
	UpdateTime     time.Time
	TargetPosition navmeshv2.RtNavmeshPosition
	HasDestination bool
	DirectionAngle float32
}
