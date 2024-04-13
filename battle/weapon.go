package battle

import (
	"slices"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
)

func FindWeaponDesignByName(name string) *WeaponDesign {
	i := slices.IndexFunc(Weapons, func(w *WeaponDesign) bool {
		return w.Name == name
	})
	return Weapons[i]
}

type Weapon struct {
	Reload float64 // Cooldown before it can fire again

	Design *WeaponDesign
}

type WeaponDesign struct {
	Name string

	EnergyCost float64 // Per shot

	Reload float64 // Base value (doesn't take multipliers into account)

	Damage Damage

	FiringType WeaponFiringType

	MaxRange float64

	ProjectileImage resource.ImageID

	ProjectileSpeed float64

	ProjectileImpactArea float64

	ImpactImage resource.ImageID
}

type WeaponFiringType int

const (
	TargetableWeapon WeaponFiringType = iota
	FixedAngleWeapon
)

type Damage struct {
	Energy  float64
	Kinetic float64
	Thermal float64
}

type Projectile struct {
	Weapon   *WeaponDesign
	Pos      gmath.Vec
	Rotation gmath.Rad
}

var Weapons = []*WeaponDesign{
	{
		Name:                 "Pulse Laser",
		EnergyCost:           5,
		Reload:               0.5,
		Damage:               Damage{Energy: 3},
		FiringType:           TargetableWeapon,
		ProjectileImage:      assets.ImageProjectileLaser,
		ProjectileSpeed:      400,
		MaxRange:             250,
		ProjectileImpactArea: 3,
		ImpactImage:          assets.ImageImpactLaser,
	},
}
