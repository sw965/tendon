package tendon

import (
	"slices"
	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Update()
	Draw(screen *ebiten.Image)
	BaseElement() *Element
}

type Components []Component

func (cs Components) SortByZAsc() {
	// a, b は Component インターフェースなので、BaseElement() を経由して Z を取得する
	slices.SortStableFunc(cs, func(a, b Component) int {
		return a.BaseElement().Z - b.BaseElement().Z
	})
}

func (cs Components) SortByZDesc() {
	slices.SortStableFunc(cs, func(a, b Component) int {
		return b.BaseElement().Z - a.BaseElement().Z
	})
}

func (cs Components) FindAllFromPoint(pointX, pointY float64, dst *Components) {
	cs.SortByZDesc()
	for _, c := range cs {
		e := c.BaseElement()
		if !e.Visible {
			continue
		}
		
		// 子要素を再帰的に探索
		e.Children.FindAllFromPoint(pointX, pointY, dst)

		// 自分自身を判定
		if !e.PassThrough && e.Contains(pointX, pointY) {
			*dst = append(*dst, c)
		}
	}
}

func (cs Components) FindOverlapping(target Component, dst *Components) {
	t := target.BaseElement()
	for _, c := range cs {
		e := c.BaseElement()
		if !e.Visible || e.PassThrough || e == t {
			continue
		}
		if e.Overlaps(target) {
			*dst = append(*dst, c)
		}
	}
}

func (cs Components) FindAllOverlapping(target Component, dst *Components) {
	t := target.BaseElement()
	for _, c := range cs {
		e := c.BaseElement()
		if !e.Visible {
			continue
		}

		e.Children.FindAllOverlapping(target, dst)

		if e != t && !e.PassThrough && e.Overlaps(target) {
			*dst = append(*dst, c)
		}
	}
}

func (cs Components) StopAllDrag() {
	for _, c := range cs {
		e := c.BaseElement()
		e.StopDrag()
		e.Children.StopAllDrag()
	}
}

func (cs Components) UpdateHover(hitTest Components) {
	hitSet := make(map[Component]bool)
	for _, c := range hitTest {
		hitSet[c] = true
	}

	var update func(Components)
	update = func(elements Components) {
		for _, c := range elements {
			e := c.BaseElement()
			isHit := hitSet[c] && e.Enabled
			e.isJustHoverIn = false
			e.isJustHoverOut = false

			if isHit && !e.isHovered {
				e.isHovered = true
				e.isJustHoverIn = true
			} else if !isHit && e.isHovered {
				e.isHovered = false
				e.isJustHoverOut = true
			}
			update(e.Children)
		}
	}
	update(cs)
}

func (cs Components) UpdateDragMove() {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float64(cursorX), float64(cursorY)

	for _, c := range cs {
		e := c.BaseElement()
		if !e.isDragging {
			continue
		}

		absWsc, absHsc := e.AbsWidthScale(), e.AbsHeightScale()
		toAbsX := cursorXf - (e.dragOffsetX * absWsc)
		toAbsY := cursorYf - (e.dragOffsetY * absHsc)
		toRelX, toRelY := e.AbsPosToRelPosToParent(toAbsX, toAbsY)

		e.DragDeltaX = toRelX - e.XRelativeToParent
		e.DragDeltaY = toRelY - e.YRelativeToParent

		if !e.ManualDrag {
			e.XRelativeToParent = toRelX
			e.YRelativeToParent = toRelY
		}
	}
}

func (cs Components) Update() {
	for _, c := range cs {
		c.Update()
	}
}

func (cs Components) Draw(screen *ebiten.Image) {
	cs.SortByZAsc()
	for _, c := range cs {
		c.Draw(screen)
	}
}