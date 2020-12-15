package model

import "time"

type MovementData struct {
	StartTime      time.Time
	UpdateTime     time.Time
	TargetPosition Position
	HasDestination bool
	DirectionAngle float32
}
