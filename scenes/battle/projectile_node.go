package battle

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/battle"
)

type projectileNode struct {
	data *battle.Projectile

	sprite *graphics.Sprite

	// This pos is bound during the projectile construction.
	targetPos gmath.Vec

	owner  *battle.Vessel
	target *battle.Vessel

	impactAreaSqr float64

	distLimit float64 // Dirty hack to avoid projectiles going through the target
	velocity  gmath.Vec
}

type projectileConfig struct {
	Data      *battle.Projectile
	Target    *battle.Vessel
	Owner     *battle.Vessel
	TargetPos gmath.Vec // Will be different from Target.Pos if using a snipe method
}

func newProjectileNode(config projectileConfig) *projectileNode {
	return &projectileNode{
		data:      config.Data,
		targetPos: config.TargetPos,
		distLimit: min(config.Data.Weapon.MaxRange, config.TargetPos.DistanceTo(config.Data.Pos)),
		target:    config.Target,
		owner:     config.Owner,
	}
}

func (p *projectileNode) Init(s *scene) {
	ctx := s.Controller().GetGameContext()

	p.impactAreaSqr = p.data.Weapon.ProjectileImpactArea * p.data.Weapon.ProjectileImpactArea

	p.sprite = ctx.NewSprite(p.data.Weapon.ProjectileImage)
	p.sprite.Pos.Base = &p.data.Pos
	p.sprite.Rotation = &p.data.Rotation
	s.AddGraphics(p.sprite)

	p.velocity = gmath.RadToVec(p.data.Rotation).Mulf(p.data.Weapon.ProjectileSpeed)
}

func (p *projectileNode) Update(delta float64) {
	v := p.velocity.Mulf(delta)
	p.distLimit -= v.Len()
	p.data.Pos = p.data.Pos.Add(v)
	if p.distLimit <= 0 {
		p.detonate()
	}
}

func (p *projectileNode) detonate() {
	if p.data.Pos.DistanceSquaredTo(p.targetPos) <= p.impactAreaSqr {
		if p.target.Health > 0 {
			p.target.OnDamage(p.data.Weapon.Damage, p.owner)
		}
	}
	p.Dispose()
}

func (p *projectileNode) Dispose() {
	p.sprite.Dispose()
}

func (p *projectileNode) IsDisposed() bool {
	return p.sprite.IsDisposed()
}
