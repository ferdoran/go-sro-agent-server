package model

const (
	WorldWidth  = 256 * RegionWidth
	WorldHeight = 128 * RegionHeight
)

type World struct {
	Continents []Continent
}
