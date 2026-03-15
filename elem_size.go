package tendon

func (e *Element) BaseWidth() float64 {
	if e.Image == nil {
		return 0
	}
	return float64(e.Image.Bounds().Dx())
}

func (e *Element) BaseHeight() float64 {
	if e.Image == nil {
		return 0
	}
	return float64(e.Image.Bounds().Dy())
}

func (e *Element) AbsWidth() float64 {
	scale := e.AbsWidthScale()
	return e.BaseWidth() * scale
}

func (e *Element) AbsHeight() float64 {
	scale := e.AbsHeightScale()
	return e.BaseHeight() * scale
}

func (e *Element) AbsWidthScale() float64 {
	if e.Parent == nil {
		return e.WidthScale
	}
	return e.Parent.AbsWidthScale() * e.WidthScale
}

func (e *Element) AbsHeightScale() float64 {
	if e.Parent == nil {
		return e.HeightScale
	}
	return e.Parent.AbsHeightScale() * e.HeightScale
}

func (e *Element) SetScale(s float64) {
	e.WidthScale = s
	e.HeightScale = s
}
