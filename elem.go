// その都度Zソートしているため、パフォーマンス改善をする必要あり
package tendon

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	// カプセル化の検討
	Children Elements
	Parent   *Element
}

func NewElement() *Element {
	e := &Element{
		Visible: true,
		Enabled: true,
	}
	e.SetScale(1.0)
	return e
}

func NewButton(relX, relY float64, w, h int, label string, bgColor color.Color) *Element {
	bg := ebiten.NewImage(w, h)
	bg.Fill(bgColor)

	btn := NewElement()
	btn.XRelativeToParent = relX
	btn.YRelativeToParent = relY
	btn.Image = bg

	txtImg := ebiten.NewImage(w, h)
	ebitenutil.DebugPrintAt(txtImg, label, 10, h/2-8)

	textElem := NewElement()
	textElem.Image = txtImg
	textElem.Filter = ebiten.FilterLinear
	textElem.PassThrough = true
	btn.AppendChild(textElem)

	return btn
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

func (e *Element) SetScale(s float64) {
	e.WidthScale = s
	e.HeightScale = s
}

func (e *Element) AbsScale() (float64, float64) {
	if e.Parent == nil {
		return e.WidthScale, e.HeightScale
	}
	parentWs, parentHs := e.Parent.AbsScale()
	return parentWs * e.WidthScale, parentHs * e.HeightScale
}

func (e *Element) Width() float64 {
	if e.Image == nil {
		return 0
	}
	absWs, _ := e.AbsScale()
	return float64(e.Image.Bounds().Dx()) * absWs
}

func (e *Element) Height() float64 {
	if e.Image == nil {
		return 0
	}
	_, absHs := e.AbsScale()
	return float64(e.Image.Bounds().Dy()) * absHs
}

func (e *Element) AbsPos() (float64, float64) {
	if e.Parent == nil {
		return e.XRelativeToParent, e.YRelativeToParent
	}
	px, py := e.Parent.AbsPos()
	parentWs, parentHs := e.Parent.AbsScale()

	return px + (e.XRelativeToParent * parentWs), py + (e.YRelativeToParent * parentHs)
}

func (e *Element) AppendChild(child *Element) {
	child.Parent = e
	e.Children = append(e.Children, child)
}

func (e *Element) AbsPosLeftOf(target *Element, margin float64, align Alignment) (absX, absY float64) {
	tx, ty := target.AbsPos()
	absX = tx - e.Width() - margin
	absY = e.calcVerticalAlign(ty, target.Height(), align)
	return absX, absY
}

func (e *Element) AbsPosRightOf(target *Element, margin float64, align Alignment) (absX, absY float64) {
	tx, ty := target.AbsPos()
	absX = tx + target.Width() + margin
	absY = e.calcVerticalAlign(ty, target.Height(), align)
	return absX, absY
}

func (e *Element) AbsPosAbove(target *Element, margin float64, align Alignment) (absX, absY float64) {
	tx, ty := target.AbsPos()
	absX = e.calcHorizontalAlign(tx, target.Width(), align)
	absY = ty - e.Height() - margin
	return absX, absY
}

func (e *Element) AbsPosBelow(target *Element, margin float64, align Alignment) (absX, absY float64) {
	tx, ty := target.AbsPos()
	absX = e.calcHorizontalAlign(tx, target.Width(), align)
	absY = ty + target.Height() + margin
	return absX, absY
}

func (e *Element) AbsPosCenterOf(target *Element) (absX, absY float64) {
	tx, ty := target.AbsPos()
	absX = tx + (target.Width()-e.Width())/2
	absY = ty + (target.Height()-e.Height())/2
	return absX, absY
}

func (e *Element) AbsXInScreenWidth(w float64, align Alignment) float64 {
	return e.calcHorizontalAlign(0, w, align)
}

func (e *Element) AbsYInScreenHeight(h float64, align Alignment) float64 {
	return e.calcVerticalAlign(0, h, align)
}

func (e *Element) AbsPosToRelPosToParent(absX, absY float64) (relX, relY float64) {
	parentX, parentY := 0.0, 0.0
	parentWs, parentHs := 1.0, 1.0

	if e.Parent != nil {
		parentX, parentY = e.Parent.AbsPos()
		parentWs, parentHs = e.Parent.AbsScale()
	}
	return (absX - parentX) / parentWs, (absY - parentY) / parentHs
}

func (e *Element) calcHorizontalAlign(targetAbsX, targetWidth float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return targetAbsX
	case AlignCenter:
		return targetAbsX + (targetWidth-e.Width())/2
	case AlignEnd:
		return targetAbsX + targetWidth - e.Width()
	default:
		return targetAbsX
	}
}

func (e *Element) calcVerticalAlign(targetAbsY, targetHeight float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return targetAbsY
	case AlignCenter:
		return targetAbsY + (targetHeight-e.Height())/2
	case AlignEnd:
		return targetAbsY + targetHeight - e.Height()
	default:
		return targetAbsY
	}
}

