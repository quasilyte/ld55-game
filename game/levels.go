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
		// Dist:  160,
		// Enemy: assets.RawLevel1EnemyJSON,
		Dist:  320,
		Enemy: assets.RawLevel5EnemyJSON,
	},
	{
		Dist:  64,
		Enemy: assets.RawLevel2EnemyJSON,
	},
	{
		Dist:  256,
		Enemy: assets.RawLevel3EnemyJSON,
	},
	{
		Dist:  32,
		Enemy: assets.RawLevel4EnemyJSON,
	},
	{
		Dist:  320,
		Enemy: assets.RawLevel5EnemyJSON,
	},
}
