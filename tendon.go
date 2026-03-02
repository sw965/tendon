package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"slices"
	"fmt"
)

var DebugMode = false

type Element struct {
	Name string

	// 親の左上を基準とした相対座標
	XRelativeToParent float64
	YRelativeToParent float64

	Image       *ebiten.Image
	WidthScale  float64
	HeightScale float64
	Filter      ebiten.Filter
	Z           int
	Visible     bool // trueなら描写される
	Enabled     bool // trueなら操作可能
	PassThrough bool

	OnLeftClick    func(e *Element)
	OnLeftPressed  func(e *Element)
	OnLeftReleased func(e *Element)

	OnRightClick    func(e *Element)
	OnRightPressed  func(e *Element)
	OnRightReleased func(e *Element)

	OnMiddleClick    func(e *Element)
	OnMiddlePressed  func(e *Element)
	OnMiddleReleased func(e *Element)

	OnMouseEnter func(e *Element)
	OnMouseLeave func(e *Element)
	OnDrag       func(e *Element)
	OnUpdate     func(e *Element)

	Draggable   bool
	ManualDrag  bool
	isDragging  bool
	dragOffsetX float64
	dragOffsetY float64
	DragDeltaX  float64
	DragDeltaY  float64

	// カプセル化の検討
	Children Elements
	Parent   *Element

	// 内部状態
	isHovered bool
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

func (e *Element) SetScale(s float64) {
	e.WidthScale = s
	e.HeightScale = s
}

func (e *Element) AppendChild(child *Element) {
	child.Parent = e
	e.Children = append(e.Children, child)
}

func (e *Element) AbsPos() (x, y float64) {
	if e.Parent == nil {
		return e.XRelativeToParent, e.YRelativeToParent
	}
	px, py := e.Parent.AbsPos()
	return px + e.XRelativeToParent, py + e.YRelativeToParent
}

func (e *Element) Update(parentX, parentY float64, target *Element) {
	if !e.Visible {
		return
	}

	isTarget := (e == target)
	isTargetAndEnabled := isTarget && e.IsEnabled()

	// 自分の画面上の絶対位置を計算する
	// 親要素の座標(parentX, parentY)が動けば、子要素の絶対座標も動く
	absX := parentX + e.XRelativeToParent
	absY := parentY + e.YRelativeToParent

	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float64(cursorX), float64(cursorY)

	isCtrlDragging := DebugMode && ebiten.IsKeyPressed(ebiten.KeyControl)
	isAltDragging := DebugMode && ebiten.IsKeyPressed(ebiten.KeyAlt)
	isDebugDragging := isCtrlDragging || isAltDragging
	isDraggable := e.Draggable || isDebugDragging

	if isDraggable {
		// ドラッグ開始
		if isTargetAndEnabled && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.isDragging = true
			e.dragOffsetX = cursorXf - absX
			e.dragOffsetY = cursorYf - absY
		}

		if e.isDragging {
			oldX, oldY := e.XRelativeToParent, e.YRelativeToParent
			targetX := (cursorXf - e.dragOffsetX) - parentX
			targetY := (cursorYf - e.dragOffsetY) - parentY
			e.DragDeltaX = targetX - oldX
			e.DragDeltaY = targetY - oldY

			// 自動で動かす
			if !e.ManualDrag {
				e.XRelativeToParent = targetX
				e.YRelativeToParent = targetY
			}

			if e.OnDrag != nil {
				e.OnDrag(e)
			}

			// 離したら終了
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				e.isDragging = false
				e.DragDeltaX = 0
				e.DragDeltaY = 0

				if DebugMode {
					fmt.Printf("Element [%s] -> XRelativeToParent: %.2f, YRelativeToParent: %.2f\n", e.Name, e.XRelativeToParent, e.YRelativeToParent)
				}
			}
		}
	}

	// 範囲外だったカーソルが重なる度に、処理を実行する (!e.isHoverdであれば、直前のフレームは範囲外であった事を意味する)
	if isTargetAndEnabled && !e.isHovered {
		e.isHovered = true
		if e.OnMouseEnter != nil {
			e.OnMouseEnter(e)
		}
		// 重なったカーソルが範囲外に移動する度に、処理を実行する (e.isHoverdであれば、直前のフレームは範囲内であった事を意味する)
	} else if !isTargetAndEnabled && e.isHovered {
		e.isHovered = false
		if e.OnMouseLeave != nil {
			e.OnMouseLeave(e)
		}
	}

	if isTargetAndEnabled && !isDebugDragging {
		// 左クリック関連
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && e.OnLeftClick != nil {
			e.OnLeftClick(e)
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && e.OnLeftPressed != nil {
			e.OnLeftPressed(e)
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && e.OnLeftReleased != nil {
			e.OnLeftReleased(e)
		}

		// 右クリック関連
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && e.OnRightClick != nil {
			e.OnRightClick(e)
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && e.OnRightPressed != nil {
			e.OnRightPressed(e)
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) && e.OnRightReleased != nil {
			e.OnRightReleased(e)
		}

		// ホイールクリック関連
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) && e.OnMiddleClick != nil {
			e.OnMiddleClick(e)
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) && e.OnMiddlePressed != nil {
			e.OnMiddlePressed(e)
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) && e.OnMiddleReleased != nil {
			e.OnMiddleReleased(e)
		}
	}

	// Updateが呼び出される度に処理する
	if e.OnUpdate != nil {
		e.OnUpdate(e)
	}

	// ドラッグ等で相対座標が変更された可能性があるため、子要素に渡す前に絶対座標を再計算する
	absX = parentX + e.XRelativeToParent
	absY = parentY + e.YRelativeToParent

	// 子要素の更新
	for _, child := range e.Children {
		child.Update(absX, absY, target)
	}
}

// Draw は相対座標を計算しながら描画します
func (e *Element) Draw(screen *ebiten.Image) {
	e.DrawAt(screen, 0, 0)
}

func (e *Element) DrawAt(screen *ebiten.Image, parentX, parentY float64) {
	if !e.Visible {
		return
	}

	// 自分の画面上の絶対位置を計算する
	// 親要素の座標(parentX, parentY)が動けば、子要素の絶対座標も動く
	absX := parentX + e.XRelativeToParent
	absY := parentY + e.YRelativeToParent

	// 親要素のZに関係なく子要素よりも先に描写する (重なる場合、親要素は子要素よりも下に描写される)
	if e.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(e.WidthScale, e.HeightScale)
		op.GeoM.Translate(absX, absY)
		op.Filter = e.Filter
		screen.DrawImage(e.Image, op)
	}

	// 子要素同士では、Zの小さい順から描写する
	e.Children.SortByZAsc()
	for _, child := range e.Children {
		child.DrawAt(screen, absX, absY)
	}
}

