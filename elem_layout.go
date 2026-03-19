package tendon

// --- 基本アライメント計算 ---

func (e *Element) calcHorizontalAlign(targetAbsX, targetWidth float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return targetAbsX
	case AlignCenter:
		return targetAbsX + (targetWidth-e.AbsWidth()) / 2
	case AlignEnd:
		return targetAbsX + targetWidth - e.AbsWidth()
	default:
		return targetAbsX
	}
}

func (e *Element) calcVerticalAlign(targetAbsY, targetHeight float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return targetAbsY
	case AlignCenter:
		return targetAbsY + (targetHeight-e.AbsHeight())/2
	case AlignEnd:
		return targetAbsY + targetHeight - e.AbsHeight()
	default:
		return targetAbsY
	}
}

// --- 標準配置 (回転無視) ---

func (e *Element) PlaceLeftOf(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	ax := tx - e.AbsWidth() - margin
	ay := e.calcVerticalAlign(ty, t.AbsHeight(), align)
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceRightOf(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	ax := tx + t.AbsWidth() + margin
	ay := e.calcVerticalAlign(ty, t.AbsHeight(), align)
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceAbove(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	ax := e.calcHorizontalAlign(tx, t.AbsWidth(), align)
	ay := ty - e.AbsHeight() - margin
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceBelow(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	ax := e.calcHorizontalAlign(tx, t.AbsWidth(), align)
	ay := ty + t.AbsHeight() + margin
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceCenterOf(target Component) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	ax := tx + (t.AbsWidth()-e.AbsWidth())/2
	ay := ty + (t.AbsHeight()-e.AbsHeight())/2
	e.SetAbsPos(ax, ay)
}

// --- Bounds配置 (回転後の外枠 AABB を考慮) ---

func (e *Element) PlaceLeftOfBounds(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	minX, minY, _, maxY := t.BoundingBox()
	ax := minX - e.AbsWidth() - margin
	ay := e.calcVerticalAlign(minY, maxY-minY, align)
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceRightOfBounds(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	_, minY, maxX, maxY := t.BoundingBox()
	ax := maxX + margin
	ay := e.calcVerticalAlign(minY, maxY-minY, align)
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceAboveBounds(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	minX, minY, maxX, _ := t.BoundingBox()
	ax := e.calcHorizontalAlign(minX, maxX-minX, align)
	ay := minY - e.AbsHeight() - margin
	e.SetAbsPos(ax, ay)
}

func (e *Element) PlaceBelowBounds(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	minX, _, maxX, maxY := t.BoundingBox()
	ax := e.calcHorizontalAlign(minX, maxX-minX, align)
	ay := maxY + margin
	e.SetAbsPos(ax, ay)
}

// --- Rotated配置 (回転に同期して辺に吸着) ---

func (e *Element) PlaceLeftOfRotated(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	e.Rotation = t.Rotation
	tw, th := t.BaseWidth(), t.BaseHeight()
	ew_abs := e.AbsWidth()
	twsc := t.AbsWidthScale()

	// ターゲットのピボットから見た、配置対象(e)のピボットの相対位置を計算
	// ew_abs と margin(画面ピクセル) はターゲットのスケールで割り戻してローカル単位にする
	lxPivot := -(ew_abs*(1-e.AnchorX)+margin)/twsc - (tw * t.AnchorX)
	lyPivot := e.calcRotatedAlign(th, e.BaseHeight()*(e.AbsHeightScale()/t.AbsHeightScale()), t.AnchorY, e.AnchorY, align)

	wx, wy := t.LocalPosToAbsPos(lxPivot, lyPivot)
	e.SetAbsPos(wx-e.AbsWidth()*e.AnchorX, wy-e.AbsHeight()*e.AnchorY)
}

func (e *Element) PlaceRightOfRotated(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	e.Rotation = t.Rotation
	tw, th := t.BaseWidth(), t.BaseHeight()
	twsc := t.AbsWidthScale()

	lxPivot := (tw * (1 - t.AnchorX)) + (margin / twsc) + (e.AbsWidth() * e.AnchorX / twsc)
	lyPivot := e.calcRotatedAlign(th, e.BaseHeight()*(e.AbsHeightScale()/t.AbsHeightScale()), t.AnchorY, e.AnchorY, align)

	wx, wy := t.LocalPosToAbsPos(lxPivot, lyPivot)
	e.SetAbsPos(wx-e.AbsWidth()*e.AnchorX, wy-e.AbsHeight()*e.AnchorY)
}

func (e *Element) PlaceAboveRotated(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	e.Rotation = t.Rotation
	tw, th := t.BaseWidth(), t.BaseHeight()
	thsc := t.AbsHeightScale()

	lxPivot := e.calcRotatedAlign(tw, e.BaseWidth()*(e.AbsWidthScale()/t.AbsWidthScale()), t.AnchorX, e.AnchorX, align)
	lyPivot := -(th * t.AnchorY) - (margin / thsc) - (e.AbsHeight() * (1 - e.AnchorY) / thsc)

	wx, wy := t.LocalPosToAbsPos(lxPivot, lyPivot)
	e.SetAbsPos(wx-e.AbsWidth()*e.AnchorX, wy-e.AbsHeight()*e.AnchorY)
}

func (e *Element) PlaceBelowRotated(target Component, margin float64, align Alignment) {
	t := target.BaseElement()
	e.Rotation = t.Rotation
	tw, th := t.BaseWidth(), t.BaseHeight()
	thsc := t.AbsHeightScale()

	lxPivot := e.calcRotatedAlign(tw, e.BaseWidth()*(e.AbsWidthScale()/t.AbsWidthScale()), t.AnchorX, e.AnchorX, align)
	lyPivot := (th * (1 - t.AnchorY)) + (margin / thsc) + (e.AbsHeight() * e.AnchorY / thsc)

	wx, wy := t.LocalPosToAbsPos(lxPivot, lyPivot)
	e.SetAbsPos(wx-e.AbsWidth()*e.AnchorX, wy-e.AbsHeight()*e.AnchorY)
}

func (e *Element) calcRotatedAlign(targetSize, elemSize, targetAnchor, elemAnchor float64, align Alignment) float64 {
	switch align {
	case AlignStart:
		return -(targetSize * targetAnchor) + (elemSize * elemAnchor)
	case AlignCenter:
		return (targetSize-elemSize)/2 - (targetSize * targetAnchor) + (elemSize * elemAnchor)
	case AlignEnd:
		return targetSize - (targetSize * targetAnchor) - elemSize + (elemSize * elemAnchor)
	default:
		return 0
	}
}

// --- ゲッター系メソッド (既存の互換性を維持) ---

func (e *Element) AbsPosLeftOf(target Component, margin float64, align Alignment) (float64, float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	return tx - e.AbsWidth() - margin, e.calcVerticalAlign(ty, t.AbsHeight(), align)
}

func (e *Element) AbsPosRightOf(target Component, margin float64, align Alignment) (float64, float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	return tx + t.AbsWidth() + margin, e.calcVerticalAlign(ty, t.AbsHeight(), align)
}

func (e *Element) AbsPosAbove(target Component, margin float64, align Alignment) (float64, float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	return e.calcHorizontalAlign(tx, t.AbsWidth(), align), ty - e.AbsHeight() - margin
}

func (e *Element) AbsPosBelow(target Component, margin float64, align Alignment) (float64, float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	return e.calcHorizontalAlign(tx, t.AbsWidth(), align), ty + t.AbsHeight() + margin
}

func (e *Element) AbsPosCenterOf(target Component) (float64, float64) {
	t := target.BaseElement()
	tx, ty := t.AbsPos()
	return tx + (t.AbsWidth()-e.AbsWidth())/2, ty + (t.AbsHeight()-e.AbsHeight())/2
}

// --- レイアウトヘルパー ---

func (e *Element) AbsXInLayoutWidth(w float64, align Alignment) float64 {
	return e.calcHorizontalAlign(0, w, align)
}

func (e *Element) AbsYInLayoutHeight(h float64, align Alignment) float64 {
	return e.calcVerticalAlign(0, h, align)
}
