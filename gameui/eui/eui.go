package eui

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ld55-game/assets"
	"github.com/quasilyte/ld55-game/styles"
	"golang.org/x/image/font"
)

type Resources struct {
	button *buttonResource
}

type buttonResource struct {
	Image      *widget.ButtonImage
	Padding    widget.Insets
	TextColors *widget.ButtonTextColor
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

	return res
}

type ButtonConfig struct {
	Text       string
	OnClick    func()
	LayoutData any
	MinWidth   int
	Font       font.Face
	AlignLeft  bool
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
	colors := buttonRes.TextColors
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

func nineSliceImage(i *ebiten.Image, offsetX, offsetY int) *image.NineSlice {
	size := i.Bounds().Size()
	w := size.X
	h := size.Y
	return image.NewNineSlice(i,
		[3]int{offsetX, w - 2*offsetX, offsetX},
		[3]int{offsetY, h - 2*offsetY, offsetY},
	)
}
