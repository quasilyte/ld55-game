package battle

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/battle"
	"github.com/quasilyte/ld55-game/progsim"
)

type vesselNode struct {
	scene *scene

	data *battle.Vessel

	sprite *graphics.Sprite

	commands progsim.VesselCommands
}

func newVesselNode(data *battle.Vessel) *vesselNode {
	return &vesselNode{
		data: data,
	}
}

func (n *vesselNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	n.scene = s

	n.sprite = ctx.NewSprite(n.data.Design.Image)
	n.sprite.Pos.Base = &n.data.Pos
	n.sprite.Rotation = &n.data.Rotation
	s.AddGraphics(n.sprite)
}

func (n *vesselNode) Dispose() {
	n.sprite.Dispose()
}

func (n *vesselNode) IsDisposed() bool {
	return false
}

func (n *vesselNode) SetCommands(c progsim.VesselCommands) {
	n.commands = c
}

func (n *vesselNode) Update(delta float64) {
	n.processRotation(delta)
	n.processMovement(delta)
}

func (n *vesselNode) processMovement(delta float64) {
	v := n.data.Velocity

	if n.commands.MoveForward {
		speed := n.data.Design.MaxSpeed * delta
		v = gmath.RadToVec(n.data.Rotation).Mulf(speed)
	} else {
		v = gmath.Vec{}
	}

	n.data.Velocity = v
	if !v.IsZero() {
		n.data.Pos = n.data.Pos.Add(v)
	}
}

func (n *vesselNode) processRotation(delta float64) {
	rotation := gmath.Rad(0)
	if n.commands.RotateLeft {
		rotation = -n.data.Design.RotationSpeed
	}
	if n.commands.RotateRight {
		rotation = n.data.Design.RotationSpeed
	}
	n.data.Rotation += rotation * gmath.Rad(delta)
}
