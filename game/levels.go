package game

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ld55-game/assets"
)

type Level struct {
	Dist  float64
	Enemy resource.RawID
}

var Levels = []*Level{
	{
		Dist:  96,
		Enemy: assets.RawLevel1EnemyJSON,
	},
	{
		Dist:  64,
		Enemy: assets.RawLevel2EnemyJSON,
	},
	{
		Dist:  256,
		Enemy: assets.RawLevel3EnemyJSON,
	},
}
