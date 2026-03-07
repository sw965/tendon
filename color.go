package tendon

import (
	"image/color"
)

var (
	Transparent = color.RGBA{0, 0, 0, 0}         // 透明
	White       = color.RGBA{255, 255, 255, 255} // しろ
	Black       = color.RGBA{0, 0, 0, 255}       // くろ
	Gray        = color.RGBA{128, 128, 128, 255} // はいいろ

	Red         = color.RGBA{255, 0, 0, 255}     // あか
	Blue        = color.RGBA{0, 0, 255, 255}     // あお
	Yellow      = color.RGBA{255, 255, 0, 255}   // きいろ
	Green       = color.RGBA{0, 128, 0, 255}     // みどり
	Orange      = color.RGBA{255, 165, 0, 255}   // だいだい
	Brown       = color.RGBA{139, 69, 19, 255}   // ちゃいろ
	Pink        = color.RGBA{255, 192, 203, 255} // ももいろ
	LightBlue   = color.RGBA{135, 206, 235, 255} // みずいろ
	YellowGreen = color.RGBA{154, 205, 50, 255}  // きみどり
	Purple      = color.RGBA{128, 0, 128, 255}   // むらさき
	PaleOrange  = color.RGBA{255, 218, 185, 255} // うすだいだい

	Ocher     = color.RGBA{184, 134, 11, 255} // おうどいろ
	Vermilion = color.RGBA{227, 66, 52, 255}  // しゅいろ
)

func WithAlpha(c color.Color, alpha float64) color.RGBA {
	r, g, b, _ := c.RGBA()
	a8 := uint8(alpha * 255)
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), a8}
}
