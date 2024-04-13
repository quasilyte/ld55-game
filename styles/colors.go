package styles

import (
	"image/color"
)

var (
	NormalTextColor   = rgb(0xFFFFFF)
	DisabledTextColor = rgb(0x777777)
	SelectedTextColor = rgb(0xF4418A)

	AlliedColor = rgb(0x65cc5e)
	EnemyColor  = rgb(0xcc5e5e)

	TransparentColor = color.NRGBA{}
)

func rgb(v uint64) color.NRGBA {
	r := uint8((v & (0xFF << (8 * 2))) >> (8 * 2))
	g := uint8((v & (0xFF << (8 * 1))) >> (8 * 1))
	b := uint8((v & (0xFF << (8 * 0))) >> (8 * 0))
	return color.NRGBA{R: r, G: g, B: b, A: 0xff}
}
