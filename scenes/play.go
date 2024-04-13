package scenes

import (
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/game"
)

type PlayController struct {
	ctx *game.Context
}

func NewPlayController(ctx *game.Context) *PlayController {
	return &PlayController{ctx: ctx}
}

func (c *PlayController) Init(scene *gscene.SimpleRootScene) {

}
