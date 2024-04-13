package battle

import (
	"fmt"
	"math"
	"time"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/battle"
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

	world *battle.World

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
	spaceBg := r.ctx.Loader.LoadImage(assets.ImageSpaceBg).Data
	for y := 0.0; y < r.ctx.WindowSize.Y; y += float64(spaceBg.Bounds().Dy()) {
		for x := 0.0; x < r.ctx.WindowSize.X; x += float64(spaceBg.Bounds().Dx()) {
			s := r.ctx.NewSprite(assets.ImageSpaceBg)
			s.Pos.Offset = gmath.Vec{X: x, Y: y}
			s.SetCentered(false)
			r.scene.AddGraphics(s)
		}
	}

	playerVessel := &battle.Vessel{
		Alliance: 0,
		Pos:      r.ctx.WindowSize.Mulf(0.5),
		Rotation: math.Pi,

		Design: battle.VesselDesign{
			Image:         assets.ImageVesselNormal1,
			RotationSpeed: 2.0,
			MaxSpeed:      150,
			Acceleration:  150,
			MaxHealth:     50,
			MaxEnergy:     50,
			EnergyRegen:   10,
			HitboxSize:    14,
			Weapons: []*battle.WeaponDesign{
				battle.FindWeaponDesignByName("Pulse Laser"),
				battle.FindWeaponDesignByName("Plasma Cannon"),
			},
		},

		Prog: &game.BotProg{
			Threads: []game.ProgThread{
				{
					Kind: game.Weapon1Thread,
					Branches: []game.ProgBranch{
						{
							Instructions: []game.ProgInstruction{
								{
									Info: game.ProgInstInfoTab[game.TargetPosInstruction],
								},
								{
									Info: game.ProgInstInfoTab[game.DistanceToInstruction],
								},
								{
									Info:   game.ProgInstInfoTab[game.IsLtInstruction],
									Params: []any{150.0},
								},
								{
									Info: game.ProgInstInfoTab[game.SnapShotInstruction],
								},
							},
						},
					},
				},
				{
					Kind: game.Weapon2Thread,
					Branches: []game.ProgBranch{
						{
							Instructions: []game.ProgInstruction{
								{
									Info: game.ProgInstInfoTab[game.TargetPosInstruction],
								},
								{
									Info: game.ProgInstInfoTab[game.DistanceToInstruction],
								},
								{
									Info:   game.ProgInstInfoTab[game.IsLtInstruction],
									Params: []any{200.0},
								},
								{
									Info: game.ProgInstInfoTab[game.NormalShotInstruction],
								},
							},
						},
					},
				},

				{
					Kind: game.MovementThread,
					Branches: []game.ProgBranch{
						{
							Instructions: []game.ProgInstruction{
								// {
								// 	Info:   game.ProgInstInfoTab[game.ChanceInstruction],
								// 	Params: []any{0.5},
								// },
								{Info: game.ProgInstInfoTab[game.TargetPosInstruction]},
								{
									Info:   game.ProgInstInfoTab[game.RandomOffsetInstruction],
									Params: []any{40.0},
								},
								{Info: game.ProgInstInfoTab[game.RotateToInstruction]},
								{
									Info:   game.ProgInstInfoTab[game.MoveForwardInstruction],
									Params: []any{100.0},
								},
							},
						},

						// {
						// 	Instructions: []game.ProgInstruction{
						// 		{Info: game.ProgInstInfoTab[game.CenterPosInstruction]},
						// 		{Info: game.ProgInstInfoTab[game.RotateToInstruction]},
						// 	},
						// },
					},
				},
			},
		},
	}

	enemyVessel := &battle.Vessel{
		Alliance: 1,
		Pos:      gmath.Vec{X: 256, Y: 256},
		Rotation: 0,

		Design: battle.VesselDesign{
			Image:         assets.ImageVesselNormal1,
			RotationSpeed: 2.0,
			MaxSpeed:      150,
			Acceleration:  150,
			EnergyResist:  0.1,
			KineticResist: 0.2,
			ThermalResist: 0.0,
			MaxHealth:     50,
			MaxEnergy:     50,
			HitboxSize:    14,
		},

		Prog: &game.BotProg{},
	}

	r.world = &battle.World{
		Vessels: []*battle.Vessel{
			playerVessel,
			enemyVessel,
		},
		Size: r.ctx.WindowSize,
	}

	r.world.Rand.SetSeed(time.Now().Unix())

	{
		playerVesselNode := newVesselNode(playerVessel)
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
		enemyVesselNode := newVesselNode(enemyVessel)
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

	// TODO: move this code somewhere else.
	for i := range r.vessels {
		v := r.vessels[i]
		v.data.Health = v.data.Design.MaxHealth
		v.data.Energy = v.data.Design.MaxEnergy
		for _, wd := range v.data.Design.Weapons {
			v.data.Weapons = append(v.data.Weapons, &battle.Weapon{
				Design: wd,
			})
		}

		v.data.EventOnDamage.Connect(nil, func(dmg float64) {
			clr := styles.AlliedColor
			if v.data.Alliance == 0 {
				clr = styles.EnemyColor
			}
			s := fmt.Sprintf("%.1f", dmg)
			ft := newFloatingTextNode(v.data.Pos, s, clr)
			r.scene.AddObject(ft)
		})

		v.data.EventDestroyed.Connect(nil, func(attacker *battle.Vessel) {
			v.Dispose()
		})
	}
}

func (r *Runner) Update(delta float64) {
	// r.vessels[0].SetCommands(progsim.VesselCommands{
	// 	FireCommands: []progsim.VesselFireCommand{
	// 		{WeaponIndex: 0, TargetPos: r.vessels[1].data.Pos},
	// 	},
	// })

	for i, e := range r.executors {
		r.vessels[i].SetCommands(e.RunTick(delta))
	}
}
