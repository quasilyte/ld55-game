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

	if n.data.Alliance != 0 {
		n.sprite.SetColorScale(graphics.ColorScale{R: 1.1, G: 0.5, B: 1.2, A: 1})
	}
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
	deceleration := 0.1

	if n.commands.RotateLeft || n.commands.RotateRight {
		deceleration = 0.4
	}

	if n.commands.MoveForward {
		accel := n.data.Design.Acceleration * delta
		accelVector := gmath.RadToVec(n.data.Rotation).Mulf(accel)
		n.data.EngineVelocity = n.data.EngineVelocity.Add(accelVector)
		n.data.EngineVelocity = n.data.EngineVelocity.ClampLen(n.data.Design.MaxSpeed)
	} else {
		n.data.EngineVelocity = n.data.EngineVelocity.Mulf(1 - (delta * deceleration))
	}

	if !n.data.ExtraVelocity.IsZero() {
		n.data.ExtraVelocity = n.data.ExtraVelocity.Mulf(1 - (delta * 0.5))
	}

	v := n.data.Velocity()
	if !v.IsZero() {
		n.data.Pos = n.data.Pos.Add(v.Mulf(delta))
	}
}

func (n *vesselNode) processRotation(delta float64) {
	rotationMultiplier := gmath.Rad(1.0)
	if n.commands.MoveForward {
		rotationMultiplier = 0.7
	}

	rotation := gmath.Rad(0)
	if n.commands.RotateLeft {
		rotation = -n.data.Design.RotationSpeed * rotationMultiplier
	}
	if n.commands.RotateRight {
		rotation = n.data.Design.RotationSpeed * rotationMultiplier
	}
	n.data.Rotation += rotation * gmath.Rad(delta)
}
