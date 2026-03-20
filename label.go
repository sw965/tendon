package tendon

import (
	"bytes"
	"math"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func NewDefaultFontSource() (*text.GoTextFaceSource, error) {
	return text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
}

type Label struct {
	*Element
	font  *text.GoTextFace
	text  string
	color color.Color
}

func NewLabel(txt string, size float64) (*Label, error) {
	src, err := NewDefaultFontSource()
	if err != nil {
		return nil, err
	}

	l := &Label{
		Element: NewElement(),
		color:Black,
	}
	l.Filter = ebiten.FilterLinear
	l.PassThrough = true

	l.SetText(src, size, Black, txt)
	return l, nil
}

func (l *Label) Font() *text.GoTextFace {
	return l.font
}

func (l *Label) Color() color.Color {
	return l.color
}

func (l *Label) Text() string {
	return l.text
}

// TODO 引数が多いから別の設計を考える？
func (l *Label) SetText(src *text.GoTextFaceSource, size float64, clr color.Color, txt string) {
	if src == nil {
		return
	}

    r1, g1, b1, a1 := l.color.RGBA()
    r2, g2, b2, a2 := clr.RGBA()
    colorUnchanged := r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2

    // 全ての条件が一致していれば計算をスキップ
    if l.font != nil && l.font.Source == src && l.font.Size == size && 
       l.text == txt && colorUnchanged && l.Image != nil {
        return
    }

	l.font = &text.GoTextFace{
		Source: src,
		Size:   size,
	}
	l.color = clr
	l.text = txt

	// テキストサイズの計測
	w, h := text.Measure(l.text, l.font, 0)
	if w <= 0 || h <= 0 {
		if l.Image != nil {
			l.Image.Dispose()
			l.Image = nil
		}
		return
	}

	// math.Ceil で切り上げることで、端数による描画欠けを防止
	img := ebiten.NewImage(int(math.Ceil(w)), int(math.Ceil(h)))
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(l.color) 
	text.Draw(img, l.text, l.font, op)

	// 古い画像があればVRAMから解放する
	if l.Image != nil {
		l.Image.Dispose()
	}
	l.Image = img
}
