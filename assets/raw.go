package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

func registerRawResources(loader *resource.Loader) {
	resources := map[resource.RawID]resource.RawInfo{
		RawLevel1EnemyJSON: {Path: "levels/level1_enemy.json"},
		RawLevel2EnemyJSON: {Path: "levels/level2_enemy.json"},
		RawLevel3EnemyJSON: {Path: "levels/level3_enemy.json"},
		RawLevel4EnemyJSON: {Path: "levels/level4_enemy.json"},
		RawLevel5EnemyJSON: {Path: "levels/level5_enemy.json"},
		RawLevel6EnemyJSON: {Path: "levels/level6_enemy.json"},
	}

	for id, info := range resources {
		loader.RawRegistry.Set(id, info)
		loader.LoadRaw(id)
	}
}

const (
	RawNone resource.RawID = iota

	RawLevel1EnemyJSON
	RawLevel2EnemyJSON
	RawLevel3EnemyJSON
	RawLevel4EnemyJSON
	RawLevel5EnemyJSON
	RawLevel6EnemyJSON
)
