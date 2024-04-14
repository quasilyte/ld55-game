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

	CollisionCheck bool
	Homing         float64

	FireSound resource.AudioID

	MaxRange float64

	ProjectileImage resource.ImageID

	Burst int

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

	DrainEnergy float64
	Slow        float64
}

type Projectile struct {
	Weapon   *WeaponDesign
	Pos      gmath.Vec
	Rotation gmath.Rad
}

var WeaponDesignList = []*WeaponDesign{
	{
		Name:                 "Pusher",
		EnergyCost:           10,
		Reload:               2.0,
		Damage:               Damage{Energy: 2, Kinetic: 2},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFirePusher1,
		ProjectileImage:      assets.ImageProjectilePusher,
		ProjectileSpeed:      500,
		MaxRange:             240,
		ProjectileImpactArea: 8,
		CollisionCheck:       true,
		ImpactImage:          assets.ImageImpactPusher,
		Burst:                1,
	},

	{
		Name:                 "Scatter Gun",
		EnergyCost:           0,
		Reload:               1.6,
		Damage:               Damage{Kinetic: 2},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFireScatter1,
		ProjectileImage:      assets.ImageProjectileScatter,
		ProjectileSpeed:      300,
		MaxRange:             240,
		ProjectileImpactArea: 5,
		Burst:                6,
	},

	{
		Name:                 "Pulse Laser",
		BuyCost:              120,
		EnergyCost:           5,
		Reload:               0.45,
		Damage:               Damage{Energy: 3},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFireLaser1,
		ProjectileImage:      assets.ImageProjectileLaser,
		ProjectileSpeed:      400,
		MaxRange:             250,
		ProjectileImpactArea: 8,
		ImpactImage:          assets.ImageImpactLaser,
		Burst:                1,
	},

	{
		Name:                 "Ion Cannon",
		EnergyCost:           4,
		Reload:               0.5,
		Damage:               Damage{Energy: 1, DrainEnergy: 6},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFireIon1,
		ProjectileImage:      assets.ImageProjectileIon,
		ProjectileSpeed:      450,
		MaxRange:             350,
		ProjectileImpactArea: 6,
		ImpactImage:          assets.ImageImpactIon,
		Burst:                1,
	},

	{
		Name:                 "Plasma Cannon",
		BuyCost:              150,
		EnergyCost:           4,
		Reload:               0.7,
		Damage:               Damage{Energy: 2, Thermal: 4},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFirePlasma1,
		ProjectileImage:      assets.ImageProjectilePlasma,
		ProjectileSpeed:      225,
		MaxRange:             300,
		ProjectileImpactArea: 10,
		ImpactImage:          assets.ImageImpactPlasma,
		Burst:                1,
	},

	{
		Name:                 "Lancer",
		EnergyCost:           16,
		Reload:               1.0,
		Damage:               Damage{Energy: 12},
		FiringType:           FixedAngleWeapon,
		FireSound:            assets.AudioFireLancer1,
		ProjectileImage:      assets.ImageProjectileLancer,
		ProjectileSpeed:      420,
		MaxRange:             400,
		ProjectileImpactArea: 8,
		CollisionCheck:       true,
		ImpactImage:          assets.ImageImpactLancer,
		Burst:                1,
	},

	{
		Name:                 "Freezer",
		EnergyCost:           20,
		Reload:               1,
		Damage:               Damage{Energy: 1, Slow: 4},
		FiringType:           TargetableWeapon,
		FireSound:            assets.AudioFireFreezer1,
		ProjectileImage:      assets.ImageProjectileFreezer,
		ProjectileSpeed:      300,
		MaxRange:             250,
		ProjectileImpactArea: 8,
		ImpactImage:          assets.ImageImpactFreezer,
		Burst:                1,
		CollisionCheck:       true,
		Homing:               120,
	},

	{
		Name:                 "Missile Launcher",
		Reload:               4,
		Damage:               Damage{Thermal: 5},
		FiringType:           FixedAngleWeapon,
		FireSound:            assets.AudioFireMissile1,
		ProjectileImage:      assets.ImageProjectileMissile,
		ProjectileSpeed:      160,
		MaxRange:             450,
		ProjectileImpactArea: 8,
		ImpactImage:          assets.ImageImpactMissile,
		Burst:                2,
		CollisionCheck:       true,
		Homing:               100,
	},
}
