package navmeshv2

import "github.com/ferdoran/go-sro-framework/utils"

const (
	RegionWidth  = 1920.0
	RegionHeight = 1920.0

	OriginX byte = 135
	OriginY byte = 92

	XSize   = 8
	XOffset = 0
	XMask   = ((1 << XSize) - 1) << XOffset

	YSize   = 7
	YOffset = XOffset + XSize
	YMask   = ((1 << YSize) - 1) << YOffset

	DungeonSize   = 1
	DungeonOffset = YOffset + YSize
	DungeonMask   = ((1 << DungeonSize) - 1) << DungeonOffset
)

var Origin Region = NewRegionFromXAndY(OriginX, OriginY)

type Region struct {
	value uint16
	ID    int16
	X     byte
	Y     byte
}

func NewRegionFromUint16(id uint16) Region {
	x, y := utils.Int16ToXAndZ(int16(id))

	return Region{
		value: id,
		ID:    int16(id),
		X:     byte(x),
		Y:     byte(y),
	}
}

func NewRegionFromInt16(id int16) Region {
	x, y := utils.Int16ToXAndZ(id)

	return Region{
		value: uint16(id),
		ID:    id,
		X:     byte(x),
		Y:     byte(y),
	}
}

func NewRegionFromXAndY(x, y byte) Region {
	id := utils.XAndZToInt16(x, y)

	isDungeon := (y & (1 << 7)) != 0
	var isDungeonByte = 0

	if isDungeon {
		isDungeonByte = 1
	} else {
		isDungeonByte = 0
	}
	val := uint16((0 & ^DungeonMask) | ((isDungeonByte << DungeonOffset) & DungeonMask))

	return Region{
		value: val,
		ID:    id,
		X:     x,
		Y:     y,
	}
}
