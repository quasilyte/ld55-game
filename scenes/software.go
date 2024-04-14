package scenes

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
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

	dragFrom gmath.Vec
	dragTo   gmath.Vec
	dragLine *graphics.Line

	draggingInst  *softwareInstSlot
	hoverInstSlot *softwareInstSlot
	instList      []*softwareInstSlot

	hasErrors bool
}

type softwareTab struct {
	index  int
	button *widget.Button
	thread *game.ProgThread
}

type softwareInstSlot struct {
	inst        game.ProgInstruction
	tooltipText *widget.Text
	button      *eui.SlotButton
}

type softwareSlot struct {
	branchIndex      int
	instructionIndex int
	warning          *graphics.Sprite
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
				if !c.hasErrors {
					game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
				} else {
					c.ctx.Audio().PlaySound(assets.AudioErrorBeep)
				}
			},
		})
		sysTabs.AddChild(saveButton)

		rows.AddChild(sysTabs)
	}

	rows.AddChild(eui.NewSeparator(nil, styles.DisabledTextColor))

	{
		numCols := game.MaxInstructions
		numRows := game.MaxBranches
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
						sb.Label.Color = styles.AlliedColor
						c.hoverSlot = slot
					},
					OnHoverEnd: func(sb *eui.SlotButton) {
						sb.Label.Color = styles.NormalTextColor
						if c.hoverSlot == slot {
							c.hoverSlot = nil
						}
					},
				})

				warning := c.ctx.NewSprite(assets.ImageWarning)
				warning.SetVisibility(false)

				slot.branchIndex = row
				slot.instructionIndex = col
				slot.button = slotButton
				slot.tooltipText = tt.Text
				slot.warning = warning
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

	{
		numCols := game.MaxInstructions + 2
		numRows := 2
		grid := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(numCols),
				widget.GridLayoutOpts.Spacing(0, 14),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
				}),
			),
		)
		for row := 0; row < numRows; row++ {
			for col := 0; col < numCols; col++ {
				tt := eui.NewTooltip(uiRes, "")
				slot := &softwareInstSlot{}
				slotButton := eui.NewSlotButton(uiRes, eui.SlotButtonConfig{
					Tooltip: tt.Container,
					OnHoverStart: func(sb *eui.SlotButton) {
						c.hoverInstSlot = slot
					},
					OnHoverEnd: func(sb *eui.SlotButton) {
						if c.hoverInstSlot == slot {
							c.hoverInstSlot = nil
						}
					},
				})
				grid.AddChild(slotButton.Container)
				slot.tooltipText = tt.Text
				slot.button = slotButton
				c.instList = append(c.instList, slot)
			}
		}
		rows.AddChild(grid)
	}

	root.AddChild(rows)

	initUI(scene, root)

	for _, row := range c.slots {
		for _, slot := range row {
			scene.AddGraphics(slot.warning)
		}
	}

	c.dragLine = graphics.NewLine(c.ctx.GraphicsCache,
		gmath.Pos{Base: &c.dragFrom},
		gmath.Pos{Base: &c.dragTo})
	c.dragLine.SetVisibility(false)
	c.dragLine.SetWidth(2)
	dragLineColor := graphics.ColorScaleFromColor(styles.SelectedTextColor)
	dragLineColor.A = 0.5
	c.dragLine.SetColorScale(dragLineColor)
	scene.AddGraphics(c.dragLine)

	c.updateInstructionSlots()
	c.updateInstBar()
}

func (c *SoftwareController) updateInstBar() {
	thread := c.selectedTab.thread

	index := 0
	for _, k := range game.InstOrder {
		instInfo := game.ProgInstInfoTab[k]
		if instInfo.Mask&thread.Kind == 0 {
			continue
		}
		c.instList[index].inst = game.MakeInst(instInfo.Kind, instInfo.DefaultParam)
		index++
	}

	for index < len(c.instList) {
		c.instList[index].inst = game.MakeInst(game.NopInstruction, 0)
		index++
	}

	for _, slot := range c.instList {
		slot.tooltipText.Label = c.instDoc(slot.inst, true)
		if slot.inst.Info.Icon != assets.ImageNone {
			slot.button.Icon.Image = c.ctx.Loader.LoadImage(slot.inst.Info.Icon).Data
		} else {
			slot.button.Icon.Image = nil
		}
	}
}

