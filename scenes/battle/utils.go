package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ld55-game/assets"
)

func playSound(s *scene, id resource.AudioID) {
	numSamples := assets.NumSamples(id)
	ctx := s.Controller().GetGameContext()
	if numSamples == 1 {
		ctx.Audio().PlaySound(id)
	} else {
		soundIndex := ctx.Rand.IntRange(0, numSamples-1)
		sound := resource.AudioID(int(id) + soundIndex)
		ctx.Audio().PlaySound(sound)
	}
}
