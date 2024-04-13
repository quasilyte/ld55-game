package scenes

import (
	"os"

	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
)

type MainMenuController struct {
	ctx *game.Context
}

func NewMainMenuController(ctx *game.Context) *MainMenuController {
	return &MainMenuController{ctx: ctx}
}

func (c *MainMenuController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	rows := eui.NewRowContainer(eui.RowContainerConfig{})

	rows.AddChild(eui.NewCenteredLabel("AstroHeart", assets.Font3))

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Play",
			OnClick: func() {
				game.ChangeScene(c.ctx, NewBattleController(c.ctx))
			},
		})
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Settings",
			OnClick: func() {
			},
		})
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Credits",
			OnClick: func() {
			},
		})
		b.GetWidget().Disabled = true
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Exit",
			OnClick: func() {
				os.Exit(0)
			},
		})
		rows.AddChild(b)
	}

	root.AddChild(rows)

	initUI(scene, root)
}

func (c *MainMenuController) Update(delta float64) {}
