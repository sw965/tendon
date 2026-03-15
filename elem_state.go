package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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

func (e *Element) StartDrag() {
	if !e.Draggable || !e.Enabled {
		return
	}

	e.isDragging = true
	cursorX, cursorY := ebiten.CursorPosition()
	
	// マウスの絶対座標を、親のローカル座標系での座標に一発で変換
	relMouseX, relMouseY := e.AbsPosToRelPosToParent(float64(cursorX), float64(cursorY))
	
	// 要素の左上(XRelativeToParent)から見た、マウスの相対的なオフセットを記録
	// これにより、要素のどこを掴んでもピタッと吸着します
	e.dragOffsetX = relMouseX - e.XRelativeToParent
	e.dragOffsetY = relMouseY - e.YRelativeToParent
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
