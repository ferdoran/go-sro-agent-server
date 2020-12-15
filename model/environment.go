package model

type WeatherType byte

const (
	Clear WeatherType = iota + 1
	Rain
	Snow
)
