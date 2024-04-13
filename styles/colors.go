package styles

import (
	"image/color"
)

var (
	NormalTextColor   = rgb(0xFFFFFF)
	DisabledTextColor = rgb(0x777777)
	SelectedTextColor = rgb(0xd29356)

	AlliedColor = rgb(0x65cc5e)
	EnemyColor  = rgb(0xcc5e5e)

	UIDarkBorder = rgb(0x1a263c)
	UINormal     = rgb(0x2b426a)

	TransparentColor = color.RGBA{}
)

func rgb(v uint64) color.RGBA {
	r := uint8((v & (0xFF << (8 * 2))) >> (8 * 2))
	g := uint8((v & (0xFF << (8 * 1))) >> (8 * 1))
	b := uint8((v & (0xFF << (8 * 0))) >> (8 * 0))
	return color.RGBA{R: r, G: g, B: b, A: 0xff}
}
