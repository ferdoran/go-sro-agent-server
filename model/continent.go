package model

type Continent struct {
	Regions map[uint16]Region
	Players map[uint32]Player
}
