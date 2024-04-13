package battle

import (
	"github.com/quasilyte/gmath"
)

type World struct {
	Vessels []*Vessel

	Size gmath.Vec

	Rand gmath.Rand
}
