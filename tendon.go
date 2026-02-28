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

	Image   *ebiten.Image
	Visible bool

	OnLeftClick   func(e *Element)
	OnLeftPressed func(e *Element)
	OnLeftReleased func(e *Element)

	OnRightClick  func(e *Element)
	OnRightPressed func(e *Element)
	OnRightReleased func(e *Element)

	OnMiddleClick func(e *Element)
	OnMiddlePressed func(e *Element)
	OnMiddleReleased func(e *Element)

	OnMouseEnter  func(e *Element)
	OnMouseLeave  func(e *Element)
	OnUpdate      func(e *Element)

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

func (e *Element) Update(parentX, parentY float64) {
	if !e.Visible {
		return
	}

	// 自分の画面上の絶対位置（グローバル座標）
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
		op.GeoM.Translate(absX, absY)
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
	w, h := e.Image.Bounds().Dx(), e.Image.Bounds().Dy()
	return cx >= globalX && cx < globalX+float64(w) && cy >= globalY && cy < globalY+float64(h)
}

func NewButton(relX, relY float64, w, h int, label string, bgColor color.Color) *Element {
	bg := ebiten.NewImage(w, h)
	bg.Fill(bgColor)

	btn := &Element{
		XRelativeToParent: relX, YRelativeToParent: relY,
		Image:   bg,
		Visible: true,
	}

	txtImg := ebiten.NewImage(w, h)
	ebitenutil.DebugPrintAt(txtImg, label, 10, h/2-8)

	btn.Children = append(btn.Children, &Element{
		XRelativeToParent: 0, YRelativeToParent: 0,
		Image:   txtImg,
		Visible: true,
	})

	return btn
}