package model

const (
	TypeID1Item     = 3
	TypeID2Wearable = 1
	TypeID3Weapon   = 6
)

type IItem interface {
	ISRObject
	GetOwner() uint32
	GetRarity() byte
}

type Item struct {
	SRObject
	ID       int
	Name     string
	Variance uint64
	TradeInfo
	PriceInfo
	LevelInfo
	StackSize  int
	SpeedClass int
}

func (i *Item) GetType() string {
	return "Item"
}
