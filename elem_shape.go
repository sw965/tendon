package tendon

import (
	"image/color"
	"math"
)

type Shape interface {
	// lx, lyはローカル座標
	// ローカル座標の定義は、elem_coord.goのElement.PointToLocalPosのコメントを参照
	Contains(lx, ly float64, e *Element) bool
}

type Rect struct{}

func (s Rect) Contains(lx, ly float64, e *Element) bool {
	w, h := e.BaseWidth(), e.BaseHeight()
	// ローカル座標における画像の左右上下の座標を求める
	// 例： 画像の幅(w) = 100 e.AnchorX = 0.5とする
	// left = -100 * 0.5 = -50
	// right = 100 * (1 - 0.5) = 50
	// これはローカル座標における左端と右端の座標と一致する
	// イメージ: -50(左端), -49, -48, ..., 0, 1, 2, ..., 50(右端)
	// topやbottomも同じ要領で求める
	left := -w * e.AnchorX
	right := w * (1 - e.AnchorX)
	top := -h * e.AnchorY
	bottom := h * (1 - e.AnchorY)

	// 点(x座標) が 画像の左端より右にあるか？
	withinLeft := lx >= left
	// 点(x座標) が 画像の右端より左にあるか？
	withinRight := lx <= right
	// 点(y座標) が 画像の上端より下にあるか？
	withinTop := ly >= top
	// 点(y座標) が 画像の下端より上にあるか？
	withinBottom := ly <= bottom

	/*
		        判定のイメージ図
		        点 が left   よりも 右 かつ
				点 が right  よりも 左 かつ
				点 が top    よりも 下 かつ
				点 が bottom よりも 上

		                  top
		           +-------^-------+
		           |               |
		      left <    (lx, ly)   > right
		           |               |
		           |               |
		           +-------v-------+
		                 bottom
	*/
	return withinLeft && withinRight && withinTop && withinBottom
}

type Circle struct {
	Radius float64
}

func (s Circle) Contains(lx, ly float64, e *Element) bool {
	r := s.Radius
	if r <= 0 {
		r = math.Min(e.BaseWidth(), e.BaseHeight()) / 2
	}

	// ローカル座標における画像の中心を求める
	// 例1： e.BaseWidth() = e.BaseHeight() = 100, e.AnchorX = e.AnchorY = 0.5 であるとき
	// 画像の中心 は (50, 50)、アンカー位置 も (50, 50)
	// よって centerX = centerY = (100 / 2) - (100 * 0.5) = 0
	// これはアンカーを原点としたとき、画像の中心 は (0, 0) である事を意味する
	//
	// 例2： e.BaseWidth() = e.BaseHeight() = 100, e.AnchorX = e.AnchorY = 0.3 であるとき
	// 画像の中心 は (50, 50)、アンカー位置 は (30, 30)
	// よって centerX = centerY = (100 / 2) - (100 * 0.3) = 20
	// これはアンカーを原点としたとき、画像の中心 は (20, 20) に位置する事を意味する。
	centerX := (e.BaseWidth() / 2.0) - (e.BaseWidth() * e.AnchorX)
	centerY := (e.BaseHeight() / 2.0) - (e.BaseHeight() * e.AnchorY)

	// 入力された点 (lx, ly) と 円(画像)の中心のズレを計算
	dx, dy := lx-centerX, ly-centerY

	// 【当たり判定のイメージ】
	// 1. 円の中心から、入力された点 (lx, ly) に向かって1本の直線を引きます。
	//    この直線の長さを「 k 」とする。
	//
	// 2. 次に、この直線 k を「斜辺」とする直角三角形を思い浮かべる。
	//    円の中心から真横に「 dx 」だけ進み、そこから真縦に「 dy 」だけ進むと、
	//    入力された点 (lx, ly) にぴったり到着して、直角三角形が出来る。
	//
	// 3. ピタゴラスの定理（底辺の2乗 + 高さの2乗 = 斜辺の2乗）により、
	//    (dx * dx) + (dy * dy) は、斜辺 k の2乗（中心からの距離の2乗）になる。
	//
	// 4. この 「k の2乗」が、円の「半径の2乗 (r * r)」以下であれば、
	//    その点 (lx, ly) は円の内側に入力されたと判定出来る。
	return (dx*dx + dy*dy) <= r*r
}

// tendon/elem_shape.go

func (e *Element) SetRectApprox(cols, rows int, overlapRatio, protrusionRatio float64) error {
	e.Shape = Rect{}

	e.rebuildCollider = func() {
		wScale := e.AbsWidthScale()
		hScale := e.AbsHeightScale()
		scaledW := e.BaseWidth() * wScale
		scaledH := e.BaseHeight() * hScale

		// ★ 動的分割の適用 (もし cols や rows が 0 以下の場合は自動計算するなどの仕様がおすすめ)
		// 例: 30ピクセルごとに1つ円を配置する
		actualCols := cols
		actualRows := rows
		if actualCols <= 0 {
			actualCols = int(math.Max(1, math.Round(scaledW/30.0)))
		}
		if actualRows <= 0 {
			actualRows = int(math.Max(1, math.Round(scaledH/30.0)))
		}

		colliders, err := NewRectCircleColliders(
			scaledW, scaledH,
			actualCols, actualRows, overlapRatio, protrusionRatio,
			e.AnchorX, e.AnchorY,
		)

		if err == nil {
			for i := range colliders {
				// ローカル座標の定義により、スケールを元に戻す
				if wScale != 0 {
					colliders[i].LocalX /= wScale
				}
				if hScale != 0 {
					colliders[i].LocalY /= hScale
				}
			}
			e.CircleColliders = colliders
		}
	}

	e.rebuildCollider()
	return nil
}

func (e *Element) AsCircle(borderColor color.Color, borderWidth float32) {
	if e.Image == nil {
		return
	}
	newImg := CreateCircularImage(e.Image, borderColor, borderWidth)
	e.Image.Dispose()
	e.Image = newImg
	e.Shape = &Circle{Radius: 0}

	e.rebuildCollider = func() {
		wScale := e.AbsWidthScale()
		hScale := e.AbsHeightScale()
		scaledW := e.BaseWidth() * wScale
		scaledH := e.BaseHeight() * hScale

		r := math.Min(scaledW, scaledH) / 2
		cx := scaledW/2 - (scaledW * e.AnchorX)
		cy := scaledH/2 - (scaledH * e.AnchorY)

		if wScale != 0 {
			cx /= wScale
		}
		if hScale != 0 {
			cy /= hScale
		}

		e.CircleColliders = []CircleCollider{{LocalX: cx, LocalY: cy, Radius: r}}
	}

	// 初回実行
	e.rebuildCollider()
}
