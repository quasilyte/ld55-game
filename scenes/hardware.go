package scenes

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
	"github.com/quasilyte/ld55-game/styles"
)

type HardwareController struct {
	ctx *game.Context

	rows []*hardwareRow
}

type hardwareRow struct {
	items []*hardwareItem
}

type hardwareItem struct {
	slot  *eui.SlotButton
	value any
}

func NewHardwareController(ctx *game.Context) *HardwareController {
	return &HardwareController{ctx: ctx}
}

func (c *HardwareController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()
	session := c.ctx.Session

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

		row := &hardwareRow{}

		for i := range game.VesselDesignList {
			vd := game.VesselDesignList[i]
			slot := eui.NewSlotButton(uiRes, eui.SlotButtonConfig{
				WithSelector: true,
				OnClick: func() {
					session.VesselDesign = vd
					c.updateSlots()
				},
			})
			slot.Icon.Image = c.ctx.Loader.LoadImage(vd.Image).Data
			grid.AddChild(slot.Container)
			row.items = append(row.items, &hardwareItem{
				slot:  slot,
				value: vd,
			})
		}

		rows.AddChild(grid)

		c.rows = append(c.rows, row)
	}

	for i := 0; i < 2; i++ {
		weaponIndex := i
		row := &hardwareRow{}

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
			slot := eui.NewSlotButton(uiRes, eui.SlotButtonConfig{
				WithSelector: true,
				Tooltip:      eui.NewSimpleTooltip(uiRes, c.weaponDoc(wd)),
				OnClick: func() {
					session.Weapons[weaponIndex] = wd
					c.updateSlots()
				},
			})
			slot.Icon.Image = c.ctx.Loader.LoadImage(wd.ProjectileImage).Data
			grid.AddChild(slot.Container)
			row.items = append(row.items, &hardwareItem{
				slot:  slot,
				value: wd,
			})
		}

		rows.AddChild(grid)

		c.rows = append(c.rows, row)
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

	c.updateSlots()
}

func (c *HardwareController) updateSlots() {
	s := c.ctx.Session

	for rowIndex, row := range c.rows {
		for _, item := range row.items {
			selected := false
			switch rowIndex {
			case 0:
				selected = s.VesselDesign == item.value.(*game.VesselDesign)
			case 1:
				selected = s.Weapons[0] == item.value.(*game.WeaponDesign)
			case 2:
				selected = s.Weapons[1] == item.value.(*game.WeaponDesign)
			}
			if selected {
				item.slot.Selector.GetWidget().Visibility = widget.Visibility_Show
			} else {
				item.slot.Selector.GetWidget().Visibility = widget.Visibility_Hide
			}
		}
	}
}

func (c *HardwareController) Update(delta float64) {}

func (c *HardwareController) back() {
	game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
}

func (c *HardwareController) weaponDoc(wd *game.WeaponDesign) string {
	var lines []string

	switch wd.Name {
	case "Pusher":
		lines = []string{
			"Every hit from this cannon pushes the target.",
			"Remember: vessels take damage when out-of-bounds.",
		}
	case "Scatter Gun":
		lines = []string{
			"A low-tech weapon that doesn't need energy to fire.",
			"It launches 6 projectiles per round, but they",
			"deal only minor Kinetic damage.",
		}
	case "Pulse Laser":
		lines = []string{
			"A well-balanced energy weapon.",
		}
	case "Ion Cannon":
		lines = []string{
			"A long-range ion cannon.",
			"Every hit burns some of the target's energy.",
		}
	case "Plasma Cannon":
		lines = []string{
			"Deals more damage than a Pulse Laser.",
			"The primary damage comes from Thermal effect.",
			"Projectiles move slowly.",
		}
	case "Lancer":
		lines = []string{
			"The heaviest-hitting laser you can find.",
			"It fires in a straight direction only.",
		}
	}

	resultLines := []string{
		wd.Name,
		"",
	}
	resultLines = append(resultLines, lines...)
	resultLines = append(resultLines, "")

	totalDamage := wd.Damage.Energy + wd.Damage.Kinetic + wd.Damage.Thermal
	dps := totalDamage * (1.0 / wd.Reload)
	multiplier := ""
	if wd.Burst != 1 {
		multiplier = fmt.Sprintf("*%d", wd.Burst)
	}
	resultLines = append(resultLines, fmt.Sprintf("DPS: %.1f%s (E%d K%d T%d)",
		dps, multiplier, int(wd.Damage.Energy), int(wd.Damage.Kinetic), int(wd.Damage.Thermal)))

	resultLines = append(resultLines, fmt.Sprintf("Max range: %d", int(wd.MaxRange)))
	resultLines = append(resultLines, fmt.Sprintf("Energy cost: %d (per shot)", int(wd.EnergyCost)))
	switch wd.FiringType {
	case game.TargetableWeapon:
		resultLines = append(resultLines, "Targeting: fires at targeted vessel")
	case game.FixedAngleWeapon:
		resultLines = append(resultLines, "Targeting: fixed frontal attack")
	}

	return strings.Join(resultLines, "\n")
}
