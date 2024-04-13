package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/game"
)

type World struct {
	Vessels []*Vessel

	Size gmath.Vec

	Rand gmath.Rand
}

type Vessel struct {
	Alliance int

	Pos      gmath.Vec
	Velocity gmath.Vec
	Rotation gmath.Rad

	Prog *game.BotProg

	Design VesselDesign
}

type VesselDesign struct {
	Image resource.ImageID

	RotationSpeed gmath.Rad

	MaxSpeed     float64
	Acceleration float64
}
