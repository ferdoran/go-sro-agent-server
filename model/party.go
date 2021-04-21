package model

import (
	"sync"
)

type Party struct {
	Number            uint32
	MasterUniqueID    uint32
	MasterJID         uint32
	MasterName        string
	CountryType       byte
	MemberCount       byte
	PartySettingsFlag PartySetting
	PurposeType       PartyPurpose
	LevelMin          byte
	LevelMax          byte
	Title             string
	Members           []uint32
	Mutex             *sync.Mutex
}

type PartyJoinRequest struct {
	PartyNumber    uint32
	RequestID      uint32
	PlayerJID      int
	PlayerUniqueID uint32
	AcceptResult   bool
	Mutex          *sync.Mutex
}
