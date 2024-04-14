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

	Slow float64

	Prog *BotProg

	Target *Vessel

	Weapons  []*Weapon
	Artifact *ArtifactDesign

	EnergyResist  float64
	KineticResist float64
	ThermalResist float64

	Design VesselDesign

	EventOnDamage  gsignal.Event[OnDamageData]
	EventDestroyed gsignal.Event[*Vessel] // Vessel arg is attacker
}

type OnDamageData struct {
	Total    float64
	Attacker *Vessel
}

func (v *Vessel) Velocity() gmath.Vec {
	return v.EngineVelocity.Add(v.ExtraVelocity)
}

func (v *Vessel) OnDamage(d Damage, attacker *Vessel) {
	totalDamage := ((1.0 - v.EnergyResist) * d.Energy) +
		((1.0 - v.KineticResist) * d.Kinetic) +
		((1.0 - v.ThermalResist) * d.Thermal)
	if totalDamage > 0 {
		v.EventOnDamage.Emit(OnDamageData{
			Total:    totalDamage,
			Attacker: attacker,
		})
		v.Health = gmath.ClampMin(v.Health-totalDamage, 0)
		if v.Health == 0 {
			v.EventDestroyed.Emit(attacker)
		}
	}

	if d.Energy > 0 {
		v.Energy = gmath.ClampMin(v.Energy-d.DrainEnergy, 0)
	}

	v.Slow = gmath.ClampMax(v.Slow+d.Slow, 20)
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
		Image:         assets.ImageVesselSmall1,
		Name:          "Machpella",
		RotationSpeed: 2,
		MaxSpeed:      140,
		Acceleration:  160,
		MaxHealth:     30,
		MaxEnergy:     70,
		EnergyRegen:   6,
		HitboxSize:    4,

		KineticResist: 0.1,
		EnergyResist:  0.1,
		ThermalResist: 0.1,
	},

	{
		Image:         assets.ImageVesselNormal1,
		Name:          "Fighter",
		RotationSpeed: 2.5,
		MaxSpeed:      110,
		Acceleration:  150,
		MaxHealth:     40,
		MaxEnergy:     50,
		EnergyRegen:   4,
		HitboxSize:    12,

		KineticResist: 0.3,
		EnergyResist:  0.2,
		ThermalResist: 0,
	},

	{
		Image:         assets.ImageVesselLarge1,
		Name:          "Destroyer",
		RotationSpeed: 1.6,
		MaxSpeed:      90,
		Acceleration:  200,
		MaxHealth:     70,
		MaxEnergy:     40,
		EnergyRegen:   3,
		HitboxSize:    20,

		KineticResist: 0.1,
		EnergyResist:  0.5,
		ThermalResist: 0.3,
	},

	// An improved Destroyed as a final boss.
	{
		Image:         assets.ImageVesselLarge2,
		Name:          "Boss",
		RotationSpeed: 1.8,
		MaxSpeed:      90,
		Acceleration:  150,
		MaxHealth:     120,
		MaxEnergy:     150,
		EnergyRegen:   6,
		HitboxSize:    18,

		KineticResist: 0.2,
		EnergyResist:  0.5,
		ThermalResist: 0.5,
	},
}
