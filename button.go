package tendon

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	*Element
	Label *Label
}

func NewButton(w, h float64, txt string, bgColor color.Color) (*Button, error) {
	base := NewElement()
	base.Image = ebiten.NewImage(int(w), int(h))
	base.Image.Fill(bgColor)

	fontSize := h * 0.6
	src, err := NewDefaultFontSource()
	if err != nil {
		return nil, err
	}
	font := &text.GoTextFace{Source: src, Size: fontSize}
	txtW, _ := text.Measure(txt, font, 0)

	maxWidth := w * 0.9
	// ボタンからテキストが横に飛び出る場合、飛び出ないように修正する
	if txtW > maxWidth {
		fontSize = fontSize * (maxWidth / txtW)
	}

	l, err := NewLabel(txt, fontSize)
	if err != nil {
		return nil, err
	}

	b := &Button{
		Element: base,
		Label:   l,
	}

	b.AppendChild(l)
	l.PlaceCenterOf(base)
	return b, nil
}

func (b *Button) SetText(txt string) {
	h := b.BaseHeight()
	w := b.BaseWidth()
	fontSize := h * 0.6

	font := &text.GoTextFace{Source: b.Label.Font().Source, Size: fontSize}
	txtW, _ := text.Measure(txt, font, 0)

	maxWidth := w * 0.9
	if txtW > maxWidth {
		fontSize = fontSize * (maxWidth / txtW)
	}

	b.Label.SetText(txt, b.Label.Font().Source, fontSize)
	b.Label.PlaceCenterOf(b.Element)
}
