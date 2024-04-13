package scenes

import (
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
	"github.com/quasilyte/ld55-game/styles"
)

type PlayController struct {
	ctx *game.Context
}

func NewPlayController(ctx *game.Context) *PlayController {
	return &PlayController{ctx: ctx}
}

func (c *PlayController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	rows := eui.NewRowContainer(eui.RowContainerConfig{
		MinWidth: 320,
	})

	rows.AddChild(eui.NewCenteredLabel("Play", assets.Font3))

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Continue",
			OnClick: func() {
			},
		})
		b.GetWidget().Disabled = true
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "New Game",
			OnClick: func() {
				c.ctx.Session = game.NewSession()
				c.createDefaultProg()
				game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
			},
		})
		rows.AddChild(b)
	}

	rows.AddChild(eui.NewSeparator(nil, styles.TransparentColor))

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

func (c *PlayController) Update(delta float64) {}

func (c *PlayController) back() {
	game.ChangeScene(c.ctx, NewMainMenuController(c.ctx))
}

func (c *PlayController) createDefaultProg() {
	prog := c.ctx.Session.Prog

	{
		b := &game.ProgBranch{
			Instructions: []game.ProgInstruction{
				game.MakeInst(game.RandomPosInstruction),
				// {Info: game.ProgInstInfoTab[game.]},
			},
		}
		prog.MovementThread.Branches = append(prog.MovementThread.Branches, b)
	}
}
