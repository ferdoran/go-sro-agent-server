package model

type AggressType byte

const (
	Aggressive AggressType = 0
	Friendly   AggressType = 1
)

type SpawnTactics struct {
	ID                      int
	RefObjID                uint32
	AIQoS                   byte
	MaxStamina              int
	MaxStaminaVariance      int
	SightRange              int
	AggressType             byte
	AggressData             int  // TODO find out what this is
	ChangeTarget            byte // TODO find out what this is
	HelpRequestTo           byte // TODO find out what this is
	HelpReponseTo           byte // TODO find out what this is
	BattleStyle             byte // TODO find out what this is
	BattleStyleData         byte // TODO find out what this is
	DiversionBasis          byte // TODO find out what this is
	DiversionBasisData1     byte // TODO find out what this is
	DiversionBasisData2     byte // TODO find out what this is
	DiversionBasisData3     byte // TODO find out what this is
	DiversionBasisData4     byte // TODO find out what this is
	DiversionBasisData5     byte // TODO find out what this is
	DiversionBasisData6     byte // TODO find out what this is
	DiversionBasisData7     byte // TODO find out what this is
	DiversionBasisData8     byte // TODO find out what this is
	DiversionKeepBasis      byte // TODO find out what this is
	DiversionKeepBasisData1 byte // TODO find out what this is
	DiversionKeepBasisData2 byte // TODO find out what this is
	DiversionKeepBasisData3 byte // TODO find out what this is
	DiversionKeepBasisData4 byte // TODO find out what this is
	DiversionKeepBasisData5 byte // TODO find out what this is
	DiversionKeepBasisData6 byte // TODO find out what this is
	DiversionKeepBasisData7 byte // TODO find out what this is
	DiversionKeepBasisData8 byte // TODO find out what this is
	KeepDistance            byte // TODO find out what this is
	KeepDistanceData        byte // TODO find out what this is
	TraceType               byte // TODO find out what this is
	TraceBoundary           byte // TODO find out what this is
	TraceData               byte // TODO find out what this is
	HomingType              byte // TODO find out what this is
	HomingData              byte // TODO find out what this is
	AggressTypeOnHoming     byte // TODO find out what this is
	FleeType                byte // TODO find out what this is
	ChampionTacticsID       int
	AdditionOptionFlag      int
	DescString128           string
}