func (c *SoftwareController) updateInstructionSlots() {
	thread := c.selectedTab.thread

	for i, row := range c.slots {
		branch := thread.Branches[i]
		for j, b := range row {
			inst := branch.Instructions[j]
			if inst.Info.Kind == game.NopInstruction {
				b.button.Icon.Image = nil
				b.button.Label.Label = ""
				b.tooltipText.Label = c.instDoc(inst, false)
				continue
			}
			b.button.Icon.Image = c.ctx.Loader.LoadImage(inst.Info.Icon).Data
			b.button.Label.Label = c.formatInstParam(inst)
			b.tooltipText.Label = c.instDoc(inst, false)
		}
	}

	// Now validate slots to reject programs that will crash at runtime.
	hasErrors := false
	for i, row := range c.slots {
		var stack []string
		branch := thread.Branches[i]
		for j, slot := range row {
			inst := branch.Instructions[j]
			var errMessage string
			if inst.Info.StackInType != "" {
				switch {
				case len(stack) == 0:
					errMessage = "argument stack is empty"
				case inst.Info.StackInType != stack[len(stack)-1]:
					errMessage = fmt.Sprintf("invalid stack arg (have %s, want %s)",
						stack[len(stack)-1], inst.Info.StackInType)
				}
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			}
			if inst.Info.StackOutType != "" {
				stack = append(stack, inst.Info.StackOutType)
			}
			if errMessage != "" {
				hasErrors = true
				slot.warning.Pos.Offset = c.widgetPos(slot.button.Button.GetWidget()).Add(gmath.Vec{X: -20, Y: 16})
				slot.warning.SetVisibility(true)
				slot.tooltipText.Label += "\n\nError: " + errMessage
			} else {
				slot.warning.SetVisibility(false)
			}
		}
	}

	c.hasErrors = hasErrors
}

func (c *SoftwareController) getSlotInst(slot *softwareSlot) *game.ProgInstruction {
	return &c.selectedTab.thread.Branches[slot.branchIndex].Instructions[slot.instructionIndex]
}

func (c *SoftwareController) selectTab(index int) {
	if c.hasErrors {
		c.ctx.Audio().PlaySound(assets.AudioErrorBeep)
		return
	}

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
	c.updateInstBar()
}

