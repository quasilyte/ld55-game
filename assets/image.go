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

		ImageSlotSelector: {Path: "image/slot-selector.png"},

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
		ImageIconIsOutOfBounds:     {Path: "image/inst/is-out-of-bounds.png"},
		ImageIconRandomOffset:      {Path: "image/inst/random-offset.png"},
		ImageIconRand:              {Path: "image/inst/rand.png"},
		ImageIconSelfHealthPercent: {Path: "image/inst/self-hp-percent.png"},
		ImageIconSelfEnergyPercent: {Path: "image/inst/self-energy-percent.png"},
		ImageIconWait:              {Path: "image/inst/wait.png"},

		ImageIconSnapShot:   {Path: "image/inst/snap-shot.png"},
		ImageIconNormalShot: {Path: "image/inst/normal-shot.png"},

		ImageSpaceBg: {Path: "image/space_bg.png"},

		ImageItemEnergyShield:  {Path: "image/item/energy-shield.png"},
		ImageItemKineticShield: {Path: "image/item/kinetic-shield.png"},
		ImageItemThermalShield: {Path: "image/item/thermal-shield.png"},

		ImageWarning: {Path: "image/warning.png"},

		ImageProjectileScatter: {Path: "image/ammo/scatter_projectile.png"},
		ImageProjectileLaser:   {Path: "image/ammo/laser_projectile.png"},
		ImageProjectilePlasma:  {Path: "image/ammo/plasma_projectile.png"},
		ImageProjectileIon:     {Path: "image/ammo/ion_projectile.png"},
		ImageProjectileLancer:  {Path: "image/ammo/lancer_projectile.png"},
		ImageProjectilePusher:  {Path: "image/ammo/pusher_projectile.png"},

		ImageImpactLaser:  {Path: "image/effects/laser_impact.png", FrameWidth: 10},
		ImageImpactPlasma: {Path: "image/effects/plasma_impact.png", FrameWidth: 11},
		ImageImpactIon:    {Path: "image/effects/ion_impact.png", FrameWidth: 5},
		ImageImpactLancer: {Path: "image/effects/lancer_impact.png", FrameWidth: 8},
		ImageImpactPusher: {Path: "image/effects/pusher_impact.png", FrameWidth: 20},

		ImageVesselSmall1:  {Path: "image/vessel/small1.png"},
		ImageVesselNormal1: {Path: "image/vessel/normal1.png"},
		ImageVesselLarge1:  {Path: "image/vessel/large1.png"},
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

	ImageSlotSelector

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
	ImageIconIsOutOfBounds
	ImageIconRandomOffset
	ImageIconRand
	ImageIconSelfHealthPercent
	ImageIconSelfEnergyPercent
	ImageIconWait

	ImageItemEnergyShield
	ImageItemKineticShield
	ImageItemThermalShield

	ImageWarning

	ImageIconSnapShot
	ImageIconNormalShot

	ImageSpaceBg

	ImageProjectileScatter
	ImageProjectileLaser
	ImageProjectilePlasma
	ImageProjectileIon
	ImageProjectileLancer
	ImageProjectilePusher

	ImageImpactLaser
	ImageImpactPlasma
	ImageImpactIon
	ImageImpactLancer
	ImageImpactPusher

	ImageVesselSmall1
	ImageVesselNormal1
	ImageVesselLarge1
)
