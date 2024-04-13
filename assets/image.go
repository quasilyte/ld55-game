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

		ImageSpaceBg: {Path: "image/space_bg.png"},

		ImageProjectileLaser: {Path: "image/ammo/laser_projectile.png"},

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

	ImageSpaceBg

	ImageProjectileLaser

	ImageVesselNormal1
)