func (c *SoftwareController) formatInstParam(inst game.ProgInstruction) string {
	if !inst.Info.Param {
		return ""
	}

	switch inst.Info.Kind {
	case game.WaitInstruction:
		return strconv.Itoa(int(inst.Param))
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

func (c *SoftwareController) instDoc(inst game.ProgInstruction, instBar bool) string {
	var lines []string

	switch inst.Info.Kind {
	case game.NopInstruction:
		lines = []string{
			"An empty instruction does nothing.",
			"It's also known as NOP instruction",
		}

	case game.WaitInstruction:
		lines = []string{
			"Wait blocks the running branch for the specified",
			"number of processing ticks.",
			"Normally, there are ~60 processing ticks in second.",
			"Therefore, a value of 30 is approximately 1/2 of a second.",
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

	case game.MoveAndRotateInstruction:
		lines = []string{
			"Combines rotation and movement.",
			"Rotates to the destination point.",
			"The engines will not be disabled while turning.",
			"Moves forward for the specified amount of units.",
			"The destination point will be popped from the stack.",
		}

	case game.RandomOffsetInstruction:
		lines = []string{
			"Adds a random offset to the pos",
			"at the top of the stack.",
			"The offset radius is controlled via parameter.",
		}

	case game.CenterPosInstruction:
		lines = []string{
			"Push center of the map pos to the stack.",
		}

	case game.VesselPosInstruction:
		lines = []string{
			"Push own vessel's pos to the stack.",
		}

	case game.TargetPosInstruction:
		lines = []string{
			"Push an enemy pos to the stack.",
		}

	case game.TargetSpeedInstruction:
		lines = []string{
			"Push an enemy speed to the stack.",
			"The enemy speed is a velocity vector length.",
		}

	case game.ChanceInstruction:
		lines = []string{
			"Go to the next branch if roll fails.",
			"The roll success rate is specified via parameter.",
			"Example: a value of 25 (%) fails 1/4 of times.",
		}

	case game.HealthPercentInstruction:
		lines = []string{
			"Push own vessel's health % to the stack.",
		}

	case game.EnergyPercentInstruction:
		lines = []string{
			"Push own vessel's energy % to the stack.",
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
	case game.IsOutBoundsInstruction:
		lines = []string{
			"Pop a top stack position value and check",
			"whether it is out of arena bounds.",
			"If it's not out-of-bounds, go to the next branch.",
		}

	case game.AimShotInstruction:
		lines = []string{
			"Fires weapon at the enemy.",
			"Aimed shot tries to predict the enemy position",
			"while targeting.",
		}

	case game.SnapShotInstruction:
		lines = []string{
			"Fires weapon at the enemy.",
			"Snap shot allows a faster rate-of-fire with",
			"a very small accuracy.",
			"It can be good against unpredictable targets.",
		}

	case game.NormalShotInstruction:
		lines = []string{
			"Fires weapon at the enemy.",
			"Normal shot takes the enemy's current pos",
			"as the aiming point.",
			"It works the best against slow or immobile targets.",
		}
	}

	lines = append(lines, "", fmt.Sprintf("Stack signature: (%s) => (%s)", inst.Info.StackInType, inst.Info.StackOutType))

	if !instBar {
		if inst.Info.Param {
			lines = append(lines, "", "Hover and start typing to change the value.")
			lines = append(lines, "Right click to remove.")
		} else {
			lines = append(lines, "", "Right click to remove.")
		}
	}

	return strings.Join(lines, "\n")
}

func (c *SoftwareController) Update(delta float64) {
	if c.maybeEditValue() {
		c.updateInstructionSlots()
		return
	}
	if c.maybeRemoveSlot() {
		c.updateInstructionSlots()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if !c.hasErrors {
			game.ChangeScene(c.ctx, NewLobbyController(c.ctx))
		} else {
			c.ctx.Audio().PlaySound(assets.AudioErrorBeep)
		}
		return
	}

	c.handleDragAndDrop()

	if runtime.GOARCH != "wasm" {
		if inpututil.IsKeyJustPressed(ebiten.KeyGraveAccent) {
			c.saveVessel()
		}
	}
}

func (c *SoftwareController) saveVessel() {
	session := c.ctx.Session

	saved := game.SavedVessel{
		VesselDesign: session.VesselDesign.Name,
		Prog:         session.Prog.Compact(),
		Artifact:     session.ArtifactDesign.Name,
	}
	for _, wd := range session.Weapons {
		saved.Weapons = append(saved.Weapons, wd.Name)
	}

	jsonData, err := json.MarshalIndent(saved, "", " ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("saved_vessel.json", jsonData, 0o644); err != nil {
		panic(err)
	}
}

func (c *SoftwareController) handleDragAndDrop() {
	if c.draggingInst == nil {
		// Maybe start dragging.

		if c.hoverInstSlot == nil {
			return
		}
		if c.hoverInstSlot.inst.Info.Kind == game.NopInstruction {
			return
		}
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return
		}
		c.draggingInst = c.hoverInstSlot
		c.dragFrom = c.widgetPos(c.hoverInstSlot.button.Button.GetWidget())
		c.dragLine.SetVisibility(true)
	}

	cursorX, cursorY := ebiten.CursorPosition()
	c.dragTo = gmath.Vec{X: float64(cursorX), Y: float64(cursorY)}
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}

	if c.hoverSlot != nil {
		inst := c.getSlotInst(c.hoverSlot)
		*inst = c.draggingInst.inst
		c.ctx.Audio().PlaySound(assets.AudioAckBeep)
		c.updateInstructionSlots()
	}

	c.draggingInst = nil
	c.dragLine.SetVisibility(false)
}

func (c *SoftwareController) widgetPos(w *widget.Widget) gmath.Vec {
	r := w.Rect
	return gmath.Vec{
		X: float64(r.Min.X + r.Dx()/2),
		Y: float64(r.Min.Y + r.Dy()/2),
	}
}

func (c *SoftwareController) maybeRemoveSlot() bool {
	if c.hoverSlot == nil {
		return false
	}

	inst := c.getSlotInst(c.hoverSlot)
	if inst.Info.Kind == game.NopInstruction {
		return false
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		*inst = game.MakeInst(game.NopInstruction, 0)
		c.ctx.Audio().PlaySound(assets.AudioDeleteBeep)
		return true
	}

	return false
}

func (c *SoftwareController) maybeEditValue() bool {
	if c.hoverSlot == nil {
		return false
	}

	inst := c.getSlotInst(c.hoverSlot)

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		c.ctx.Audio().PlaySound(assets.AudioClickBeep)
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
		c.ctx.Audio().PlaySound(assets.AudioClickBeep)
		kv := k - ebiten.KeyDigit0
		s := strconv.Itoa(int(inst.Param))
		s += strconv.Itoa(int(kv))
		v, _ := strconv.ParseFloat(s, 64)
		inst.SetParam(v)
		return true
	}

	return false
}
