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
		return e.widthScale
	}
	return e.Parent.AbsWidthScale() * e.widthScale
}

func (e *Element) AbsHeightScale() float64 {
	if e.Parent == nil {
		return e.heightScale
	}
	return e.Parent.AbsHeightScale() * e.heightScale
}

func (e *Element) SetWidthScale(s float64) {
    if e.widthScale == s {
        return
    }
    e.widthScale = s
    e.markAllScaleDirty()
}

func (e *Element) SetHeightScale(s float64) {
    if e.heightScale == s {
        return
    }
    e.heightScale = s
    e.markAllScaleDirty()
}

func (e *Element) SetScale(s float64) {
    if e.widthScale == s && e.heightScale == s {
        return
    }
    e.widthScale = s
    e.heightScale = s
    e.markAllScaleDirty()
}

func (e *Element) markAllScaleDirty() {
    e.isScaleDirty = true
    for _, child := range e.Children {
        child.BaseElement().markAllScaleDirty()
    }
}