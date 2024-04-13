package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ld55-game/game"
)

type Vessel struct {
	Alliance int

	Pos            gmath.Vec
	EngineVelocity gmath.Vec
	ExtraVelocity  gmath.Vec
	Rotation       gmath.Rad

	Health float64
	Energy float64

	Prog *game.BotProg

	Target *Vessel

	Weapons []*Weapon

	Design VesselDesign

	EventDestroyed gsignal.Event[*Vessel] // Vessel arg is attacker
}

func (v *Vessel) Velocity() gmath.Vec {
	return v.EngineVelocity.Add(v.ExtraVelocity)
}

func (v *Vessel) OnDamage(d Damage, attacker *Vessel) {
	totalDamage := ((1.0 - v.Design.EnergyResist) * d.Energy) +
		((1.0 - v.Design.KineticResist) * d.Kinetic) +
		((1.0 - v.Design.ThermalResist) * d.Thermal)
	if totalDamage > 0 {
		v.Health = gmath.ClampMin(v.Health-totalDamage, 0)
		if v.Health == 0 {
			v.EventDestroyed.Emit(attacker)
		}
	}
}

type VesselDesign struct {
	Image resource.ImageID

	RotationSpeed gmath.Rad

	MaxHealth float64
	MaxEnergy float64

	EnergyRegen float64

	MaxSpeed     float64
	Acceleration float64

	EnergyResist  float64
	KineticResist float64
	ThermalResist float64

	Weapons []*WeaponDesign
}
