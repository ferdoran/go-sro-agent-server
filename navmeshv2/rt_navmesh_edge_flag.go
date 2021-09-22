package navmeshv2

type RtNavmeshEdgeFlag byte

func (f RtNavmeshEdgeFlag) IsNone() bool {
	return f == 0
}

func (f RtNavmeshEdgeFlag) IsBlockedDstToSrc() bool {
	return f&(1<<0) != 0
}

func (f RtNavmeshEdgeFlag) IsBlockedSrcToDst() bool {
	return f&(1<<1) != 0
}

func (f RtNavmeshEdgeFlag) IsBlocked() bool {
	return f&((1<<0)|(1<<1)) != 0
}

func (f RtNavmeshEdgeFlag) IsInternal() bool {
	return f&(1<<2) != 0
}

func (f RtNavmeshEdgeFlag) IsGlobal() bool {
	return f&(1<<3) != 0
}

func (f RtNavmeshEdgeFlag) IsBridge() bool {
	return f&(1<<4) != 0
}

func (f RtNavmeshEdgeFlag) IsEntrance() bool {
	return f&(1<<5) != 0
}

func (f RtNavmeshEdgeFlag) IsBit6() bool {
	return f&(1<<6) != 0
}

func (f RtNavmeshEdgeFlag) IsSiege() bool {
	return f&(1<<7) != 0
}
