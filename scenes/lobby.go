package scenes

import (
	"fmt"
	"strings"

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
		tt := strings.Join([]string{
			"Configure the hardware.",
			"You can change your equipment here.",
		}, "\n")
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Hardware",
			OnClick: func() {
			},
			Tooltip: eui.NewSimpleTooltip(uiRes, tt),
		})
		rows.AddChild(b)
	}

	{
		tt := strings.Join([]string{
			"Configure the software.",
			"You can program your bot's algorithm here.",
		}, "\n")
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Software",
			OnClick: func() {
				game.ChangeScene(c.ctx, NewSoftwareController(c.ctx))
			},
			Tooltip: eui.NewSimpleTooltip(uiRes, tt),
		})
		rows.AddChild(b)
	}

	{
		tt := strings.Join([]string{
			"Journal contains hints and recon info.",
			"Knowledge is the key!",
		}, "\n")
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Journal",
			OnClick: func() {
			},
			Tooltip: eui.NewSimpleTooltip(uiRes, tt),
		})
		rows.AddChild(b)
	}

	{
		s := fmt.Sprintf("Level %d", c.ctx.Session.Level+1)
		rows.AddChild(eui.NewCenteredLabel(s, assets.Font1))
	}

	{
		tt := strings.Join([]string{
			fmt.Sprintf("Start level %d battle.", c.ctx.Session.Level+1),
			"Make sure that you're prepared!",
			"Consult the Journal for hints.",
		}, "\n")
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Start Battle",
			OnClick: func() {
				game.ChangeScene(c.ctx, NewBattleController(c.ctx))
			},
			Tooltip: eui.NewSimpleTooltip(uiRes, tt),
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
		b.GetWidget().Disabled = true
	}

	root.AddChild(rows)

	initUI(scene, root)
}

func (c *LobbyController) Update(delta float64) {}
