package tendon

import (
	"bytes"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	defaultFaceSource *text.GoTextFaceSource
	DefaultFontSize   float64 = 16
)

// NewDefaultFace は指定されたサイズのフォントフェイスを生成します。
func NewDefaultFace(size float64) *text.GoTextFace {
	if defaultFaceSource == nil {
		var err error
		defaultFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			panic(err)
		}
	}
	return &text.GoTextFace{
		Source: defaultFaceSource,
		Size:   size,
	}
}

type Label struct {
	*Element
	text     string
	fontSize float64
}

// NewLabel は初期テキストとフォントサイズを指定してラベルを作成します。
func NewLabel(text string, size float64) *Label {
	l := &Label{
		Element:  NewElement(),
		fontSize: size,
	}
	l.Filter = ebiten.FilterLinear
	l.PassThrough = true
	l.SetText(text)
	return l
}

// SetSize はフォントサイズを変更し、画像をクッキリした状態で再生成します。
func (l *Label) SetSize(size float64) {
	if l.fontSize == size {
		return
	}
	l.fontSize = size
	l.render()
}

// SetText は現在のサイズを維持したままテキストを更新します。
func (l *Label) SetText(txt string) {
	if l.text == txt && l.Image != nil {
		return
	}
	l.text = txt
	l.render()
}

// render は現在のテキストとサイズで、最適な解像度の画像を生成します。
func (l *Label) render() {
	face := NewDefaultFace(l.fontSize)
	
	w, h := text.Measure(l.text, face, 0)
	if w <= 0 || h <= 0 {
		l.Image = nil
		return
	}

	// 拡大によるボヤけを防ぐため、常にネイティブサイズで描画
	img := ebiten.NewImage(int(math.Ceil(w)), int(math.Ceil(h)))
	text.Draw(img, l.text, face, &text.DrawOptions{})
	
	l.Image = img
}

func (l *Label) Text() string {
	return l.text
}