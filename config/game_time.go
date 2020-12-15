package config

type GameTimeConfig struct {
	TicksPerSecond int     `json:"ticks_per_second"`
	DaySpeed       float32 `json:"day_speed"`
}
