package tendon

// AbsPos は要素の「左上」の画面上の絶対座標を返します。
func (e *Element) AbsPos() (float64, float64) {
	m := e.AbsGeoM()
	return m.Apply(0, 0)
}

func (e *Element) AbsPosToRelPosToParent(absX, absY float64) (relX, relY float64) {
	if e.Parent == nil {
		return absX, absY
	}
	pm := e.Parent.AbsGeoM()
	pm.Invert()
	return pm.Apply(absX, absY)
}

func (e *Element) SetAbsPos(absX, absY float64) {
	e.XRelativeToParent, e.YRelativeToParent = e.AbsPosToRelPosToParent(absX, absY)
}

func (e *Element) LocalPosToAbsPos(lx, ly float64) (float64, float64) {
	m := e.AbsGeoM()
	w, h := e.BaseWidth(), e.BaseHeight()
	xRelToImg := lx + w*e.AnchorX
	yRelToImg := ly + h*e.AnchorY
	// 画像サイズを w = 100, h = 50, e.AnchorX = e.AnchorY = 0.5 としたとき
	// lx = ly = 0 のとき、ローカル座標の定義より、これは画像の中心を示す
	// よって画像の中心の絶対座標を知りたい (ローカル座標から絶対座標に変換したい)
	// m.Applyは、引数に画像から見たときの座標を渡せば、絶対座標に変換してくれる
	// 例えば、 m.Apply(0, 0)であれば、画像の左上の絶対座標を戻り値として返す
	// 上記の例では、画像の中心の絶対座標を知りたいため、m.Apply(50, 25)を渡せばいい。
	// これを計算式にあてはめると、
	// xRelToImg = 0 + 100 * 0.5 = 50
	// yRelToImg = 0 + 50 * 0.5 = 25
	// これにより、ローカル座標から絶対座標へ変換出来る
	return m.Apply(xRelToImg, yRelToImg)
}

func (e *Element) PointToLocalPos(pointX, pointY float64) (float64, float64) {
	m := e.AbsGeoM()
	m.Invert()
	// (pointX, pointY) を 自身(e) の スケールや回転を元に戻した状態の画像の左上を 原点(0, 0) としたときの相対座標へ変換する
	xRelToImg, yRelToImg := m.Apply(pointX, pointY)
	w, h := e.BaseWidth(), e.BaseHeight()
	// ローカル座標は「スケールや回転を元に戻した状態の画像に対するアンカーを原点とみなした座標」が定義
	// 例えば、画像サイズを w = 200, h = 200, e.AnchorX = e.AnchorY = 0.5 (画像の中心) としたとき
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

	// 画像の左上の絶対座標を取得
	topLeftX, topLeftY := m.Apply(0, 0)
	// 画像の右上の絶対座標を取得
	topRightX, topRightY := m.Apply(w, 0)
	// 画像の左下の絶対座標を取得
	bottomLeftX, bottomLeftY := m.Apply(0, h)
	// 画像の右下の絶対座標を取得
	bottomRightX, bottomRightY := m.Apply(w, h)

    minX := min(topLeftX, topRightX, bottomLeftX, bottomRightX)
	maxX := max(topLeftX, topRightX, bottomLeftX, bottomRightX)
    minY := min(topLeftY, topRightY, bottomLeftY, bottomRightY)
    maxY := max(topLeftY, topRightY, bottomLeftY, bottomRightY)
    return minX, minY, maxX, maxY
}
