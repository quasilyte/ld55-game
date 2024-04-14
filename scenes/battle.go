package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/scenes/battle"
)

type BattleController struct {
	ctx *game.Context

	pause       bool
	exitPending bool
	exitDelay   float64

	runner *battle.Runner
}

func NewBattleController(ctx *game.Context) *BattleController {
	return &BattleController{ctx: ctx}
}

func (c *BattleController) Init(scene *gscene.RootScene[battle.ControllerAccessor]) {
	c.ctx.CRT = true
	c.ctx.Audio().PauseCurrentMusic()

	c.runner = battle.NewRunner(battle.RunnerConfig{
		Context: c.ctx,
		Scene:   scene,
	})
	c.runner.Init()
	c.runner.EventBattleOver.Connect(nil, func(victory bool) {
		c.finishBattle(victory)
	})
}

func (c *BattleController) GetGameContext() *game.Context {
	return c.ctx
}

func (c *BattleController) Update(delta float64) {
	if !c.pause {
		c.runner.Update(delta)
	}

	if c.exitPending {
		c.exitDelay -= delta
		if c.exitDelay <= 0 {
			if c.ctx.Session.Level+1 > len(game.Levels) {
				// TODO: go to credits?
				game.ChangeScene(c.ctx, NewMainMenuController(c.ctx))
				return
			}
			game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
			return
		}
	}

	if !c.exitPending && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		c.finishBattle(true)
	}
}

func (c *BattleController) finishBattle(victory bool) {
	if c.exitPending {
		return
	}

	c.exitPending = true
	c.exitDelay = 2

	if victory {
		c.ctx.Session.Level++
	}
}
