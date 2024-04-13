package battle

import (
	"image/color"

	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
)

type floatingTextNode struct {
	pos        gmath.Vec
	text       string
	decaySpeed float64
	clr        graphics.ColorScale
	label      *graphics.Label
}

func newFloatingTextNode(pos gmath.Vec, text string, clr color.RGBA) *floatingTextNode {
	return &floatingTextNode{
		pos:  pos,
		text: text,
		clr:  graphics.ColorScaleFromColor(clr),
	}
}

func (t *floatingTextNode) IsDisposed() bool {
	return t.label.IsDisposed()
}

func (t *floatingTextNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	w := 128
	h := 32

	t.decaySpeed = ctx.Rand.FloatRange(0.9, 1.4)
	t.label = ctx.NewLabel(t.text, assets.Font1)
	t.label.SetSize(w, h)
	t.label.Pos.Base = &t.pos
	t.label.Pos.Offset.X -= float64(w / 2)
	t.label.Pos.Offset.Y -= float64(h / 2)
	t.label.SetAlignHorizontal(graphics.AlignHorizontalCenter)
	t.label.SetAlignVertical(graphics.AlignVerticalCenter)
	t.label.SetColorScale(t.clr)
	s.AddGraphics(t.label)
}

func (t *floatingTextNode) Update(delta float64) {
	alpha := t.label.GetAlpha() - float32(delta*t.decaySpeed)
	if alpha < 0.05 {
		t.label.Dispose()
		return
	}
	t.label.SetAlpha(alpha)
	t.pos.Y -= 100 * delta
}
