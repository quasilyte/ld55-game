package game

import (
	"slices"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ld55-game/assets"
)

func FindWeaponDesignByName(name string) *WeaponDesign {
	i := slices.IndexFunc(WeaponDesignList, func(w *WeaponDesign) bool {
		return w.Name == name
	})
	return WeaponDesignList[i]
}

type Weapon struct {
	Reload float64 // Cooldown before it can fire again

	Design *WeaponDesign
}

type WeaponDesign struct {
	Name string

	EnergyCost float64 // Per shot

	BuyCost int

	Reload float64 // Base value (doesn't take multipliers into account)

	Damage Damage

	FiringType WeaponFiringType

	FireSound resource.AudioID

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

var WeaponDesignList = []*WeaponDesign{
	{
		Name:                 "Pulse Laser",
		BuyCost:              120,
		EnergyCost:           5,
		Reload:               0.5,
		Damage:               Damage{Energy: 3},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFireLaser1,
		ProjectileImage:      assets.ImageProjectileLaser,
		ProjectileSpeed:      400,
		MaxRange:             250,
		ProjectileImpactArea: 8,
		ImpactImage:          assets.ImageImpactLaser,
	},

	{
		Name:                 "Plasma Cannon",
		BuyCost:              150,
		EnergyCost:           4,
		Reload:               0.7,
		Damage:               Damage{Energy: 1, Thermal: 3},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFirePlasma1,
		ProjectileImage:      assets.ImageProjectilePlasma,
		ProjectileSpeed:      250,
		MaxRange:             300,
		ProjectileImpactArea: 10,
		ImpactImage:          assets.ImageImpactPlasma,
	},
}
