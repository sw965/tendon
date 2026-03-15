package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
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

	// 再帰的にドラッグ状態を更新する内部関数
	var update func(Components)
	update = func(elements Components) {
		for _, c := range elements {
			e := c.BaseElement()
			if e.isDragging {
				// 【行列ベースのリファクタリング】
				// マウスの絶対座標を、親のローカル空間（回転・スケール適用前）の座標に逆変換する
				toRelX, toRelY := e.AbsPosToRelPosToParent(cursorXf, cursorYf)
				
				// 開始時のマウスと要素の距離（オフセット）を維持する位置を計算
				targetX := toRelX - e.dragOffsetX
				targetY := toRelY - e.dragOffsetY

				e.DragDeltaX = targetX - e.XRelativeToParent
				e.DragDeltaY = targetY - e.YRelativeToParent

				if !e.ManualDrag {
					e.XRelativeToParent = targetX
					e.YRelativeToParent = targetY
				}
			}
			// 子要素も忘れずにチェック（これで tmp3_test.go が動くようになります）
			update(e.Children)
		}
	}
	update(cs)
}

func (cs Components) Update() {
	for _, c := range cs {
		c.Update()
	}
}

func (cs Components) Draw(screen *ebiten.Image) {
	// TODO ダーディーを追加して計算量を省く
	cs.SortByZAsc()
	for _, c := range cs {
		c.Draw(screen)
	}
}
