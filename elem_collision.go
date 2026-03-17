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

	if len(e.CircleColliders) > 0 && len(t.CircleColliders) > 0 {
		// ⭕️ 追加: 双方のスケール値を取得（幅と高さで大きい方を採用）
		eScale := math.Max(e.AbsWidthScale(), e.AbsHeightScale())
		tScale := math.Max(t.AbsWidthScale(), t.AbsHeightScale())

		for _, mc := range e.CircleColliders {
			mx, my := e.LocalPosToAbsPos(mc.LocalX, mc.LocalY)
			// ⭕️ 半径にスケールを適用
			mRadius := mc.Radius * eScale

			for _, tc := range t.CircleColliders {
				tx, ty := t.LocalPosToAbsPos(tc.LocalX, tc.LocalY)
				// ⭕️ 半径にスケールを適用
				tRadius := tc.Radius * tScale

				dx, dy := mx-tx, my-ty
				distSq := dx*dx + dy*dy
				// ⭕️ スケール済みの半径同士で比較する
				radSum := mRadius + tRadius
				if distSq <= radSum*radSum {
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
