package model

import (
	"sync"
)

type ExchangeRequest struct {
	RequestingUniqueID uint32
	RequestedUniqueID  uint32
	IsStarted          bool
	Mutex              *sync.Mutex
}