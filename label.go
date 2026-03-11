package tendon

import (
	"bytes"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func NewDefaultFontSource() (*text.GoTextFaceSource, error) {
	return text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
}

type Label struct {
	*Element
	font *text.GoTextFace
	text string
}

func NewLabel(txt string, size float64) (*Label, error) {
	src, err := NewDefaultFontSource()
	if err != nil {
		return nil, err
	}

	l := &Label{
		Element: NewElement(),
	}
	l.Filter = ebiten.FilterLinear
	l.PassThrough = true

	l.SetText(txt, src, size)
	return l, nil
}

func (l *Label) Font() *text.GoTextFace {
	return l.font
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) SetText(txt string, src *text.GoTextFaceSource, size float64) {
	if src == nil {
		return
	}

	// 変更がなければ、計算をスキップして負荷を下げる
	if l.text == txt && l.font != nil && l.font.Source == src && l.font.Size == size && l.Image != nil {
		return
	}

	l.text = txt
	l.font = &text.GoTextFace{
		Source: src,
		Size:   size,
	}

	// テキストサイズの計測
	w, h := text.Measure(l.text, l.font, 0)
	if w <= 0 || h <= 0 {
		l.Image = nil
		return
	}

	// math.Ceil で切り上げることで、端数による描画欠けを防止
	img := ebiten.NewImage(int(math.Ceil(w)), int(math.Ceil(h)))
	text.Draw(img, l.text, l.font, &text.DrawOptions{})
	l.Image = img
}