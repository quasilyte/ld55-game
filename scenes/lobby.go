package scenes

import (
	"fmt"

	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
	"github.com/quasilyte/ld55-game/styles"
)

type LobbyController struct {
	ctx *game.Context
}

func NewLobbyController(ctx *game.Context) *LobbyController {
	return &LobbyController{ctx: ctx}
}

func (c *LobbyController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	rows := eui.NewRowContainer(eui.RowContainerConfig{
		MinWidth: 320,
	})

	rows.AddChild(eui.NewCenteredLabel("Hangar", assets.Font3))

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Hardware",
			OnClick: func() {
			},
		})
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Software",
			OnClick: func() {
			},
		})
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Journal",
			OnClick: func() {
			},
		})
		b.GetWidget().Disabled = true
		rows.AddChild(b)
	}

	{
		s := fmt.Sprintf("Level %d", c.ctx.Session.Level+1)
		rows.AddChild(eui.NewCenteredLabel(s, assets.Font1))
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Start Battle",
			OnClick: func() {
				game.ChangeScene(c.ctx, NewBattleController(c.ctx))
			},
		})
		rows.AddChild(b)
	}

	rows.AddChild(eui.NewSeparator(nil, styles.TransparentColor))

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Save and Exit",
			OnClick: func() {
			},
		})
		rows.AddChild(b)
	}

	root.AddChild(rows)

	initUI(scene, root)
}

func (c *LobbyController) Update(delta float64) {}
