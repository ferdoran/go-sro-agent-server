package navmesh

type NavMeshData struct {
	ObjectList
	NavigationCells
	NavMeshGlobalEdges
	NavMeshInternalEdges
	TileMap
	HeightMap
	SurfaceTypeMap
	SurfaceHeightMap
}
