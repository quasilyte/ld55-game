package scenes

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
	"github.com/quasilyte/ld55-game/styles"
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
		grid := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(len(game.VesselDesignList)+1),
				widget.GridLayoutOpts.Spacing(8, 14),
			)),
		)

		l := eui.NewLabel("Vessel Design", assets.Font1,
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(108, 0)))
		grid.AddChild(l)

		for i := range game.VesselDesignList {
			vd := game.VesselDesignList[i]
			slot := eui.NewSlotButton(uiRes, eui.SlotButtonConfig{})
			slot.Icon.Image = c.ctx.Loader.LoadImage(vd.Image).Data
			grid.AddChild(slot.Container)
		}

		rows.AddChild(grid)
	}

	for i := 0; i < 2; i++ {
		grid := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(len(game.WeaponDesignList)+1),
				widget.GridLayoutOpts.Spacing(8, 14),
			)),
		)

		l := eui.NewLabel(fmt.Sprintf("Weapon %d", i+1), assets.Font1,
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(108, 0)))
		grid.AddChild(l)

		for i := range game.WeaponDesignList {
			wd := game.WeaponDesignList[i]
			slot := eui.NewSlotButton(uiRes, eui.SlotButtonConfig{})
			slot.Icon.Image = c.ctx.Loader.LoadImage(wd.ProjectileImage).Data
			grid.AddChild(slot.Container)
		}

		rows.AddChild(grid)
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

func (c *HardwareController) Update(delta float64) {}

func (c *HardwareController) back() {
	game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
}
