package navmeshv2

type RtNavmesh interface {
	GetNavmeshType() RtNavmeshType
	GetFilename() string
	GetCell(index int) RtNavmeshCell
}

type RtNavmeshBase struct {
	Filename string
}

func (base RtNavmeshBase) GetFilename() string {
	return base.Filename
}
