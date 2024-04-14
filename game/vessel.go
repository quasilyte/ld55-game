package game

import (
	"slices"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ld55-game/assets"
)

type Vessel struct {
	Alliance int

	Pos            gmath.Vec
	EngineVelocity gmath.Vec
	ExtraVelocity  gmath.Vec
	Rotation       gmath.Rad

	Health float64
	Energy float64

	Prog *BotProg

	Target *Vessel

	Weapons []*Weapon

	Design VesselDesign

	EventOnDamage  gsignal.Event[float64]
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
		v.EventOnDamage.Emit(totalDamage)
		v.Health = gmath.ClampMin(v.Health-totalDamage, 0)
		if v.Health == 0 {
			v.EventDestroyed.Emit(attacker)
		}
	}
}

func FindVesselDesignByName(name string) *VesselDesign {
	i := slices.IndexFunc(VesselDesignList, func(w *VesselDesign) bool {
		return w.Name == name
	})
	return VesselDesignList[i]
}

type VesselDesign struct {
	Image resource.ImageID

	Name string

	RotationSpeed gmath.Rad

	MaxHealth float64
	MaxEnergy float64

	EnergyRegen float64

	MaxSpeed     float64
	Acceleration float64

	EnergyResist  float64
	KineticResist float64
	ThermalResist float64

	HitboxSize float64
}

var VesselDesignList = []*VesselDesign{
	{
		Image:         assets.ImageVesselNormal1,
		Name:          "Destroyer",
		RotationSpeed: 2.0,
		MaxSpeed:      150,
		Acceleration:  150,
		MaxHealth:     50,
		MaxEnergy:     50,
		EnergyRegen:   10,
		HitboxSize:    14,
	},
}
