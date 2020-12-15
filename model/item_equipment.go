package model

type IEquipment interface {
	IItem
	GetItemClass() int
	SetItemClass(itemClass int)
	GetOptLevel() byte
}

type Equipment struct {
	Item
	AlchemyInfo
	ItemClass int
	OptLevel  byte
}

func (e *Equipment) GetItemClass() int {
	return e.ItemClass
}

func (e *Equipment) SetItemClass(itemClass int) {
	e.ItemClass = itemClass
}

func (e *Equipment) GetOptLevel() byte {
	return e.OptLevel
}