func (e *Element) PlaceLeftOf(target *Element, margin float64, align Alignment) {
	absX, absY := e.AbsPosLeftOf(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceRightOf(target *Element, margin float64, align Alignment) {
	absX, absY := e.AbsPosRightOf(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceAbove(target *Element, margin float64, align Alignment) {
	absX, absY := e.AbsPosAbove(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceBelow(target *Element, margin float64, align Alignment) {
	absX, absY := e.AbsPosBelow(target, margin, align)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) PlaceCenterOf(target *Element) {
	absX, absY := e.AbsPosCenterOf(target)
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) Contains(pointX, pointY float64) bool {
	if e.Image == nil {
		return false
	}

	absX, absY := e.AbsPos()
	w, h := e.Width(), e.Height()

	isRightOfLeft := pointX >= absX
	isLeftOfRight := pointX < absX+w
	isBelowTop := pointY >= absY
	isAboveBottom := pointY < absY+h

	return isRightOfLeft && isLeftOfRight && isBelowTop && isAboveBottom
}

func (e *Element) Overlaps(other *Element) bool {
	if e.Image == nil || other.Image == nil {
		return false
	}

	xa, ya := e.AbsPos()
	wa, ha := e.Width(), e.Height()

	xb, yb := other.AbsPos()
	wb, hb := other.Width(), other.Height()

	isLeftOfBRight := xa < xb+wb
	isRightOfBLeft := xa+wa > xb
	isAboveBBottom := ya < yb+hb
	isBelowBTop := ya+ha > yb

	return isLeftOfBRight && isRightOfBLeft && isAboveBBottom && isBelowBTop
}

func (e *Element) FindAllFromPoint(pointX, pointY float64, dst *Elements) {
	// 子要素も弾く
	if !e.Visible {
		return
	}

	// 親よりも子要素を先に判定
	e.Children.SortByZDesc()
	for _, child := range e.Children {
		child.FindAllFromPoint(pointX, pointY, dst)
	}

	// 最後に自分自身を判定
	if !e.PassThrough && e.Contains(pointX, pointY) {
		*dst = append(*dst, e)
	}
}

func (e *Element) FindAllOverlapping(target *Element, dst *Elements) {
	if !e.Visible {
		return
	}

	// 子要素を先に判定
	e.Children.FindAllOverlapping(target, dst)

	// 自分自身を判定
	if e != target && !e.PassThrough && e.Overlaps(target) {
		*dst = append(*dst, e)
	}
}

func (e *Element) StartDrag() {
	if !e.Draggable || !e.Enabled {
		return
	}
	e.isDragging = true
	cursorX, cursorY := ebiten.CursorPosition()
	absX, absY := e.AbsPos()
	absWs, absHs := e.AbsScale()
	e.dragOffsetX = (float64(cursorX) - absX) / absWs
	e.dragOffsetY = (float64(cursorY) - absY) / absHs
}

func (e *Element) StopDrag() {
	e.isDragging = false
	e.DragDeltaX = 0
	e.DragDeltaY = 0
}

func (e *Element) Update() {
	if !e.Visible {
		return
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
	absWs, absHs := e.AbsScale()

	if e.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(absWs, absHs)
		op.GeoM.Translate(absX, absY)
		op.Filter = e.Filter
		screen.DrawImage(e.Image, op)
	}

	e.Children.SortByZAsc()
	for _, child := range e.Children {
		child.Draw(screen)
	}
}

type Elements []*Element

func (es Elements) SortByZAsc() {
	slices.SortStableFunc(es, func(a, b *Element) int {
		return a.Z - b.Z
	})
}

func (es Elements) SortByZDesc() {
	slices.SortStableFunc(es, func(a, b *Element) int {
		return b.Z - a.Z
	})
}

func (es Elements) FindAllHitTest(pointX, pointY float64, dst *Elements) {
	es.SortByZDesc()
	for _, e := range es {
		e.FindAllFromPoint(pointX, pointY, dst)
	}
}

func (es Elements) FindOverlapping(target *Element, dst *Elements) {
	for _, e := range es {
		if !e.Visible || e.PassThrough || e == target {
			continue
		}
		if e.Overlaps(target) {
			*dst = append(*dst, e)
		}
	}
}

func (es Elements) FindAllOverlapping(target *Element, dst *Elements) {
	for _, e := range es {
		e.FindAllOverlapping(target, dst)
	}
}

func (es Elements) StopAllDrag() {
	for _, e := range es {
		e.StopDrag()
		e.Children.StopAllDrag()
	}
}

func (es Elements) UpdateHover(hitTest Elements) {
	// アロケーションを減らすような設計に変える
	hitSet := make(map[*Element]bool)
	for _, e := range hitTest {
		hitSet[e] = true
	}

	var update func(Elements)
	update = func(elements Elements) {
		for _, e := range elements {
			isHit := hitSet[e] && e.Enabled
			e.isJustHoverIn = false
			e.isJustHoverOut = false

			// 現在ヒットテスト かつ 直前のフレームでは範囲外 (新しく侵入)
			if isHit && !e.isHovered {
				e.isHovered = true
				e.isJustHoverIn = true
				// 現在ヒットテストではない かつ 直前のフレームでは範囲内（今抜けた)
			} else if !isHit && e.isHovered {
				e.isHovered = false
				e.isJustHoverOut = true
			}
			update(e.Children)
		}
	}
	update(es)
}

func (es Elements) UpdateDragMove() {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float64(cursorX), float64(cursorY)

	for _, e := range es {
		if !e.isDragging {
			continue
		}

		absWs, absHs := e.AbsScale()
		toAbsX := cursorXf - (e.dragOffsetX * absWs)
		toAbsY := cursorYf - (e.dragOffsetY * absHs)
		// 目標の絶対座標を、自身の相対座標に変換
		toRelX, toRelY := e.AbsPosToRelPosToParent(toAbsX, toAbsY)

		e.DragDeltaX = toRelX - e.XRelativeToParent
		e.DragDeltaY = toRelY - e.YRelativeToParent

		if !e.ManualDrag {
			e.XRelativeToParent = toRelX
			e.YRelativeToParent = toRelY
		}
	}
}

func (es Elements) Draw(screen *ebiten.Image) {
	// Zが小さい順から描写する
	es.SortByZAsc()
	for _, e := range es {
		e.Draw(screen)
	}
}
