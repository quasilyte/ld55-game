package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"

	_ "image/png"
)

func registerImageResources(loader *resource.Loader) {
	resources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonIdle:     {Path: "image/button-idle.png"},
		ImageUIButtonHover:    {Path: "image/button-hover.png"},
		ImageUIButtonPressed:  {Path: "image/button-pressed.png"},
		ImageUIButtonDisabled: {Path: "image/button-disabled.png"},
		ImageUITooltip:        {Path: "image/tooltip.png"},
		ImageUISlotIdle:       {Path: "image/slot-idle.png"},
		ImageUISlotHover:      {Path: "image/slot-hover.png"},
		ImageUISlotDisabled:   {Path: "image/slot-disabled.png"},

		ImageSpaceBg: {Path: "image/space_bg.png"},

		ImageProjectileLaser:  {Path: "image/ammo/laser_projectile.png"},
		ImageProjectilePlasma: {Path: "image/ammo/plasma_projectile.png"},

		ImageImpactLaser:  {Path: "image/effects/laser_impact.png", FrameWidth: 10},
		ImageImpactPlasma: {Path: "image/effects/plasma_impact.png", FrameWidth: 11},

		ImageVesselNormal1: {Path: "image/vessel/normal1.png"},
	}

	for id, info := range resources {
		loader.ImageRegistry.Set(id, info)
		loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUIButtonDisabled
	ImageUITooltip
	ImageUISlotIdle
	ImageUISlotHover
	ImageUISlotDisabled

	ImageSpaceBg

	ImageProjectileLaser
	ImageProjectilePlasma

	ImageImpactLaser
	ImageImpactPlasma

	ImageVesselNormal1
)
