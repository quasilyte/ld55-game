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

type SoftwareController struct {
	ctx *game.Context

	selectedTab *softwareTab
	tabs        []*softwareTab
}

type softwareTab struct {
	index  int
	button *widget.Button
}

func NewSoftwareController(ctx *game.Context) *SoftwareController {
	return &SoftwareController{ctx: ctx}
}

func (c *SoftwareController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	rows := eui.NewRowContainer(eui.RowContainerConfig{})

	{
		sysTabs := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(4),
			)),
		)

		sysTabs.AddChild(
			eui.NewLabel("SYS ", assets.Font2,
				widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}))))

		movButton := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "MOV",
			Tooltip: eui.NewTooltip(uiRes, strings.Join([]string{
				"Configure [Movement] beharior.",
				"This subprogram controls the vessel's movement.",
				"",
				"Click to select.",
			}, "\n")),
			OnClick: func() {
				c.selectTab(0)
			},
		})
		c.tabs = append(c.tabs, &softwareTab{
			index:  0,
			button: movButton,
		})
		c.selectedTab = c.tabs[0]
		sysTabs.AddChild(movButton)
		movButton.TextColor.Idle = styles.SelectedTextColor

		w1button := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "WP1",
			Tooltip: eui.NewTooltip(uiRes, strings.Join([]string{
				"Configure [Weapon 1] beharior.",
				"This subprogram controls the weapon usage.",
				"",
				"Click to select.",
			}, "\n")),
			OnClick: func() {
				c.selectTab(1)
			},
		})
		c.tabs = append(c.tabs, &softwareTab{
			index:  1,
			button: w1button,
		})
		sysTabs.AddChild(w1button)

		w2button := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "WP2",
			Tooltip: eui.NewTooltip(uiRes, strings.Join([]string{
				"Configure [Weapon 2] beharior.",
				"This subprogram controls the weapon usage.",
				"",
				"Click to select.",
			}, "\n")),
			OnClick: func() {
				c.selectTab(2)
			},
		})
		c.tabs = append(c.tabs, &softwareTab{
			index:  2,
			button: w2button,
		})
		sysTabs.AddChild(w2button)

		defButton := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "DEF",
			Tooltip: eui.NewTooltip(uiRes, strings.Join([]string{
				"Configure [Defense] beharior.",
				"This subprogram controls the shield/dodge behavior.",
				"",
				"Click to select.",
			}, "\n")),
			OnClick: func() {
				c.selectTab(3)
			},
		})
		c.tabs = append(c.tabs, &softwareTab{
			index:  3,
			button: defButton,
		})
		sysTabs.AddChild(defButton)

		sysTabs.AddChild(
			eui.NewLabel("  |  ", assets.Font2,
				widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}))))

		saveButton := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Save",
			Tooltip: eui.NewTooltip(uiRes, strings.Join([]string{
				"Save edits and go back.",
			}, "\n")),
		})
		sysTabs.AddChild(saveButton)

		rows.AddChild(sysTabs)
	}

	rows.AddChild(eui.NewSeparator(nil, styles.DisabledTextColor))

	{
		numCols := 8
		numRows := 3
		grid := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(numCols+1+(numCols-1)),
				widget.GridLayoutOpts.Spacing(0, 14),
			)),
		)
		for row := 0; row < numRows; row++ {
			grid.AddChild(eui.NewLabel(fmt.Sprintf("Branch %d", row+1), assets.Font1))
			for col := 0; col < numCols; col++ {
				grid.AddChild(eui.NewButton(uiRes, eui.ButtonConfig{Slot: true}))
				if col != numCols-1 {
					grid.AddChild(eui.NewCenteredLabel(">", assets.Font1))
				}
			}
		}
		rows.AddChild(grid)
	}

	rows.AddChild(eui.NewSeparator(nil, styles.DisabledTextColor))

	root.AddChild(rows)

	initUI(scene, root)
}

func (c *SoftwareController) selectTab(index int) {
	for _, t := range c.tabs {
		selected := t.index == index
		if selected {
			c.selectedTab = t
			t.button.TextColor.Idle = styles.SelectedTextColor
		} else {
			t.button.TextColor.Idle = styles.NormalTextColor
		}
	}
}

func (c *SoftwareController) Update(delta float64) {}
