package battle

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/progsim"
	"github.com/quasilyte/ld55-game/styles"
)

type scene = gscene.Scene[ControllerAccessor]

type ControllerAccessor interface {
	GetGameContext() *game.Context
}

type Runner struct {
	ctx *game.Context

	world *game.World

	vessels   []*vesselNode
	executors []*progsim.Executor

	scene *gscene.RootScene[ControllerAccessor]
}

type RunnerConfig struct {
	Context *game.Context

	Scene *gscene.RootScene[ControllerAccessor]
}

func NewRunner(config RunnerConfig) *Runner {
	return &Runner{
		ctx:   config.Context,
		scene: config.Scene,
	}
}

func (r *Runner) Init() {
	session := r.ctx.Session
	level := game.Levels[session.Level]

	spaceBg := r.ctx.Loader.LoadImage(assets.ImageSpaceBg).Data
	for y := 0.0; y < r.ctx.WindowSize.Y; y += float64(spaceBg.Bounds().Dy()) {
		for x := 0.0; x < r.ctx.WindowSize.X; x += float64(spaceBg.Bounds().Dx()) {
			s := r.ctx.NewSprite(assets.ImageSpaceBg)
			s.Pos.Offset = gmath.Vec{X: x, Y: y}
			s.SetCentered(false)
			r.scene.AddGraphics(s)
		}
	}

	playerVessel := &game.Vessel{
		Alliance: 0,
		Pos: gmath.Vec{
			X: level.Dist,
			Y: r.ctx.WindowSize.Y * 0.5,
		},
		Rotation: math.Pi,

		Design:   *r.ctx.Session.VesselDesign,
		Prog:     r.ctx.Session.Prog,
		Artifact: r.ctx.Session.ArtifactDesign,
	}
	for _, wd := range r.ctx.Session.Weapons {
		v := playerVessel
		v.Weapons = append(v.Weapons, &game.Weapon{
			Design: wd,
		})
	}

	enemyVessel := &game.Vessel{
		Alliance: 1,
		Pos: gmath.Vec{
			X: r.ctx.WindowSize.X - level.Dist,
			Y: r.ctx.WindowSize.Y * 0.5,
		},
	}
	{
		enemyData := r.ctx.Loader.LoadRaw(level.Enemy).Data
		var vesselData game.SavedVessel
		if err := json.Unmarshal(enemyData, &vesselData); err != nil {
			panic(err)
		}
		enemyVessel.Design = *game.FindVesselDesignByName(vesselData.VesselDesign)
		enemyVessel.Artifact = game.FindArtifactDesignByName(vesselData.Artifact)
		for _, weaponName := range vesselData.Weapons {
			enemyVessel.Weapons = append(enemyVessel.Weapons, &game.Weapon{
				Design: game.FindWeaponDesignByName(weaponName),
			})
		}
		enemyVessel.Prog = vesselData.Prog
	}

	// 	Pos:      gmath.Vec{X: 256, Y: 256},
	// 	Rotation: 0,

	// 	Design: game.VesselDesign{
	// 		Image:         assets.ImageVesselNormal1,
	// 		RotationSpeed: 2.0,
	// 		MaxSpeed:      150,
	// 		Acceleration:  150,
	// 		EnergyResist:  0.1,
	// 		KineticResist: 0.2,
	// 		ThermalResist: 0.0,
	// 		MaxHealth:     50,
	// 		MaxEnergy:     50,
	// 		HitboxSize:    14,
	// 	},

	// 	Prog: game.NewBotProg(),
	// }

	r.world = &game.World{
		Vessels: []*game.Vessel{
			playerVessel,
			enemyVessel,
		},
		Size: r.ctx.WindowSize,
	}

	r.world.Rand.SetSeed(time.Now().Unix())

	{
		playerVesselNode := newVesselNode(r.world, playerVessel)
		r.vessels = append(r.vessels, playerVesselNode)
		r.scene.AddObject(playerVesselNode)

		e := progsim.NewExecutor(progsim.ExecutorConfig{
			Prog:   playerVessel.Prog,
			World:  r.world,
			Vessel: playerVessel,
		})
		r.executors = append(r.executors, e)
	}
	{
		enemyVesselNode := newVesselNode(r.world, enemyVessel)
		r.vessels = append(r.vessels, enemyVesselNode)
		r.scene.AddObject(enemyVesselNode)

		e := progsim.NewExecutor(progsim.ExecutorConfig{
			Prog:   enemyVessel.Prog,
			World:  r.world,
			Vessel: enemyVessel,
		})
		r.executors = append(r.executors, e)
	}
	{
		r.vessels[0].data.Target = r.world.Vessels[1]
		r.vessels[1].data.Target = r.world.Vessels[0]
	}

	// TODO: move this code somewhere else?
	for i := range r.vessels {
		v := r.vessels[i]
		v.data.Health = v.data.Design.MaxHealth
		v.data.Energy = v.data.Design.MaxEnergy

		v.data.Rotation = v.data.Pos.AngleToPoint(v.world.Size.Mulf(0.5))

		v.data.EnergyResist = v.data.Design.EnergyResist
		v.data.KineticResist = v.data.Design.KineticResist
		v.data.ThermalResist = v.data.Design.ThermalResist

		if v.data.Artifact != nil {
			v.data.Artifact.ApplyBonus(v.data)
		}

		v.data.EventOnDamage.Connect(nil, func(data game.OnDamageData) {
			if data.Attacker == nil {
				return
			}
			clr := styles.AlliedColor
			if v.data.Alliance == 0 {
				clr = styles.EnemyColor
			}
			s := fmt.Sprintf("%.1f", data.Total)
			ft := newFloatingTextNode(v.data.Pos, s, clr)
			r.scene.AddObject(ft)
		})

		v.data.EventDestroyed.Connect(nil, func(attacker *game.Vessel) {
			v.Dispose()
		})
	}
}

func (r *Runner) Update(delta float64) {
	for i, e := range r.executors {
		r.vessels[i].SetCommands(e.RunTick(delta))
	}
}
