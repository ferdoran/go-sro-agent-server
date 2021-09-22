package navmeshv2

const (
	RtNavmeshTypeNone = iota
	RtNavmeshTypeTerrain
	RtNavmeshTypeObject
	RtNavmeshTypeDungeon
)

type RtNavmeshType byte

func (t RtNavmeshType) IsNone() bool {
	return t == RtNavmeshTypeNone
}

func (t RtNavmeshType) IsTerrain() bool {
	return t == RtNavmeshTypeTerrain
}

func (t RtNavmeshType) IsObject() bool {
	return t == RtNavmeshTypeObject
}

func (t RtNavmeshType) IsDungeon() bool {
	return t == RtNavmeshTypeDungeon
}
