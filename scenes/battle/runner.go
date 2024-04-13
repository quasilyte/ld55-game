package battle

import (
	"math"
	"time"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/battle"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/progsim"
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
		},

		Prog: &game.BotProg{
			Threads: []game.ProgThread{
				{
					Kind: game.MovementThread,
					Branches: []game.ProgBranch{
						{
							Instructions: []game.ProgInstruction{
								// {Info: game.ProgInstInfoTab[game.RandomPosInstruction]},
								// {Info: game.ProgInstInfoTab[game.RotateToInstruction]},
								// {
								// 	Info:   game.ProgInstInfoTab[game.MoveForwardInstruction],
								// 	Params: []any{100.0},
								// },

								{
									Info:   game.ProgInstInfoTab[game.ChanceInstruction],
									Params: []any{0.5},
								},
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

						{
							Instructions: []game.ProgInstruction{
								{Info: game.ProgInstInfoTab[game.CenterPosInstruction]},
								{Info: game.ProgInstInfoTab[game.RotateToInstruction]},
							},
						},
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
}

func (r *Runner) Update(delta float64) {
	for i, e := range r.executors {
		r.vessels[i].SetCommands(e.RunTick(delta))
	}
}
