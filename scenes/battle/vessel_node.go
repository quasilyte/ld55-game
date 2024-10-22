package battle

import (
	"fmt"

	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/progsim"
	"github.com/quasilyte/ld55-game/styles"
)

type vesselNode struct {
	scene *scene

	data  *game.Vessel
	world *game.World

	sprite    *graphics.Sprite
	aura      *graphics.Rect
	hpBar     *valueBarNode
	energyBar *valueBarNode

	commands progsim.VesselCommands
}

func newVesselNode(world *game.World, data *game.Vessel) *vesselNode {
	return &vesselNode{
		world: world,
		data:  data,
	}
}

func (n *vesselNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	n.scene = s

	n.sprite = ctx.NewSprite(n.data.Design.Image)
	n.sprite.Pos.Base = &n.data.Pos
	n.sprite.Rotation = &n.data.Rotation
	s.AddGraphics(n.sprite)

	n.aura = ctx.NewRect(64, 64)
	n.aura.Pos.Base = &n.data.Pos
	s.AddGraphics(n.aura)
	auraColor := styles.AlliedColor
	if n.data.Alliance != 0 {
		auraColor = styles.EnemyColor
	}
	auraColor.A = 0xff / 2
	n.aura.SetFillColorScale(graphics.ColorScale{}) // Transparent
	n.aura.SetOutlineColorScale(graphics.ColorScaleFromColor(auraColor))

	{
		pos := gmath.Pos{Base: &n.data.Pos, Offset: gmath.Vec{Y: 37}}
		n.hpBar = newValueBarNode(&n.data.Health, n.data.Design.MaxHealth, pos, styles.HealthColor)
		s.AddObject(n.hpBar)
	}
	{
		pos := gmath.Pos{Base: &n.data.Pos, Offset: gmath.Vec{Y: 43}}
		n.energyBar = newValueBarNode(&n.data.Energy, n.data.Design.MaxEnergy, pos, styles.EnergyColor)
		s.AddObject(n.energyBar)
	}

	if n.data.Alliance != 0 {
		n.sprite.SetColorScale(graphics.ColorScale{R: 1.1, G: 0.5, B: 1.2, A: 1})
	}
}

func (n *vesselNode) Dispose() {
	n.sprite.Dispose()
	n.aura.Dispose()
	n.hpBar.Dispose()
	n.energyBar.Dispose()
}

func (n *vesselNode) IsDisposed() bool {
	return n.sprite.IsDisposed()
}

func (n *vesselNode) SetCommands(c progsim.VesselCommands) {
	n.commands = c
}

func (n *vesselNode) Update(delta float64) {
	{
		rect := gmath.Rect{Max: n.world.Size}
		if !rect.Contains(n.data.Pos) {
			dmg := game.Damage{Thermal: 2 * delta}
			n.data.OnDamage(dmg, nil)
		}
	}

	n.data.Energy = gmath.ClampMax(n.data.Energy+(n.data.Design.EnergyRegen*delta), n.data.Design.MaxEnergy)
	n.data.Slow = gmath.ClampMin(n.data.Slow-delta, 0)

	n.processRotation(delta)
	n.processMovement(delta)
	n.processWeapons(delta)
}

func (n *vesselNode) processMovement(delta float64) {
	deceleration := 0.25

	if n.commands.RotateLeft || n.commands.RotateRight {
		deceleration = 0.6
	}

	if n.commands.MoveForward {
		accel := n.data.Design.Acceleration * delta
		maxSpeed := n.data.Design.MaxSpeed
		if n.data.Slow > 0 {
			accel *= 0.5
			switch {
			case n.data.Slow >= 10:
				maxSpeed *= 0.5
			case n.data.Slow >= 5:
				maxSpeed *= 0.6
			default:
				maxSpeed *= 0.7
			}
		}
		accelVector := gmath.RadToVec(n.data.Rotation).Mulf(accel)
		n.data.EngineVelocity = n.data.EngineVelocity.Add(accelVector)
		n.data.EngineVelocity = n.data.EngineVelocity.ClampLen(maxSpeed)
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

func (n *vesselNode) processWeapons(delta float64) {
	for _, w := range n.data.Weapons {
		w.Reload = gmath.ClampMin(w.Reload-delta, 0)
	}

	ctx := n.scene.Controller().GetGameContext()

	for _, c := range n.commands.FireCommands {
		if c.WeaponIndex >= uint(len(n.data.Weapons)) {
			fmt.Printf("warning: invalid weapon index %d\n", c.WeaponIndex)
			continue
		}

		w := n.data.Weapons[c.WeaponIndex]
		if w.Reload > 0 {
			continue
		}
		if n.data.Energy < w.Design.EnergyCost {
			continue
		}

		for i := 0; i < w.Design.Burst; i++ {
			// TODO: handle different weapon fire modes, etc.
			targetPos := c.TargetPos
			firePos := n.data.Pos
			rotation := n.data.Rotation
			if w.Design.FiringType == game.TargetableWeapon {
				if i != 0 {
					targetPos = targetPos.Add(ctx.Rand.Offset(-28, 28))
					firePos = firePos.Add(ctx.Rand.Offset(-8, 8))
				}
				rotation = n.data.Pos.AngleToPoint(targetPos)
			} else {
				targetPos = n.data.Pos.MoveInDirection(w.Design.MaxRange, rotation)
				if i != 0 {
					targetPos = targetPos.Add(ctx.Rand.Offset(-28, 28))
					firePos = firePos.Add(ctx.Rand.Offset(-8, 8))
				}
			}
			pd := &game.Projectile{
				Pos:      firePos,
				Rotation: rotation,
				Weapon:   w.Design,
			}
			p := newProjectileNode(projectileConfig{
				Data:      pd,
				TargetPos: targetPos,
				Target:    n.data.Target,
				Owner:     n.data,
			})
			n.scene.AddObject(p)
		}

		if w.Design.FireSound != assets.AudioNone {
			playSound(n.scene, w.Design.FireSound)
		}

		n.data.Energy -= w.Design.EnergyCost
		w.Reload = w.Design.Reload * c.ReloadMultiplier
	}
}
