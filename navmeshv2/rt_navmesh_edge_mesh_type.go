package navmeshv2

const (
	RtNavmeshEdgeMeshTypeTerrain = iota
	RtNavmeshEdgeMeshTypeObject
)

type RtNavmeshEdgeMeshType byte

func (edgeType RtNavmeshEdgeMeshType) IsTerrain() bool {
	return edgeType == RtNavmeshEdgeMeshTypeTerrain
}

func (edgeType RtNavmeshEdgeMeshType) IsObject() bool {
	return edgeType == RtNavmeshEdgeMeshTypeObject
}
