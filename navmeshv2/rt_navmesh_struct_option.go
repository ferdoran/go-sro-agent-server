package navmeshv2

const (
	RtNavmeshStructOptionNone = iota
	RtNavmeshStructOptionEdge
	RtNavmeshStructOptionCell
	_
	RtNavmeshStructOptionEvent
)

type RtNavmeshStructOption int

func (option RtNavmeshStructOption) IsNone() bool {
	return option&RtNavmeshStructOptionNone != 0
}

func (option RtNavmeshStructOption) IsEdge() bool {
	return option&RtNavmeshStructOptionEdge != 0
}

func (option RtNavmeshStructOption) IsCell() bool {
	return option&RtNavmeshStructOptionCell != 0
}

func (option RtNavmeshStructOption) IsEvent() bool {
	return option&RtNavmeshStructOptionEvent != 0
}
