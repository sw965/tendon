// その都度Zソートしているため、パフォーマンス改善をする必要あり
package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
)

// Element はUIツリーを構成する基本コンポーネントです。
// 相対座標の計算、親子関係の構築、描画、およびマウス操作の判定を担います。
//
// 【3つの主要フラグの違いと、子要素への影響】
//
// 以下のフラグを組み合わせることで、様々なUIの要件（コンテナ、装飾、無効化など）を実現します。
//
//		フラグ名    | 意味           | 描画 | 当たり判定 | 状態更新 | 子要素への影響
//		------------|----------------|------|------------|----------|----------------------------------------
//		Visible     | 可視性         |  ❌  |     ❌     |    ❌    | falseの場合、子要素も全て非表示・判定外になる
//		Enabled     | 有効/無効      |  ⭕️  |     ⭕️     |    ❌    | falseにしても、子要素のEnabledには影響しない
//		PassThrough | 当たり判定透過 |  ⭕️  |     ❌     |    ❌    | trueにしても、子要素の当たり判定は通常通り行われる
//
//	  - Visible:
//	    要素を完全にツリーから除外したい場合に使用します（例：非表示になったウィンドウ）。
//
//	  - Enabled:
//	    クリック等の対象として「重なり」は検知させたいが、操作は弾きたい場合に使用します。
//	    （例：グレーアウトして押せなくなったボタン）。
//
//	  - PassThrough:
//	    見た目だけで当たり判定を下に貫通させたい場合や、当たり判定計算の負荷を下げたい場合に使用します。
//	    （例：ボタンの上に乗っているテキスト、複数のボタンを並べるための透明なレイアウト用親コンテナ）。
type Element struct {
	Id                int
	XRelativeToParent float64
	YRelativeToParent float64

	Image       *ebiten.Image
	WidthScale  float64
	HeightScale float64
	Filter      ebiten.Filter

	Visible     bool
	Enabled     bool
	PassThrough bool
	Z           int

	// ドラッグ関連
	Draggable   bool
	ManualDrag  bool
	isDragging  bool
	dragOffsetX float64
	dragOffsetY float64
	DragDeltaX  float64
	DragDeltaY  float64

	isHovered      bool
	isJustHoverIn  bool
	isJustHoverOut bool

	// 移動アニメーション関連
	toX                float64
	toY                float64
	isMoving           bool
	isJustMoveFinished bool
	EasingFunc         EasingFunc

	// カプセル化の検討
	Children Components
	Parent   *Element
}

func NewElement() *Element {
	e := &Element{
		Visible: true,
		Enabled: true,
		EasingFunc: func(current, target float64) float64 { return target },
	}
	e.SetScale(1.0)
	return e
}

func (e *Element) IsDragging() bool {
	return e.isDragging
}

func (e *Element) DragDelta() (float64, float64) {
	return e.DragDeltaX, e.DragDeltaY
}

func (e *Element) IsHovered() bool {
	return e.isHovered
}

func (e *Element) IsJustHoverIn() bool {
	return e.isJustHoverIn
}

func (e *Element) IsJustHoverOut() bool {
	return e.isJustHoverOut
}

func (e *Element) IsMoving() bool {
	return e.isMoving
}

func (e *Element) IsJustMoveFinished() bool {
	return e.isJustMoveFinished
}

func (e *Element) SetScale(s float64) {
	e.WidthScale = s
	e.HeightScale = s
}

func (e *Element) AbsWidthScale() float64 {
	if e.Parent == nil {
		return e.WidthScale
	}
	return e.Parent.AbsWidthScale() * e.WidthScale
}

func (e *Element) AbsHeightScale() float64 {
	if e.Parent == nil {
		return e.HeightScale
	}
	return e.Parent.AbsHeightScale() * e.HeightScale
}

func (e *Element) BaseWidth() float64 {
	if e.Image == nil {
		return 0
	}
	return float64(e.Image.Bounds().Dx())
}

