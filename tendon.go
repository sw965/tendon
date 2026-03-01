package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type Element struct {
	// 親の左上を基準とした相対座標
	XRelativeToParent float64
	YRelativeToParent float64

	Image       *ebiten.Image
	WidthScale  float64
	HeightScale float64
	Filter      ebiten.Filter
	Visible     bool

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
	OnUpdate     func(e *Element)

	// ドラッグ機能
	Draggable   bool
	isDragging  bool
	dragOffsetX float64
	dragOffsetY float64

	// 子要素
	Children []*Element

	// 内部状態
	isHovered bool
}

func NewElement() *Element {
	return &Element{
		WidthScale:  1.0,
		HeightScale: 1.0,
		Visible:     true,
	}
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
	btn.Children = append(btn.Children, textElem)

	return btn
}

func (e *Element) Update(parentX, parentY float64) {
	if !e.Visible {
		return
	}

	// 自分の画面上の絶対位置
	absX := parentX + e.XRelativeToParent
	absY := parentY + e.YRelativeToParent

	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float64(cursorX), float64(cursorY)
	// 引数の順番を検討
	isCursorInside := e.isHit(cursorXf, cursorYf, absX, absY)

	if e.Draggable {
		if isCursorInside && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.isDragging = true
			e.dragOffsetX = cursorXf - absX
			e.dragOffsetY = cursorYf - absY
		}

		// ドラッグ中
		if e.isDragging {
			e.XRelativeToParent = (cursorXf - e.dragOffsetX) - parentX
			e.YRelativeToParent = (cursorYf - e.dragOffsetY) - parentY

			// 左クリックを離したらドラッグ終了
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				e.isDragging = false
			}
		}
	}

	// カーソルが早く動いてホバー判定が外れないように、ドラッグ中はホバー中と見なす
	isCursorInsideOrDragging := isCursorInside || e.isDragging

	// 範囲外だったカーソルが重なる度に、処理を実行する
	if isCursorInsideOrDragging && !e.isHovered {
		e.isHovered = true
		if e.OnMouseEnter != nil {
			e.OnMouseEnter(e)
		}
		// 重なったカーソルが範囲外に移動する度に、処理を実行する
	} else if !isCursorInsideOrDragging && e.isHovered {
		e.isHovered = false
		if e.OnMouseLeave != nil {
			e.OnMouseLeave(e)
		}
	}

	if isCursorInside {
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

	// 子要素の更新
	for _, child := range e.Children {
		child.Update(absX, absY)
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

	absX := parentX + e.XRelativeToParent
	absY := parentY + e.YRelativeToParent

	if e.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(e.WidthScale, e.HeightScale)
		op.GeoM.Translate(absX, absY)
		op.Filter = e.Filter
		screen.DrawImage(e.Image, op)
	}

	for _, child := range e.Children {
		child.DrawAt(screen, absX, absY)
	}
}

// isHit は絶対座標(globalX, globalY)を基準に当たり判定を行います
// 後で命名や引数の順番を検討する
func (e *Element) isHit(cx, cy, globalX, globalY float64) bool {
	if e.Image == nil {
		return false
	}
	w, h := e.Width(), e.Height()
	return cx >= globalX && cx < globalX+w && cy >= globalY && cy < globalY+h
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

func (e *Element) XYLeftOf(target *Element, margin float64, align Alignment) (x, y float64) {
	x = target.XRelativeToParent - e.Width() - margin
	y = e.calcVerticalAlign(target, align)
	return x, y
}

func (e *Element) XYRightOf(target *Element, margin float64, align Alignment) (x, y float64) {
	x = target.XRelativeToParent + target.Width() + margin
	y = e.calcVerticalAlign(target, align)
	return x, y
}

func (e *Element) XYAbove(target *Element, margin float64, align Alignment) (x, y float64) {
	x = e.calcHorizontalAlign(target, align)
	y = target.YRelativeToParent - e.Height() - margin
	return x, y
}

func (e *Element) XYBelow(target *Element, margin float64, align Alignment) (x, y float64) {
	x = e.calcHorizontalAlign(target, align)
	y = target.YRelativeToParent + target.Height() + margin
	return x, y
}

func (e *Element) calcHorizontalAlign(target *Element, align Alignment) float64 {
	switch align {
	case AlignStart:
		return target.XRelativeToParent
	case AlignCenter:
		return target.XRelativeToParent + (target.Width()-e.Width())/2
	case AlignEnd:
		return target.XRelativeToParent + target.Width() - e.Width()
	default:
		return target.XRelativeToParent
	}
}

func (e *Element) calcVerticalAlign(target *Element, align Alignment) float64 {
	switch align {
	case AlignStart:
		return target.YRelativeToParent
	case AlignCenter:
		return target.YRelativeToParent + (target.Height()-e.Height())/2
	case AlignEnd:
		return target.YRelativeToParent + target.Height() - e.Height()
	default:
		return target.YRelativeToParent
	}
}

func (e *Element) PlaceRightOf(target *Element, margin float64, align Alignment) {
	e.XRelativeToParent, e.YRelativeToParent = e.XYRightOf(target, margin, align)
}

func (e *Element) PlaceLeftOf(target *Element, margin float64, align Alignment) {
	e.XRelativeToParent, e.YRelativeToParent = e.XYLeftOf(target, margin, align)
}

func (e *Element) PlaceBelow(target *Element, margin float64, align Alignment) {
	e.XRelativeToParent, e.YRelativeToParent = e.XYBelow(target, margin, align)
}

func (e *Element) PlaceAbove(target *Element, margin float64, align Alignment) {
	e.XRelativeToParent, e.YRelativeToParent = e.XYAbove(target, margin, align)
}

type Alignment int

const (
	AlignStart  Alignment = iota // 上端 または 左端
	AlignCenter                  // 中央
	AlignEnd                     // 下端 または 右端
)
