package tendon

import (
	"math"
)

// AbsPos は要素の「左上」の画面上の絶対座標を返します。
func (e *Element) AbsPos() (float64, float64) {
	m := e.TransformMatrix()
	return m.Apply(0, 0)
}

// AbsPosToRelPosToParent は画面上の絶対座標を、親要素のローカル空間（XRelativeToParentなどに使える座標）に変換します。
func (e *Element) AbsPosToRelPosToParent(absX, absY float64) (relX, relY float64) {
	if e.Parent == nil {
		return absX, absY
	}
	// 親の行列を逆変換(Invert)して通すだけで、親がどんなに回転・拡大していても一発でローカル座標に戻ります
	pm := e.Parent.TransformMatrix()
	pm.Invert()
	return pm.Apply(absX, absY)
}

func (e *Element) SetAbsPos(absX, absY float64) {
	// 内部で RelPos への変換を隠蔽する
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) LocalPosToAbsPos(lx, ly float64) (float64, float64) {
	m := e.TransformMatrix()
	w, h := e.BaseWidth(), e.BaseHeight()
	tx := lx + w*e.AnchorX
	ty := ly + h*e.AnchorY
	return m.Apply(tx, ty)
}

func (e *Element) PointToLocalPos(pointX, pointY float64) (float64, float64) {
	m := e.TransformMatrix()
	m.Invert()
	tx, ty := m.Apply(pointX, pointY)
	w, h := e.BaseWidth(), e.BaseHeight()
	return tx - w*e.AnchorX, ty - h*e.AnchorY
}

func (e *Element) BoundingBox() (float64, float64, float64, float64) {
	m := e.TransformMatrix()
	w, h := e.BaseWidth(), e.BaseHeight()

	p1x, p1y := m.Apply(0, 0)
	p2x, p2y := m.Apply(w, 0)
	p3x, p3y := m.Apply(0, h)
	p4x, p4y := m.Apply(w, h)

	minX := math.Min(math.Min(p1x, p2x), math.Min(p3x, p4x))
	maxX := math.Max(math.Max(p1x, p2x), math.Max(p3x, p4x))
	minY := math.Min(math.Min(p1y, p2y), math.Min(p3y, p4y))
	maxY := math.Max(math.Max(p1y, p2y), math.Max(p3y, p4y))

	return minX, minY, maxX, maxY
}