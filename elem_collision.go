package tendon

import (
	"fmt"
	"math"
)

type CircleCollider struct {
	LocalX float64
	LocalY float64
	Radius float64
}

func NewRectCircleColliders(w, h float64, cols, rows int, overlapRatio, protrusionRatio, anchorX, anchorY float64) ([]CircleCollider, error) {
	if cols <= 0 || rows <= 0 {
		return nil, fmt.Errorf("cols <= 0 || rows <= 0")
	}

	colliders := make([]CircleCollider, 0, cols*rows)
	// w, h には BaseWidth, BaseHeight が渡される想定
	baseW := w / float64(cols)
	baseH := h / float64(rows)

	radius := (math.Max(baseW, baseH) / 2) * (1 + protrusionRatio + overlapRatio)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cx := (float64(c) + 0.5) * baseW
			cy := (float64(r) + 0.5) * baseH

			lx := cx - (w * anchorX)
			ly := cy - (h * anchorY)

			colliders = append(colliders, CircleCollider{
				LocalX: lx,
				LocalY: ly,
				Radius: radius,
			})
		}
	}
	return colliders, nil
}

func (e *Element) Contains(pointX, pointY float64) bool {
	if e.Image == nil {
		return false
	}

	// 回転行列を使ってローカル座標に変換
	lx, ly := e.PointToLocalPos(pointX, pointY)

	// Shapeが設定されていればそのロジックを使う
	if e.Shape != nil {
		return e.Shape.Contains(lx, ly, e)
	}

	// 設定されていなければデフォルトの矩形判定
	return Rect{}.Contains(lx, ly, e)
}

func (e *Element) FindAllFromPoint(pointX, pointY float64, dst *Components) {
	if !e.Visible {
		return
	}

	e.sortChildren()

	// 手前の要素（インデックスが大きい方）から逆順に調べる
	for i := len(e.Children) - 1; i >= 0; i-- {
		child := e.Children[i]
		child.BaseElement().FindAllFromPoint(pointX, pointY, dst)
	}

	// 自分自身の判定
	if !e.PassThrough && e.Contains(pointX, pointY) {
		*dst = append(*dst, e)
	}
}

func (e *Element) Overlaps(other Component) bool {
	t := other.BaseElement()
	if e.Image == nil || t.Image == nil {
		return false
	}

	e.resolveDirtyScale()
	t.resolveDirtyScale()
	
	if len(e.CircleColliders) > 0 && len(t.CircleColliders) > 0 {
		for _, ec := range e.CircleColliders {
			ex, ey := e.LocalPosToAbsPos(ec.LocalX, ec.LocalY)
			for _, tc := range t.CircleColliders {
				tx, ty := t.LocalPosToAbsPos(tc.LocalX, tc.LocalY)

				// 2つの円の中心のx座標とy座標のズレをそれぞれ計算
				dx, dy := ex-tx, ey-ty

				// 【円と円の当たり判定（重なり判定）】
				// 1. 円Aの中心 (ex, ey) と、円Bの中心 (tx, ty) を結ぶ1本の直線を引く。
				//    この直線の長さを「 k 」とする。
				//
				// 2. 次に、この直線 k を「斜辺」とする直角三角形を思い浮かべる。
				//    円Aの中心から真横に「 dx 」だけ進み、そこから真縦に「 dy 」だけ進むと、
				//    円Bの中心に到着し、直角三角形が完成する。
				//
				// 3. ピタゴラスの定理により、(dx * dx) + (dy * dy) = (k の 2乗) になる。
				//    なお (k の 2乗) の 変数名はdistSqである。
				//
				// 4. 2つの円が重なっている状態とは、2つの円の中心を結んだ 直線 k が
				//    2つの円の半径の合計以下のとき (k <= 円Aの半径 + 円Bの半径)
				//    ここでは、kの2乗 <= (円Aの半径 + 円Bの半径)の2乗 で判定する
				distSq := dx*dx + dy*dy
				rSum := ec.Radius + tc.Radius
				if distSq <= rSum*rSum {
					return true
				}
			}
		}
		return false
	}

	// === これ以降は既存のコードそのまま（フォールバック用のAABB判定） ===
	xaMin, yaMin, xaMax, yaMax := e.BoundingBox()
	xbMin, ybMin, xbMax, ybMax := t.BoundingBox()

	isLeftOfBRight := xaMin < xbMax
	isRightOfBLeft := xaMax > xbMin
	isAboveBBottom := yaMin < ybMax
	isBelowBTop := yaMax > ybMin

	return isLeftOfBRight && isRightOfBLeft && isAboveBBottom && isBelowBTop
}

func (e *Element) FindAllOverlapping(target Component, dst *Components) {
	if !e.Visible {
		return
	}

	e.sortChildren()

	// 手前から奥へ再帰探索
	for i := len(e.Children) - 1; i >= 0; i-- {
		e.Children[i].BaseElement().FindAllOverlapping(target, dst)
	}

	// 自身の判定
	if e != target.BaseElement() && !e.PassThrough && e.Overlaps(target) {
		*dst = append(*dst, e)
	}
}
