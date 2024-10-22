package battle

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
)

type projectileNode struct {
	data *game.Projectile

	scene *scene

	sprite *graphics.Sprite

	// This pos is bound during the projectile construction.
	targetPos gmath.Vec

	owner  *game.Vessel
	target *game.Vessel

	impactAreaSqr float64

	distLimit float64 // Dirty hack to avoid projectiles going through the target
	velocity  gmath.Vec
}

type projectileConfig struct {
	Data      *game.Projectile
	Target    *game.Vessel
	Owner     *game.Vessel
	TargetPos gmath.Vec // Will be different from Target.Pos if using a snipe method
}

func newProjectileNode(config projectileConfig) *projectileNode {
	return &projectileNode{
		data:      config.Data,
		targetPos: config.TargetPos,
		distLimit: min(config.Data.Weapon.MaxRange, config.TargetPos.DistanceTo(config.Data.Pos)) * 1.1,
		target:    config.Target,
		owner:     config.Owner,
	}
}

func (p *projectileNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	p.scene = s

	p.impactAreaSqr = p.data.Weapon.ProjectileImpactArea * p.data.Weapon.ProjectileImpactArea

	p.sprite = ctx.NewSprite(p.data.Weapon.ProjectileImage)
	p.sprite.Pos.Base = &p.data.Pos
	p.sprite.Rotation = &p.data.Rotation
	s.AddGraphics(p.sprite)

	p.velocity = gmath.RadToVec(p.data.Rotation).Mulf(p.data.Weapon.ProjectileSpeed)
}

func (p *projectileNode) Update(delta float64) {

	var v gmath.Vec
	if p.data.Weapon.Homing == 0 {
		v = p.velocity.Mulf(delta)
		p.data.Pos = p.data.Pos.Add(v)
	} else {
		accel := p.seek()
		p.velocity = p.velocity.Add(accel.Mulf(delta)).ClampLen(p.data.Weapon.ProjectileSpeed)
		p.data.Rotation = p.velocity.Angle()
		p.data.Pos = p.data.Pos.Add(p.velocity.Mulf(delta))
		v = p.velocity.Mulf(delta)
	}
	travelled := v.Len()
	p.distLimit -= v.Len()

	if p.distLimit <= 0 || p.data.Pos.DistanceTo(p.targetPos) < travelled {
		p.detonate()
	}

	if !p.data.Weapon.CollisionCheck {
		return
	}

	hitboxSize := p.data.Weapon.ProjectileImpactArea + p.target.Design.HitboxSize
	if p.data.Pos.DistanceTo(p.target.Pos) < hitboxSize {
		p.detonate()
	}
}

func (p *projectileNode) seek() gmath.Vec {
	dst := p.target.Pos.Sub(p.data.Pos).Normalized().Mulf(p.data.Weapon.ProjectileSpeed)
	return dst.Sub(p.velocity).Normalized().Mulf(p.data.Weapon.Homing)
}

func (p *projectileNode) detonate() {
	hitboxSize := p.data.Weapon.ProjectileImpactArea + p.target.Design.HitboxSize
	if p.data.Pos.DistanceTo(p.target.Pos) <= hitboxSize {
		if p.target.Health > 0 {
			p.target.OnDamage(p.data.Weapon.Damage, p.owner)
			p.createImpactEffect()
			if p.data.Weapon.Name == "Pusher" {
				p.target.ExtraVelocity = p.target.ExtraVelocity.Add(p.velocity.Mulf(0.2))
			}
		}
	}
	p.Dispose()
}

func (p *projectileNode) createImpactEffect() {
	if p.data.Weapon.ImpactImage == assets.ImageNone {
		return
	}
	effect := newEffectNode(p.data.Pos, p.data.Weapon.ImpactImage)
	p.scene.AddObject(effect)
}

func (p *projectileNode) Dispose() {
	p.sprite.Dispose()
}

func (p *projectileNode) IsDisposed() bool {
	return p.sprite.IsDisposed()
}
