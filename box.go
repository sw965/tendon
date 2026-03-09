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
	CrossAlignment Alignment // 交差軸方向の配置
	AutoCompress   bool
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
		CrossAlignment: AlignCenter, // デフォルトは中央揃え
		AutoCompress:   true,
	}
}

func (b *Box) RelPositions(o Orientation) ([]float64, []float64) {
	if len(b.Children) == 0 {
		return nil, nil
	}
	if o == Horizontal {
		return b.horizontalRelPositions()
	}
	return b.verticalRelPositions()
}

func (b *Box) horizontalRelPositions() ([]float64, []float64) {
	n := len(b.Children)
	childrenW := 0.0
	childWs := make([]float64, n)

	for i, child := range b.Children {
		w := child.BaseWidth() * child.WidthScale
		childWs[i] = w
		childrenW += w
	}

	parentW := b.BaseWidth()
	gap := b.Gap
	totalW := childrenW + float64(n-1)*b.Gap

	// (要素幅の合計 + 隙間の合計幅) が コンテナの幅より大きい場合 (コンテナから要素がはみ出た場合)
	if b.AutoCompress && totalW > parentW && n > 1 {
		// コンテナの幅に収まるようなgapを逆算する
		// (要素幅の合計) + (隙間 × 個数-1) = 親の幅 この方程式をgap(隙間)について解く
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

	for i, child := range b.Children {
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

		// 自分の幅と隙間を加算
		currentX += childWs[i] + gap
	}
	return xs, ys
}

func (b *Box) verticalRelPositions() ([]float64, []float64) {
	n := len(b.Children)
	childrenH := 0.0
	childHs := make([]float64, n)

	for i, child := range b.Children {
		h := child.BaseHeight() * child.HeightScale
		childHs[i] = h
		childrenH += h
	}

	parentH := b.BaseHeight()
	gap := b.Gap
	totalH := childrenH + float64(n-1)*b.Gap

	// (要素高さの合計 + 隙間の合計高さ) が コンテナの高さより大きい場合 (コンテナから要素がはみ出た場合)
	if b.AutoCompress && totalH > parentH && n > 1 {
		// コンテナの高さに収まるようなgapを逆算する
		// (要素高さの合計) + (隙間 × 個数-1) = 親の高さ この方程式をgap(隙間)について解く
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

	for i, child := range b.Children {
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

		// 自分の高さと隙間を加算
		currentY += childHs[i] + gap
	}
	return xs, ys
}

func (b *Box) Update(o Orientation) {
	xs, ys := b.RelPositions(o)
	for i, x := range xs {
		if i < len(b.Children) {
			b.Children[i].XRelativeToParent = x
			b.Children[i].YRelativeToParent = ys[i]
		}
	}
}
