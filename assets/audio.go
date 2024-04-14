package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
)

func registerAudioResources(loader *resource.Loader) {
	resources := map[resource.AudioID]resource.AudioInfo{
		AudioFireScatter1: {Path: "audio/scatter1.wav", Volume: -0.1},
		AudioFireScatter2: {Path: "audio/scatter2.wav", Volume: -0.1},

		AudioFireLaser1: {Path: "audio/laser1.wav", Volume: -0.4},
		AudioFireLaser2: {Path: "audio/laser2.wav", Volume: -0.4},
		AudioFireLaser3: {Path: "audio/laser3.wav", Volume: -0.4},
		AudioFireLaser4: {Path: "audio/laser4.wav", Volume: -0.4},

		AudioFirePlasma1: {Path: "audio/plasma1.wav", Volume: -0.2},
		AudioFirePlasma2: {Path: "audio/plasma2.wav", Volume: -0.2},
		AudioFirePlasma3: {Path: "audio/plasma3.wav", Volume: -0.2},

		AudioFireIon1: {Path: "audio/ion1.wav", Volume: -0.5},
		AudioFireIon2: {Path: "audio/ion2.wav", Volume: -0.5},

		AudioFireMissile1: {Path: "audio/missile1.wav", Volume: -0.3},
		AudioFireMissile2: {Path: "audio/missile2.wav", Volume: -0.3},

		AudioFireLancer1: {Path: "audio/lancer1.wav", Volume: -0.4},
		AudioFireLancer2: {Path: "audio/lancer2.wav", Volume: -0.4},
		AudioFireLancer3: {Path: "audio/lancer3.wav", Volume: -0.4},
		AudioFireLancer4: {Path: "audio/lancer4.wav", Volume: -0.4},

		AudioFirePusher1: {Path: "audio/pusher1.wav", Volume: -0.3},
		AudioFirePusher2: {Path: "audio/pusher2.wav", Volume: -0.3},
		AudioFirePusher3: {Path: "audio/pusher3.wav", Volume: -0.3},

		AudioFireFreezer1: {Path: "audio/freezer1.wav", Volume: -0.1},
		AudioFireFreezer2: {Path: "audio/freezer2.wav", Volume: -0.1},
		AudioFireFreezer3: {Path: "audio/freezer3.wav", Volume: -0.1},
	}

	for id, info := range resources {
		loader.AudioRegistry.Set(id, info)
		loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioFireScatter1, AudioFireIon1, AudioFireMissile1:
		return 2
	case AudioFireLaser1, AudioFireLancer1:
		return 4
	case AudioFirePlasma1, AudioFirePusher1, AudioFireFreezer1:
		return 3
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioFireScatter1
	AudioFireScatter2
	AudioFireLaser1
	AudioFireLaser2
	AudioFireLaser3
	AudioFireLaser4
	AudioFirePlasma1
	AudioFirePlasma2
	AudioFirePlasma3
	AudioFireIon1
	AudioFireIon2
	AudioFireLancer1
	AudioFireLancer2
	AudioFireLancer3
	AudioFireLancer4
	AudioFirePusher1
	AudioFirePusher2
	AudioFirePusher3
	AudioFireFreezer1
	AudioFireFreezer2
	AudioFireFreezer3
	AudioFireMissile1
	AudioFireMissile2
)
