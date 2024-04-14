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
		ImageIconTargetSpeed:       {Path: "image/inst/target-speed.png"},

		ImageIconSnapShot:   {Path: "image/inst/snap-shot.png"},
		ImageIconNormalShot: {Path: "image/inst/normal-shot.png"},
		ImageIconAimShot:    {Path: "image/inst/aim-shot.png"},

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
		ImageProjectileFreezer: {Path: "image/ammo/freezer_projectile.png"},
		ImageProjectileMissile: {Path: "image/ammo/missile_projectile.png"},

		ImageImpactLaser:   {Path: "image/effects/laser_impact.png", FrameWidth: 10},
		ImageImpactPlasma:  {Path: "image/effects/plasma_impact.png", FrameWidth: 11},
		ImageImpactIon:     {Path: "image/effects/ion_impact.png", FrameWidth: 5},
		ImageImpactLancer:  {Path: "image/effects/lancer_impact.png", FrameWidth: 8},
		ImageImpactPusher:  {Path: "image/effects/pusher_impact.png", FrameWidth: 20},
		ImageImpactFreezer: {Path: "image/effects/freezer_impact.png", FrameWidth: 8},
		ImageImpactMissile: {Path: "image/effects/missile_impact.png", FrameWidth: 11},

		ImageVesselSmall1:  {Path: "image/vessel/small1.png"},
		ImageVesselNormal1: {Path: "image/vessel/normal1.png"},
		ImageVesselLarge1:  {Path: "image/vessel/large1.png"},
		ImageVesselLarge2:  {Path: "image/vessel/large2.png"},
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
	ImageIconTargetSpeed

	ImageItemEnergyShield
	ImageItemKineticShield
	ImageItemThermalShield

	ImageWarning

	ImageIconSnapShot
	ImageIconNormalShot
	ImageIconAimShot

	ImageSpaceBg

	ImageProjectileScatter
	ImageProjectileLaser
	ImageProjectilePlasma
	ImageProjectileIon
	ImageProjectileLancer
	ImageProjectilePusher
	ImageProjectileFreezer
	ImageProjectileMissile

	ImageImpactLaser
	ImageImpactPlasma
	ImageImpactIon
	ImageImpactLancer
	ImageImpactPusher
	ImageImpactFreezer
	ImageImpactMissile

	ImageVesselSmall1
	ImageVesselNormal1
	ImageVesselLarge1
	ImageVesselLarge2
)
