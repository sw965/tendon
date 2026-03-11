package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// BoxPosition は計算された各要素の相対座標を保持します。
type BoxPosition struct {
	RelX float64
	RelY float64
}

type Box struct {
	*Element
	Gap            float64
	MainAlignment  Alignment
	CrossAlignment Alignment
	AutoCompress   bool
	Orientation    Orientation // ★追加：Update() を引数なしにするため方向を保持する
}

// NewBox は指定されたサイズ(w, h)と隙間(gap)で新しい整列コンテナを作成します。
func NewBox(w, h, gap float64) *Box {
	container := NewElement()
	if w > 0 && h > 0 {
		container.Image = ebiten.NewImage(int(w), int(h))
	}

	return &Box{
		Element:        container,
		Gap:            gap,
		MainAlignment:  AlignCenter,
		CrossAlignment: AlignCenter,
		AutoCompress:   true,
		Orientation:    Horizontal, // ★デフォルトは水平方向に設定
	}
}

func (b *Box) RelPositions() ([]float64, []float64) {
	if len(b.Children) == 0 {
		return nil, nil
	}
	// ★引数ではなく自分自身の保持する方向を使う
	if b.Orientation == Horizontal {
		return b.horizontalRelPositions()
	}
	return b.verticalRelPositions()
}

func (b *Box) horizontalRelPositions() ([]float64, []float64) {
	n := len(b.Children)
	childrenW := 0.0
	childWs := make([]float64, n)

	for i, comp := range b.Children {
		child := comp.BaseElement() // ★ BaseElement() で実体を取り出す
		w := child.BaseWidth() * child.WidthScale
		childWs[i] = w
		childrenW += w
	}

	parentW := b.BaseWidth()
	gap := b.Gap
	totalW := childrenW + float64(n-1)*b.Gap

	if b.AutoCompress && totalW > parentW && n > 1 {
		gap = (parentW - childrenW) / float64(n-1)
		totalW = parentW
	}

	var currentX float64
	switch b.MainAlignment {
	case AlignStart:
		currentX = 0
	case AlignCenter:
		currentX = (parentW - totalW) / 2
	case AlignEnd:
		currentX = parentW - totalW
	}

	parentH := b.BaseHeight()
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i, c := range b.Children {
		child := c.BaseElement() // ★ BaseElement() で実体を取り出す
		childH := child.BaseHeight() * child.HeightScale

		var y float64
		switch b.CrossAlignment {
		case AlignStart:
			y = 0
		case AlignCenter:
			y = (parentH - childH) / 2
		case AlignEnd:
			y = parentH - childH
		}

		xs[i] = currentX
		ys[i] = y
		currentX += childWs[i] + gap
	}
	return xs, ys
}

func (b *Box) verticalRelPositions() ([]float64, []float64) {
	n := len(b.Children)
	childrenH := 0.0
	childHs := make([]float64, n)

	for i, c := range b.Children {
		child := c.BaseElement() // ★ BaseElement() で実体を取り出す
		h := child.BaseHeight() * child.HeightScale
		childHs[i] = h
		childrenH += h
	}

	parentH := b.BaseHeight()
	gap := b.Gap
	totalH := childrenH + float64(n-1)*b.Gap

	if b.AutoCompress && totalH > parentH && n > 1 {
		gap = (parentH - childrenH) / float64(n-1)
		totalH = parentH
	}

	var currentY float64
	switch b.MainAlignment {
	case AlignStart:
		currentY = 0
	case AlignCenter:
		currentY = (parentH - totalH) / 2
	case AlignEnd:
		currentY = parentH - totalH
	}

	parentW := b.BaseWidth()
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i, c := range b.Children {
		child := c.BaseElement()
		childW := child.BaseWidth() * child.WidthScale

		var x float64
		switch b.CrossAlignment {
		case AlignStart:
			x = 0
		case AlignCenter:
			x = (parentW - childW) / 2
		case AlignEnd:
			x = parentW - childW
		}

		xs[i] = x
		ys[i] = currentY
		currentY += childHs[i] + gap
	}
	return xs, ys
}

func (b *Box) Update() {
	b.Element.Update()

	xs, ys := b.RelPositions()
	for i, x := range xs {
		if i < len(b.Children) {
			child := b.Children[i].BaseElement()
			child.XRelativeToParent = x
			child.YRelativeToParent = ys[i]
		}
	}
}