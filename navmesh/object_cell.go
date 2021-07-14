package navmesh

import "github.com/ferdoran/go-sro-framework/math"

type ObjectCell struct {
	*math.Triangle
	Index int
	Flag  uint16
}
