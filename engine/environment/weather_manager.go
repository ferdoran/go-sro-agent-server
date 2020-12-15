package environment

import "github.com/ferdoran/go-sro-agent-server/model"

type WeatherManager struct{}

func GetCelestialPositionForPlayer(player *model.Player) CelestialPosition {
	// TODO how is it calculated?
	return CelestialPosition{
		CharUniqueID: player.UniqueID,
	}
}

func GetCurrentWeather() (weatherType model.WeatherType, intensity byte) {
	// TODO is it based on region?
	return model.Clear, 0
}
