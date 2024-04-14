package scenes

import (
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
)

type HardwareController struct {
	ctx *game.Context
}

func NewHardwareController(ctx *game.Context) *HardwareController {
	return &HardwareController{ctx: ctx}
}

func (c *HardwareController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	rows := eui.NewRowContainer(eui.RowContainerConfig{})

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Back",
			OnClick: func() {
				c.back()
			},
		})
		rows.AddChild(b)
	}

	root.AddChild(rows)

	initUI(scene, root)
}

func (c *HardwareController) Update(delta float64) {}

func (c *HardwareController) back() {
	game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
}
