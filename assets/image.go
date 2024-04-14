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

		ImageIconRandomPos:         {Path: "image/inst/random-pos.png"},
		ImageIconRotateTo:          {Path: "image/inst/rotate-to.png"},
		ImageIconMoveForward:       {Path: "image/inst/move-forward.png"},
		ImageIconMoveAndRotate:     {Path: "image/inst/move-and-rotate.png"},
		ImageIconDistanceTo:        {Path: "image/inst/distance-to.png"},
		ImageIconTargetPos:         {Path: "image/inst/target-pos.png"},
		ImageIconSelfPos:           {Path: "image/inst/self-pos.png"},
		ImageIconCenterPos:         {Path: "image/inst/center-pos.png"},
		ImageIconIsLt:              {Path: "image/inst/is-lt.png"},
		ImageIconIsGt:              {Path: "image/inst/is-gt.png"},
		ImageIconRandomOffset:      {Path: "image/inst/random-offset.png"},
		ImageIconRand:              {Path: "image/inst/rand.png"},
		ImageIconSelfHealthPercent: {Path: "image/inst/self-hp-percent.png"},
		ImageIconSelfEnergyPercent: {Path: "image/inst/self-energy-percent.png"},

		ImageIconSnapShot:   {Path: "image/inst/snap-shot.png"},
		ImageIconNormalShot: {Path: "image/inst/normal-shot.png"},

		ImageSpaceBg: {Path: "image/space_bg.png"},

		ImageWarning: {Path: "image/warning.png"},

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

	ImageIconRandomPos
	ImageIconRotateTo
	ImageIconMoveForward
	ImageIconMoveAndRotate
	ImageIconDistanceTo
	ImageIconTargetPos
	ImageIconCenterPos
	ImageIconSelfPos
	ImageIconIsLt
	ImageIconIsGt
	ImageIconRandomOffset
	ImageIconRand
	ImageIconSelfHealthPercent
	ImageIconSelfEnergyPercent

	ImageWarning

	ImageIconSnapShot
	ImageIconNormalShot

	ImageSpaceBg

	ImageProjectileLaser
	ImageProjectilePlasma

	ImageImpactLaser
	ImageImpactPlasma

	ImageVesselNormal1
)