func (e *Element) BaseHeight() float64 {
	if e.Image == nil {
		return 0
	}
	return float64(e.Image.Bounds().Dy())
}

func (e *Element) AbsWidth() float64 {
	scale := e.AbsWidthScale()
	return e.BaseWidth() * scale
}

func (e *Element) AbsHeight() float64 {
	scale := e.AbsHeightScale()
	return e.BaseHeight() * scale
}

func (e *Element) AbsPos() (float64, float64) {
	if e.Parent == nil {
		return e.XRelativeToParent, e.YRelativeToParent
	}
	px, py := e.Parent.AbsPos()
	parentWsc := e.Parent.AbsWidthScale()
	parentHsc := e.Parent.AbsHeightScale()
	return px + (e.XRelativeToParent * parentWsc), py + (e.YRelativeToParent * parentHsc)
}

func (e *Element) AppendChild(child Component) {
	child.BaseElement().Parent = e
	e.Children = append(e.Children, child)
}

func (e *Element) RemoveChild(target Component) bool {
	index := slices.Index(e.Children, target)
	if index == -1 {
		return false
	}

	// 子要素の親参照を解除
	target.BaseElement().Parent = nil

	// スライスから削除
	e.Children = slices.Delete(e.Children, index, index+1)
	return true
}

func (e *Element) AbsPosLeftOf(target Component, margin float64, align Alignment) (absX, absY float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	absX = tx - e.AbsWidth() - margin
	absY = e.calcVerticalAlign(ty, t.AbsHeight(), align)
	return absX, absY
}

func (e *Element) AbsPosRightOf(target Component, margin float64, align Alignment) (absX, absY float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	absX = tx + t.AbsWidth() + margin
	absY = e.calcVerticalAlign(ty, t.AbsHeight(), align)
	return absX, absY
}

func (e *Element) AbsPosAbove(target Component, margin float64, align Alignment) (absX, absY float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	absX = e.calcHorizontalAlign(tx, t.AbsWidth(), align)
	absY = ty - e.AbsHeight() - margin
	return absX, absY
}

func (e *Element) AbsPosBelow(target Component, margin float64, align Alignment) (absX, absY float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	absX = e.calcHorizontalAlign(tx, t.AbsWidth(), align)
	absY = ty + t.AbsHeight() + margin
	return absX, absY
}

func (e *Element) AbsPosCenterOf(target Component) (absX, absY float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	absX = tx + (t.AbsWidth()-e.AbsWidth())/2
	absY = ty + (t.AbsHeight()-e.AbsHeight())/2
	return absX, absY
}

