package battle

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ld55-game/gfx"
)

type effectNode struct {
	pos     gmath.Vec
	image   resource.ImageID
	anim    *gfx.Animation
	rotates bool
	noFlip  bool

	rotation gmath.Rad

	EventCompleted gsignal.Event[gsignal.Void]
}

func newEffectNode(pos gmath.Vec, image resource.ImageID) *effectNode {
	return &effectNode{
		pos:   pos,
		image: image,
	}
}

func (e *effectNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	var sprite *graphics.Sprite
	if e.anim == nil {
		sprite = ctx.NewSprite(e.image)
		sprite.Pos.Base = &e.pos
	} else {
		sprite = e.anim.Sprite()
	}
	sprite.Rotation = &e.rotation
	if !e.noFlip {
		sprite.SetHorizontalFlip(ctx.Rand.Bool())
	}
	s.AddGraphics(sprite)
	if e.anim == nil {
		e.anim = gfx.NewAnimation(sprite, -1)
	}
}

func (e *effectNode) IsDisposed() bool {
	return e.anim.IsDisposed()
}

func (e *effectNode) Dispose() {
	e.anim.Sprite().Dispose()
}

func (e *effectNode) Update(delta float64) {
	if e.anim.Tick(delta) {
		e.EventCompleted.Emit(gsignal.Void{})
		e.Dispose()
		return
	}
	if e.rotates {
		e.rotation += gmath.Rad(delta * 2)
	}
}
