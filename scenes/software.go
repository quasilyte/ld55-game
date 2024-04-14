package scenes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/game"
	"github.com/quasilyte/ld55-game/gameui/eui"
	"github.com/quasilyte/ld55-game/styles"
)

type SoftwareController struct {
	ctx *game.Context

	prog *game.BotProg

	hoverSlot   *softwareSlot
	selectedTab *softwareTab
	tabs        []*softwareTab

	slots [3][]*softwareSlot
}

type softwareTab struct {
	index  int
	button *widget.Button
	thread *game.ProgThread
}

type softwareSlot struct {
	branchIndex      int
	instructionIndex int
	tooltipText      *widget.Text
	button           *eui.SlotButton
}

func NewSoftwareController(ctx *game.Context) *SoftwareController {
	return &SoftwareController{ctx: ctx}
}

func (c *SoftwareController) Init(scene *gscene.SimpleRootScene) {
	uiRes := c.ctx.UIResources
	root := eui.NewRootContainer()

	c.prog = c.ctx.Session.Prog

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
			Tooltip: eui.NewSimpleTooltip(uiRes, strings.Join([]string{
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
			thread: c.prog.MovementThread,
		})
		c.selectedTab = c.tabs[0]
		sysTabs.AddChild(movButton)
		movButton.TextColor.Idle = styles.SelectedTextColor

		w1button := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "WP1",
			Tooltip: eui.NewSimpleTooltip(uiRes, strings.Join([]string{
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
			thread: c.prog.Weapon1Thread,
		})
		sysTabs.AddChild(w1button)

		w2button := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "WP2",
			Tooltip: eui.NewSimpleTooltip(uiRes, strings.Join([]string{
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
			thread: c.prog.Weapon2Thread,
		})
		sysTabs.AddChild(w2button)

		defButton := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "DEF",
			Tooltip: eui.NewSimpleTooltip(uiRes, strings.Join([]string{
				"Configure [Defense] beharior.",
				"This subprogram controls the shield/dodge behavior.",
				"",
				"Click to select.",
			}, "\n")),
			OnClick: func() {
				c.selectTab(3)
			},
		})
		defButton.GetWidget().Disabled = true
		c.tabs = append(c.tabs, &softwareTab{
			index:  3,
			button: defButton,
			thread: c.prog.DefThread,
		})
		sysTabs.AddChild(defButton)

		sysTabs.AddChild(
			eui.NewLabel("  |  ", assets.Font2,
				widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}))))

		testButton := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Test",
			Tooltip: eui.NewSimpleTooltip(uiRes, strings.Join([]string{
				"Do a quick test of your program.",
			}, "\n")),
		})
		testButton.GetWidget().Disabled = true
		sysTabs.AddChild(testButton)

		saveButton := eui.NewButton(uiRes, eui.ButtonConfig{
			Text: "Save",
			Tooltip: eui.NewSimpleTooltip(uiRes, strings.Join([]string{
				"Save edits and go back.",
			}, "\n")),
			OnClick: func() {
				game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
			},
		})
		sysTabs.AddChild(saveButton)

		rows.AddChild(sysTabs)
	}

	rows.AddChild(eui.NewSeparator(nil, styles.DisabledTextColor))

	{
		numCols := 10
		numRows := 3
		grid := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(numCols+1+(numCols-1)),
				widget.GridLayoutOpts.Spacing(0, 14),
			)),
		)
		for row := 0; row < numRows; row++ {
			grid.AddChild(eui.NewLabel(fmt.Sprintf("Branch %d ", row+1), assets.Font1))
			for col := 0; col < numCols; col++ {
				tt := eui.NewTooltip(uiRes, "")
				slot := &softwareSlot{}
				slotButton := eui.NewSlotButton(uiRes, eui.SlotButtonConfig{
					Tooltip:   tt.Container,
					WithLabel: true,
					OnHoverStart: func(sb *eui.SlotButton) {
						if sb.Label == nil || sb.Label.Label == "" {
							return
						}
						sb.Label.Color = styles.AlliedColor
						c.hoverSlot = slot
					},
					OnHoverEnd: func(sb *eui.SlotButton) {
						if sb.Label == nil || sb.Label.Label == "" {
							return
						}
						sb.Label.Color = styles.NormalTextColor
						c.hoverSlot = nil
					},
				})
				slot.branchIndex = row
				slot.instructionIndex = col
				slot.button = slotButton
				slot.tooltipText = tt.Text
				c.slots[row] = append(c.slots[row], slot)
				grid.AddChild(slotButton.Container)
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

	c.updateInstructionSlots()
}

func (c *SoftwareController) updateInstructionSlots() {
	thread := c.selectedTab.thread

	for i, row := range c.slots {
		var branch *game.ProgBranch
		if i < len(thread.Branches) {
			branch = thread.Branches[i]
		}
		for j, b := range row {
			if branch == nil || j >= len(branch.Instructions) {
				b.button.Icon.Image = nil
				b.button.Label.Label = ""
				b.tooltipText.Label = c.instDoc(game.MakeInst(game.NopInstruction, 0))
				continue
			}
			inst := branch.Instructions[j]
			b.button.Icon.Image = c.ctx.Loader.LoadImage(inst.Info.Icon).Data
			b.button.Label.Label = c.formatInstParam(inst)
			b.tooltipText.Label = c.instDoc(inst)
		}
	}
}

func (c *SoftwareController) getSlotInst(slot *softwareSlot) *game.ProgInstruction {
	return &c.selectedTab.thread.Branches[slot.branchIndex].Instructions[slot.instructionIndex]
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

	c.updateInstructionSlots()
}

func (c *SoftwareController) formatInstParam(inst game.ProgInstruction) string {
	if !inst.Info.Param {
		return ""
	}

	switch inst.Info.Kind {
	case game.RandomOffsetInstruction, game.IsLtInstruction, game.IsGtInstruction:
		return strconv.Itoa(int(inst.Param))
	case game.ChanceInstruction:
		return strconv.Itoa(int(inst.Param)) + "%"
	case game.MoveForwardInstruction, game.MoveAndRotateInstruction:
		return strconv.Itoa(int(inst.Param))
	default:
		return ""
	}
}

func (c *SoftwareController) instDoc(inst game.ProgInstruction) string {
	var lines []string

	switch inst.Info.Kind {
	case game.NopInstruction:
		lines = []string{
			"An empty instruction does nothing.",
			"It's also known as NOP instruction",
		}

	case game.RandomPosInstruction:
		lines = []string{
			"Push a random pos to the stack.",
		}

	case game.RotateToInstruction:
		lines = []string{
			"Rotate to the destination point.",
			"Keeps the engines offline while rotating.",
			"The destination point will be popped from the stack.",
		}
	case game.MoveForwardInstruction:
		lines = []string{
			"Turns on the engine.",
			"Moves forward for the specified amount of units.",
		}

	case game.TargetPosInstruction:
		lines = []string{
			"Push a target pos to the stack.",
		}

	case game.DistanceToInstruction:
		lines = []string{
			"Pop a stack value and push back a",
			"distance between that pos and the current",
			"possition of the vessel.",
		}

	case game.IsLtInstruction:
		lines = []string{
			"Pop a stack value and check",
			"whether it's less than the value specified.",
			"If not, go to the next branch.",
		}
	case game.IsGtInstruction:
		lines = []string{
			"Pop a top stack value and check",
			"whether it's greater than the value specified.",
			"If not, go to the next branch.",
		}

	case game.SnapShotInstruction:
		lines = []string{
			"Fires weapon at the target.",
			"Snap shot allows a faster rate-of-fire with",
			"a very small accuracy.",
			"It can be good against unpredictable targets.",
		}
	}

	if inst.Info.Param {
		lines = append(lines, "", "Hover and start typing to change the value.")
	}

	return strings.Join(lines, "\n")
}

func (c *SoftwareController) Update(delta float64) {
	if c.maybeEditValue() {
		c.updateInstructionSlots()
	}
}

func (c *SoftwareController) maybeEditValue() bool {
	if c.hoverSlot == nil {
		return false
	}

	inst := c.getSlotInst(c.hoverSlot)

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		s := strconv.Itoa(int(inst.Param))
		if len(s) >= 2 {
			v, _ := strconv.ParseFloat(s[:len(s)-1], 64)
			inst.SetParam(v)
		} else {
			inst.SetParam(0)
		}
		return true
	}

	for k := ebiten.KeyDigit0; k <= ebiten.KeyDigit9; k++ {
		if !inpututil.IsKeyJustPressed(k) {
			continue
		}
		kv := k - ebiten.KeyDigit0
		s := strconv.Itoa(int(inst.Param))
		s += strconv.Itoa(int(kv))
		v, _ := strconv.ParseFloat(s, 64)
		inst.SetParam(v)
		return true
	}

	return false
}
