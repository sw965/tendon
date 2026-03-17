package tendon

import (
	"image/color"
	"math"
)

type Shape interface {
	// lx, ly は逆回転適用後の座標
	Contains(lx, ly float64, e *Element) bool
}

type Rect struct{}

func (s Rect) Contains(lx, ly float64, e *Element) bool {
	w, h := e.BaseWidth(), e.BaseHeight()
	left := -w * e.AnchorX
	right := w * (1 - e.AnchorX)
	top := -h * e.AnchorY
	bottom := h * (1 - e.AnchorY)
	return lx >= left && lx <= right && ly >= top && ly <= bottom
}

type Circle struct {
	Radius float64
}

func (s Circle) Contains(lx, ly float64, e *Element) bool {
	r := s.Radius
	if r <= 0 {
		r = math.Min(e.BaseWidth(), e.BaseHeight()) / 2
	}

	// アンカー(0,0)から見た「画像中央」の座標を出す
	cx := e.BaseWidth()/2 - (e.BaseWidth() * e.AnchorX)
	cy := e.BaseHeight()/2 - (e.BaseHeight() * e.AnchorY)

	// 中央からの距離で判定
	dx, dy := lx-cx, ly-cy
	return (dx*dx + dy*dy) <= r*r
}

func (e *Element) SetRectApprox(cols, rows int, overlapRatio, protrusionRatio float64) error {
	e.Shape = Rect{}
	// AbsWidth() ではなく BaseWidth() を渡す！
	// スケールは LocalPosToAbsPos が計算時に自動で掛けてくれるため
	colliders, err := NewRectCircleColliders(e.BaseWidth(), e.BaseHeight(), cols, rows, overlapRatio, protrusionRatio, e.AnchorX, e.AnchorY)
	if err != nil {
		return err
	}
	e.CircleColliders = colliders
	return nil
}

func (e *Element) AsCircle(borderColor color.Color, borderWidth float32) {
	if e.Image == nil {
		return
	}
	e.Image = CreateCircularImage(e.Image, borderColor, borderWidth)
	e.Shape = &Circle{Radius: 0}

	r := math.Min(e.BaseWidth(), e.BaseHeight()) / 2

	// 要素の「本当の中心座標」が、アンカー(原点)から見てどこにあるかを計算する
	cx := e.BaseWidth()/2 - (e.BaseWidth() * e.AnchorX)   // ⭕️ 修正後
	cy := e.BaseHeight()/2 - (e.BaseHeight() * e.AnchorY) // ⭕️ 修正後

	// LocalX: 0, LocalY: 0 ではなく、計算した中心座標をセット
	e.CircleColliders = []CircleCollider{{LocalX: cx, LocalY: cy, Radius: r}}
}
