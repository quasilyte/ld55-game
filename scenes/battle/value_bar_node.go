package battle

import (
	"image/color"

	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
)

type valueBarNode struct {
	pos      gmath.Pos
	maxValue float64
	value    *float64
	clr      graphics.ColorScale
	rect     *graphics.Rect
}

func newValueBarNode(v *float64, maxValue float64, pos gmath.Pos, clr color.RGBA) *valueBarNode {
	return &valueBarNode{
		pos:      pos,
		value:    v,
		maxValue: maxValue,
		clr:      graphics.ColorScaleFromColor(clr),
	}
}

func (n *valueBarNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	n.rect = ctx.NewRect(64, 3)
	n.rect.Pos = n.pos
	n.rect.SetFillColorScale(n.clr)
	s.AddGraphics(n.rect)
}

func (n *valueBarNode) Update(delta float64) {
	percent := *n.value / n.maxValue
	width := 64.0 * percent
	n.rect.SetWidth(width)
}

func (n *valueBarNode) Dispose() {
	n.rect.Dispose()
}

func (n *valueBarNode) IsDisposed() bool {
	return n.rect.IsDisposed()
}
