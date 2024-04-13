package scenes

import (
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/scenes/battle"
)

type BattleController struct {
	ctx *game.Context

	pause bool

	runner *battle.Runner
}

func NewBattleController(ctx *game.Context) *BattleController {
	return &BattleController{ctx: ctx}
}

func (c *BattleController) Init(scene *gscene.RootScene[battle.ControllerAccessor]) {
	c.runner = battle.NewRunner(battle.RunnerConfig{
		Context: c.ctx,
		Scene:   scene,
	})
	c.runner.Init()
}

func (c *BattleController) GetGameContext() *game.Context {
	return c.ctx
}

func (c *BattleController) Update(delta float64) {
	if !c.pause {
		c.runner.Update(delta)
	}
}
