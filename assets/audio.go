package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

func registerAudioResources(loader *resource.Loader) {
	resources := map[resource.AudioID]resource.AudioInfo{
		AudioImpactLaser1: {Path: "audio/laser1.wav"},
		AudioImpactLaser2: {Path: "audio/laser2.wav"},
		AudioImpactLaser3: {Path: "audio/laser3.wav"},
		AudioImpactLaser4: {Path: "audio/laser4.wav"},
	}

	for id, info := range resources {
		loader.AudioRegistry.Set(id, info)
		loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioImpactLaser1:
		return 4
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioImpactLaser1
	AudioImpactLaser2
	AudioImpactLaser3
	AudioImpactLaser4
)
