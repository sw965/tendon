package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

type Box struct {
	*Element
	Gap            float64
	MainAlignment  Alignment
	CrossAlignment Alignment
	AutoCompress   bool
	Orientation    Orientation
	LayoutChildren Components
}

func NewBox(w, h, gap float64) *Box {
	base := NewElement()
	if w > 0 && h > 0 {
		base.Image = ebiten.NewImage(int(w), int(h))
	}

	return &Box{
		Element:        base,
		Gap:            gap,
		MainAlignment:  AlignCenter,
		CrossAlignment: AlignCenter,
		AutoCompress:   true,
		Orientation:    Horizontal,
		LayoutChildren: Components{},
	}
}

func (b *Box) AppendChild(child Component) {
	b.Element.AppendChild(child)
	b.LayoutChildren = append(b.LayoutChildren, child)
}

func (b *Box) RemoveChild(target Component) bool {
	removed := b.Element.RemoveChild(target)
	if !removed {
		return false
	}

	t := target.BaseElement()
	index := slices.IndexFunc(b.LayoutChildren, func(c Component) bool {
		return c.BaseElement() == t
	})

	if index != -1 {
		b.LayoutChildren = slices.Delete(b.LayoutChildren, index, index+1)
	}

	return true
}

func (b *Box) RelPositions() ([]float64, []float64) {
	if len(b.LayoutChildren) == 0 {
		return nil, nil
	}

	if b.Orientation == Horizontal {
		return b.horizontalRelPositions()
	}
	return b.verticalRelPositions()
}

func (b *Box) horizontalRelPositions() ([]float64, []float64) {
	n := len(b.LayoutChildren)
	childrenW := 0.0
	childWs := make([]float64, n)

	for i, comp := range b.LayoutChildren {
		child := comp.BaseElement()
		w := child.BaseWidth() * child.widthScale
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

	for i, c := range b.LayoutChildren {
		child := c.BaseElement() // ★ BaseElement() で実体を取り出す
		childH := child.BaseHeight() * child.heightScale

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
	n := len(b.LayoutChildren)
	childrenH := 0.0
	childHs := make([]float64, n)

	for i, c := range b.LayoutChildren {
		child := c.BaseElement() // ★ BaseElement() で実体を取り出す
		h := child.BaseHeight() * child.heightScale
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

	for i, c := range b.LayoutChildren {
		child := c.BaseElement()
		childW := child.BaseWidth() * child.widthScale

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

func (b *Box) Reflow() {
	xs, ys := b.RelPositions()
	for i, x := range xs {
		if i < len(b.LayoutChildren) {
			child := b.LayoutChildren[i].BaseElement()
			child.XRelativeToParent = x
			child.YRelativeToParent = ys[i]
		}
	}
}
