package model

type PartyFormRequest struct {
	MasterUniqueID    uint32
	MasterJID         uint32
	MasterName        string
	CountryType       byte
	PartySettingsFlag PartySetting
	PurposeType       PartyPurpose
	LevelMin          byte
	LevelMax          byte
	Title             string
}
