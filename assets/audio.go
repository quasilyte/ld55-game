package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

func registerAudioResources(loader *resource.Loader) {
	resources := map[resource.AudioID]resource.AudioInfo{
		AudioFireLaser1: {Path: "audio/laser1.wav", Volume: -0.3},
		AudioFireLaser2: {Path: "audio/laser2.wav", Volume: -0.3},
		AudioFireLaser3: {Path: "audio/laser3.wav", Volume: -0.3},
		AudioFireLaser4: {Path: "audio/laser4.wav", Volume: -0.3},

		AudioFirePlasma1: {Path: "audio/plasma1.wav", Volume: -0.3},
		AudioFirePlasma2: {Path: "audio/plasma2.wav", Volume: -0.3},
		AudioFirePlasma3: {Path: "audio/plasma3.wav", Volume: -0.3},
	}

	for id, info := range resources {
		loader.AudioRegistry.Set(id, info)
		loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioFireLaser1:
		return 4
	case AudioFirePlasma1:
		return 3
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioFireLaser1
	AudioFireLaser2
	AudioFireLaser3
	AudioFireLaser4
	AudioFirePlasma1
	AudioFirePlasma2
	AudioFirePlasma3
)
