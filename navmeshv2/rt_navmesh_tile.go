package navmeshv2

const (
	RtNavmeshTileFlagNone    = 0
	RtNavmeshTileFlagBlocked = 1 << 0
	RtNavmeshTileFlagBit1    = 1 << 1
	RtNavmeshTileFlagBit2    = 1 << 2
	RtNavmeshTileFlagBit3    = 1 << 3
	RtNavmeshTileFlagBit4    = 1 << 4
	RtNavmeshTileFlagBit5    = 1 << 5
	RtNavmeshTileFlagBit6    = 1 << 6
	RtNavmeshTileFlagBit7    = 1 << 7
	RtNavmeshTileFlagBit8    = 1 << 8
	RtNavmeshTileFlagBit9    = 1 << 9
	RtNavmeshTileFlagBit10   = 1 << 10
	RtNavmeshTileFlagBit11   = 1 << 11
	RtNavmeshTileFlagBit12   = 1 << 12
	RtNavmeshTileFlagBit13   = 1 << 13
	RtNavmeshTileFlagBit14   = 1 << 14
	RtNavmeshTileFlagBit15   = 1 << 15

	TileWidth  = 20.0
	TileHeight = 20.0
)

type RtNavmeshTile struct {
	CellIndex int
	Flag      RtNavmeshTileFlag
	TextureID int16
}

func (t RtNavmeshTile) GetCellIndex() int {
	return t.CellIndex
}

func (t RtNavmeshTile) GetTextureID() int16 {
	return t.TextureID
}

type RtNavmeshTileFlag uint16

func (flag RtNavmeshTileFlag) IsBlocked() bool {
	return flag == RtNavmeshTileFlagBlocked
}
