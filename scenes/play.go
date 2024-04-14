package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
				game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
			},
		})
		b.GetWidget().Disabled = c.ctx.Session == nil
		rows.AddChild(b)
	}

	{
		b := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "New Game",
			OnClick: func() {
				c.ctx.Session = game.NewSession()
				c.createDefaultProg()
				c.createDefaultVesselDesign()
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

func (c *PlayController) Update(delta float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.ChangeScene(c.ctx, NewMainMenuController(c.ctx))
		return
	}
}

func (c *PlayController) back() {
	game.ChangeScene(c.ctx, NewMainMenuController(c.ctx))
}

func (c *PlayController) createDefaultVesselDesign() {
	s := c.ctx.Session

	s.Weapons = []*game.WeaponDesign{
		game.FindWeaponDesignByName("Pulse Laser"),
		game.FindWeaponDesignByName("Scatter Gun"),
	}

	s.VesselDesign = game.FindVesselDesignByName("Fighter")

	s.ArtifactDesign = game.FindArtifactDesignByName("E-Shield")
}

func (c *PlayController) createDefaultProg() {
	prog := c.ctx.Session.Prog

	{
		b := &game.ProgBranch{
			Instructions: []game.ProgInstruction{
				game.MakeInst(game.RandomPosInstruction, 0),
				game.MakeInst(game.RotateToInstruction, 0),
				game.MakeInst(game.MoveForwardInstruction, 100),
			},
		}
		prog.MovementThread.Branches = append(prog.MovementThread.Branches, b)
	}

	{
		b := &game.ProgBranch{
			Instructions: []game.ProgInstruction{
				game.MakeInst(game.TargetPosInstruction, 0),
				game.MakeInst(game.DistanceToInstruction, 0),
				game.MakeInst(game.IsLtInstruction, 150),
				game.MakeInst(game.SnapShotInstruction, 0),
			},
		}
		prog.Weapon1Thread.Branches = append(prog.Weapon1Thread.Branches, b)
	}

	// For convenience, pad everything with NOPs.
	prog.EachThread(func(i int, t *game.ProgThread) {
		for len(t.Branches) < game.MaxBranches {
			t.Branches = append(t.Branches, &game.ProgBranch{})
		}
		for _, b := range t.Branches {
			for len(b.Instructions) < game.MaxInstructions {
				b.Instructions = append(b.Instructions, game.MakeInst(game.NopInstruction, 0))
			}
		}
	})
}
