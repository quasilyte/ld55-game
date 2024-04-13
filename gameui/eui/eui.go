package eui

import (
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/styles"
	"golang.org/x/image/font"
)

type Resources struct {
	button     *buttonResource
	slotButton *buttonResource
	tooltip    *tooltipResources
}

type buttonResource struct {
	Image      *widget.ButtonImage
	Padding    widget.Insets
	TextColors *widget.ButtonTextColor
}

type tooltipResources struct {
	Background *image.NineSlice
	Padding    widget.Insets
	FontFace   font.Face
	TextColor  color.Color
}

func LoadResources(loader *resource.Loader) *Resources {
	res := &Resources{}

	{
		disabled := nineSliceImage(loader.LoadImage(assets.ImageUIButtonDisabled).Data, 16, 16)
		idle := nineSliceImage(loader.LoadImage(assets.ImageUIButtonIdle).Data, 16, 16)
		hover := nineSliceImage(loader.LoadImage(assets.ImageUIButtonHover).Data, 16, 16)
		pressed := nineSliceImage(loader.LoadImage(assets.ImageUIButtonPressed).Data, 16, 16)
		buttonPadding := widget.Insets{
			Left:   24,
			Right:  24,
			Top:    8,
			Bottom: 8,
		}
		buttonColors := &widget.ButtonTextColor{
			Idle:     styles.NormalTextColor,
			Disabled: styles.DisabledTextColor,
		}
		res.button = &buttonResource{
			Image: &widget.ButtonImage{
				Idle:     idle,
				Hover:    hover,
				Pressed:  pressed,
				Disabled: disabled,
			},
			Padding:    buttonPadding,
			TextColors: buttonColors,
		}
	}

	{
		idle := nineSliceImage(loader.LoadImage(assets.ImageUISlotIdle).Data, 16, 16)
		hover := nineSliceImage(loader.LoadImage(assets.ImageUISlotHover).Data, 16, 16)
		disabled := nineSliceImage(loader.LoadImage(assets.ImageUISlotDisabled).Data, 16, 16)
		buttonPadding := widget.Insets{
			Left:   8,
			Right:  8,
			Top:    8,
			Bottom: 8,
		}
		buttonColors := &widget.ButtonTextColor{
			Idle:     styles.NormalTextColor,
			Disabled: styles.DisabledTextColor,
		}
		res.slotButton = &buttonResource{
			Image: &widget.ButtonImage{
				Idle:     idle,
				Hover:    hover,
				Pressed:  hover,
				Disabled: disabled,
			},
			Padding:    buttonPadding,
			TextColors: buttonColors,
		}
	}

	{
		res.tooltip = &tooltipResources{
			Background: nineSliceImage(loader.LoadImage(assets.ImageUITooltip).Data, 18, 18),
			Padding: widget.Insets{
				Left:   16,
				Right:  16,
				Top:    10,
				Bottom: 10,
			},
			FontFace:  assets.Font1,
			TextColor: styles.NormalTextColor,
		}
	}

	return res
}

type ButtonConfig struct {
	Text       string
	OnClick    func()
	LayoutData any
	MinWidth   int
	Font       font.Face
	AlignLeft  bool

	Tooltip *widget.Container
}

func NewSeparator(ld interface{}, clr color.RGBA) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(clr)),
	))

	return c
}

func NewCenteredLabel(text string, ff font.Face) *widget.Text {
	return NewCenteredLabelWithMaxWidth(text, ff, -1)
}

func NewCenteredLabelWithMaxWidth(text string, ff font.Face, width float64) *widget.Text {
	options := []widget.TextOpt{
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Text(text, ff, styles.NormalTextColor),
	}
	if width != -1 {
		options = append(options, widget.TextOpts.MaxWidth(width))
	}
	return widget.NewText(options...)
}

type SlotButtonConfig struct {
	OnClick func()

	Tooltip *widget.Container
}

type SlotButton struct {
	Icon      *widget.Graphic
	Button    *widget.Button
	Container *widget.Container
}

func NewSlotButton(res *Resources, config SlotButtonConfig) *SlotButton {
	buttonRes := res.slotButton

	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(64, 64),
		),
	)

	options := []widget.ButtonOpt{
		widget.ButtonOpts.Image(buttonRes.Image),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if config.OnClick != nil {
				config.OnClick()
			}
		}),
		// widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(64, 64)),
	}

	if config.Tooltip != nil {
		tt := widget.NewToolTip(
			widget.ToolTipOpts.Content(config.Tooltip),
			widget.ToolTipOpts.Delay(time.Second),
		)
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(tt)))
	}

	b := widget.NewButton(options...)
	container.AddChild(b)

	g := widget.NewGraphic()
	container.AddChild(g)

	return &SlotButton{
		Button:    b,
		Icon:      g,
		Container: container,
	}
}

func NewButton(res *Resources, config ButtonConfig) *widget.Button {
	buttonRes := res.button

	ff := config.Font
	if ff == nil {
		ff = assets.Font2
	}
	options := []widget.ButtonOpt{
		widget.ButtonOpts.Image(buttonRes.Image),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if config.OnClick != nil {
				config.OnClick()
			}
		}),
	}
	if config.Tooltip != nil {
		tt := widget.NewToolTip(
			widget.ToolTipOpts.Content(config.Tooltip),
			widget.ToolTipOpts.Delay(time.Second),
		)
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(tt)))
	}
	colors := &widget.ButtonTextColor{
		Idle:     buttonRes.TextColors.Idle,
		Disabled: buttonRes.TextColors.Disabled,
	}
	options = append(options,
		widget.ButtonOpts.Text(config.Text, ff, colors),
		widget.ButtonOpts.TextPadding(buttonRes.Padding))
	if config.AlignLeft {
		options = append(options, widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter))
	}
	if config.LayoutData != nil {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(config.LayoutData)))
	}
	if config.MinWidth != 0 {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)))
	}

	return widget.NewButton(options...)
}

type Tooltip struct {
	Container *widget.Container
	Text      *widget.Text
}

func NewTooltip(res *Resources, s string) *Tooltip {
	tt := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.tooltip.Background),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.tooltip.Padding),
			widget.RowLayoutOpts.Spacing(2),
		)))
	text := widget.NewText(
		widget.TextOpts.MaxWidth(800),
		widget.TextOpts.Text(s, res.tooltip.FontFace, res.tooltip.TextColor),
	)
	tt.AddChild(text)
	return &Tooltip{
		Container: tt,
		Text:      text,
	}
}

func NewSimpleTooltip(res *Resources, text string) *widget.Container {
	return NewTooltip(res, text).Container
}

func NewColoredLabel(text string, ff font.Face, clr color.RGBA, options ...widget.TextOpt) *widget.Text {
	opts := []widget.TextOpt{
		widget.TextOpts.Text(text, ff, clr),
	}
	if len(options) != 0 {
		opts = append(opts, options...)
	}
	return widget.NewText(opts...)
}

func NewLabel(text string, ff font.Face, options ...widget.TextOpt) *widget.Text {
	return NewColoredLabel(text, ff, styles.NormalTextColor, options...)
}

func nineSliceImage(i *ebiten.Image, offsetX, offsetY int) *image.NineSlice {
	size := i.Bounds().Size()
	w := size.X
	h := size.Y
	return image.NewNineSlice(i,
		[3]int{offsetX, w - 2*offsetX, offsetX},
		[3]int{offsetY, h - 2*offsetY, offsetY},
	)
}