func (e *Element) Contains(absX, absY, px, py float64) bool {
	if e.Image == nil {
		return false
	}
	w, h := e.Width(), e.Height()
	return px >= absX && px < absX+w && py >= absY && py < absY+h
}

func (e *Element) FindDeepestFromPoint(parentX, parentY, px, py float64) *Element {
	if !e.Visible {
		return nil
	}

	if e.IsPassThrough() {
		return nil
	}

	// 自分の画面上の絶対位置を計算する
	// 親要素の座標(parentX, parentY)が動けば、子要素の絶対座標も動く
	absX := parentX + e.XRelativeToParent
	absY := parentY + e.YRelativeToParent

	// 子要素から先に判定する (子要素は親要素よりも先に描写されるため)
	e.Children.SortByZDesc() // Zが大きい順に並び変える (手前に描写される順)
	for _, child := range e.Children {
		if target := child.FindDeepestFromPoint(absX, absY, px, py); target != nil {
			return target
		}
	}

	// 子要素がヒットしなかったら、自分自身を判定する
	if e.Contains(absX, absY, px, py) {
		return e
	}
	return nil
}

func (e *Element) Width() float64 {
	if e.Image != nil {
		return float64(e.Image.Bounds().Dx()) * e.WidthScale
	}
	return 0
}

func (e *Element) Height() float64 {
	if e.Image != nil {
		return float64(e.Image.Bounds().Dy()) * e.HeightScale
	}
	return 0
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

func (e *Element) AbsPosToRelPosToParent(absX, absY float64) (relX, relY float64) {
	px, py := 0.0, 0.0
	if e.Parent != nil {
		// 自分の親 (e.Parent) の絶対座標
		px, py = e.Parent.AbsPos()
	}
	return absX - px, absY - py
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

func (e *Element) IsEnabled() bool {
	isDebugDragging := DebugMode && (ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyAlt))
	return e.Enabled || isDebugDragging
}

func (e *Element) IsPassThrough() bool {
	isDebugDragging := DebugMode && ebiten.IsKeyPressed(ebiten.KeyAlt)
	return e.PassThrough && !isDebugDragging
}

type Elements []*Element

func (es Elements) FindDeepestFromDragging() *Element {
	for _, e := range es {
		if e.isDragging {
			return e
		}
		if found := e.Children.FindDeepestFromDragging(); found != nil {
			return found
		}
	}
	return nil
}

func (es Elements) SortByZAsc() {
	slices.SortFunc(es, func(a, b *Element) int {
		return a.Z - b.Z
	})
}

func (es Elements) SortByZDesc() {
	slices.SortFunc(es, func(a, b *Element) int {
		return b.Z - a.Z
	})
}

func (es Elements) Update(parentX, parentY float64) {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float64(cursorX), float64(cursorY)

	// ドラッグ中の要素があれば、その要素を最優先ターゲットとする
	target := es.FindDeepestFromDragging()
	es.SortByZDesc()

	if target == nil {
		for _, e := range es {
			if found := e.FindDeepestFromPoint(parentX, parentY, cursorXf, cursorYf); found != nil {
				target = found
				break
			}
		}
	}

	for _, e := range es {
		e.Update(parentX, parentY, target)
	}
}

func (es Elements) Draw(screen *ebiten.Image) {
	// Zが小さい順から描写する
	es.SortByZAsc()
	for _, e := range es {
		e.Draw(screen)
	}
}

type Alignment int

const (
	AlignStart  Alignment = iota // 上端 または 左端
	AlignCenter                  // 中央
	AlignEnd                     // 下端 または 右端
)
