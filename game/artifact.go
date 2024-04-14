package game

import (
	"slices"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ld55-game/assets"
)

func FindArtifactDesignByName(name string) *ArtifactDesign {
	i := slices.IndexFunc(ArtifactDesignList, func(w *ArtifactDesign) bool {
		return w.Name == name
	})
	return ArtifactDesignList[i]
}

type ArtifactDesign struct {
	Name string

	Icon resource.ImageID

	ApplyBonus func(v *Vessel)
}

var ArtifactDesignList = []*ArtifactDesign{
	{
		Name: "E-Shield",
		Icon: assets.ImageItemEnergyShield,
		ApplyBonus: func(v *Vessel) {
			v.EnergyResist += 0.15
		},
	},
	{
		Name: "K-Shield",
		Icon: assets.ImageItemKineticShield,
		ApplyBonus: func(v *Vessel) {
			v.KineticResist += 0.25
		},
	},
	{
		Name: "T-Shield",
		Icon: assets.ImageItemThermalShield,
		ApplyBonus: func(v *Vessel) {
			v.ThermalResist += 0.2
		},
	},
}