func (e *Element) PlaceLeftOf(target Component, margin float64, align Alignment) {
	absX, absY := e.AbsPosLeftOf(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceRightOf(target Component, margin float64, align Alignment) {
	absX, absY := e.AbsPosRightOf(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceAbove(target Component, margin float64, align Alignment) {
	absX, absY := e.AbsPosAbove(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceBelow(target Component, margin float64, align Alignment) {
	absX, absY := e.AbsPosBelow(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceCenterOf(target Component) {
	absX, absY := e.AbsPosCenterOf(target)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) AbsXInLayoutWidth(w float64, align Alignment) float64 {
	return e.calcHorizontalAlign(0, w, align)
}

func (e *Element) AbsYInLayoutHeight(h float64, align Alignment) float64 {
	return e.calcVerticalAlign(0, h, align)
}

func (e *Element) AbsPosToRelPosToParent(absX, absY float64) (relX, relY float64) {
	parentX, parentY := 0.0, 0.0
	parentWsc, parentHsc := 1.0, 1.0

	if e.Parent != nil {
		parentX, parentY = e.Parent.AbsPos()
		parentWsc, parentHsc = e.Parent.AbsWidthScale(), e.Parent.AbsHeightScale()
	}
	return (absX - parentX) / parentWsc, (absY - parentY) / parentHsc
}

func (e *Element) calcHorizontalAlign(targetAbsX, targetWidth float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return targetAbsX
	case AlignCenter:
		return targetAbsX + (targetWidth-e.AbsWidth())/2
	case AlignEnd:
		return targetAbsX + targetWidth - e.AbsWidth()
	default:
		return targetAbsX
	}
}

func (e *Element) calcVerticalAlign(targetAbsY, targetHeight float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return targetAbsY
	case AlignCenter:
		return targetAbsY + (targetHeight-e.AbsHeight())/2
	case AlignEnd:
		return targetAbsY + targetHeight - e.AbsHeight()
	default:
		return targetAbsY
	}
}

func (e *Element) Contains(pointX, pointY float64) bool {
	if e.Image == nil {
		return false
	}

	absX, absY := e.AbsPos()
	w, h := e.AbsWidth(), e.AbsHeight()

	isRightOfLeft := pointX >= absX
	isLeftOfRight := pointX < absX+w
	isBelowTop := pointY >= absY
	isAboveBottom := pointY < absY+h

	return isRightOfLeft && isLeftOfRight && isBelowTop && isAboveBottom
}

func (e *Element) Overlaps(other Component) bool {
	t := other.BaseElement()
	if e.Image == nil || t.Image == nil {
		return false
	}

	xa, ya := e.AbsPos()
	wa, ha := e.AbsWidth(), e.AbsHeight()

	xb, yb := t.AbsPos()
	wb, hb := t.AbsWidth(), t.AbsHeight()

	isLeftOfBRight := xa < xb+wb
	isRightOfBLeft := xa+wa > xb
	isAboveBBottom := ya < yb+hb
	isBelowBTop := ya+ha > yb

	return isLeftOfBRight && isRightOfBLeft && isAboveBBottom && isBelowBTop
}

func (e *Element) StartDrag() {
	if !e.Draggable || !e.Enabled {
		return
	}

	e.isDragging = true
	cursorX, cursorY := ebiten.CursorPosition()
	absX, absY := e.AbsPos()
	absWsc, absHsc := e.AbsWidthScale(), e.AbsHeightScale()
	e.dragOffsetX = (float64(cursorX) - absX) / absWsc
	e.dragOffsetY = (float64(cursorY) - absY) / absHsc
}

func (e *Element) StopDrag() {
	e.isDragging = false
	e.DragDeltaX = 0
	e.DragDeltaY = 0
}

func (e *Element) MoveTo(x, y float64) {
	e.toX = x
	e.toY = y
	e.isMoving = true
	e.isJustMoveFinished = false
}

func (e *Element) Update() {
	if !e.Visible {
		return
	}

	e.isJustMoveFinished = false
	if e.isMoving {
		if e.EasingFunc == nil {
			e.XRelativeToParent = e.toX
			e.YRelativeToParent = e.toY
		} else {
			e.XRelativeToParent = e.EasingFunc(e.XRelativeToParent, e.toX)
			e.YRelativeToParent = e.EasingFunc(e.YRelativeToParent, e.toY)
		}

		if e.XRelativeToParent == e.toX && e.YRelativeToParent == e.toY {
			e.isMoving = false
			e.isJustMoveFinished = true
		}
	}

	for _, child := range e.Children {
		child.Update()
	}
}

func (e *Element) Draw(screen *ebiten.Image) {
	if !e.Visible {
		return
	}

	absX, absY := e.AbsPos()
	absWsc, absHsc := e.AbsWidthScale(), e.AbsHeightScale()

	if e.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(absWsc, absHsc)
		op.GeoM.Translate(absX, absY)
		op.Filter = e.Filter
		screen.DrawImage(e.Image, op)
	}

	e.Children.SortByZAsc()
	for _, child := range e.Children {
		child.Draw(screen)
	}
}

func (e *Element) BaseElement() *Element {
	return e
}