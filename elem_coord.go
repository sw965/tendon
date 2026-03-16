package tendon

import (
	"math"
)

// AbsPos は要素の「左上」の画面上の絶対座標を返します。
func (e *Element) AbsPos() (float64, float64) {
	m := e.AbsGeoM()
	return m.Apply(0, 0)
}

// AbsPosToRelPosToParent は画面上の絶対座標を、親要素のローカル空間（XRelativeToParentなどに使える座標）に変換します。
func (e *Element) AbsPosToRelPosToParent(absX, absY float64) (relX, relY float64) {
	if e.Parent == nil {
		return absX, absY
	}
	// 親の行列を逆変換(Invert)して通すだけで、親がどんなに回転・拡大していても一発でローカル座標に戻ります
	pm := e.Parent.AbsGeoM()
	pm.Invert()
	return pm.Apply(absX, absY)
}

func (e *Element) SetAbsPos(absX, absY float64) {
	// 内部で RelPos への変換を隠蔽する
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) LocalPosToAbsPos(lx, ly float64) (float64, float64) {
	m := e.AbsGeoM()
	w, h := e.BaseWidth(), e.BaseHeight()
	tx := lx + w*e.AnchorX
	ty := ly + h*e.AnchorY
	return m.Apply(tx, ty)
}

func (e *Element) PointToLocalPos(pointX, pointY float64) (float64, float64) {
	m := e.AbsGeoM()
	m.Invert()
	// (pointX, pointY) を 自身(e) の スケールや回転を元に戻した状態の画像の左上を 原点(0, 0) としたときの相対座標へ変換する
	xRelToImg, yRelToImg := m.Apply(pointX, pointY)
	w, h := e.BaseWidth(), e.BaseHeight()
	// ローカル座標は「スケールや回転を元に戻した状態の画像に対するアンカーを原点とみなした座標」が定義
	// 例えば、画像サイズを w = 200, h = 200, e.AnchorX, e.AnchorY = 0.5 (画像の中心) としたとき
	// 画像から見て、(50, 50) の地点をクリックしたとするならば、
	// lx = 50 - 200 * 0.5 = -50
	// これは、アンカーから-50ズレているX座標をクリックしたことを意味する
	// すなわち、アンカーを原点としたとき、-50のX座標をクリックした事と同義であり、ローカル座標の定義と一致する 
	lx := xRelToImg - w*e.AnchorX
	ly := yRelToImg - h*e.AnchorY
	return lx, ly
}

func (e *Element) BoundingBox() (float64, float64, float64, float64) {
	m := e.AbsGeoM()
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