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
	c.ctx.Audio().ContinueMusic(assets.AudioMusicMenu)

	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	rows := eui.NewRowContainer(eui.RowContainerConfig{
		MinWidth: 320,
	})

	rows.AddChild(eui.NewCenteredLabel("Astro Heart", assets.Font3))

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Play",
			OnClick: func() {
				game.ChangeScene(c.ctx, NewPlayController(c.ctx))
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
		b.GetWidget().Disabled = true
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
