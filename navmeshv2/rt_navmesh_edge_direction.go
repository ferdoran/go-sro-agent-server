package navmeshv2

type RtNavmeshEdgeDirection int8

func (d RtNavmeshEdgeDirection) IsNone() bool {
	return d == -1
}

func (d RtNavmeshEdgeDirection) IsSouth() bool {
	return d == 0
}

func (d RtNavmeshEdgeDirection) IsWest() bool {
	return d == 1
}

func (d RtNavmeshEdgeDirection) IsNorth() bool {
	return d == 2
}

func (d RtNavmeshEdgeDirection) IsEast() bool {
	return d == 3
}
