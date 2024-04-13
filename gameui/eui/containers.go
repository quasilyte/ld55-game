package eui

import (
	"github.com/ebitenui/ebitenui/widget"
)

func NewRootContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
			StretchVertical:   true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
}

type RowContainerConfig struct {
	RowScale []bool
	Spacing  int
	MinWidth int
}

func NewRowContainer(config RowContainerConfig) *widget.Container {
	if config.Spacing == 0 {
		config.Spacing = 4
	}

	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
			StretchVertical:   true,
		})),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, config.RowScale),
			widget.GridLayoutOpts.Spacing(config.Spacing, config.Spacing),
		)),
	)
}
