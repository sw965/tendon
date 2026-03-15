package tendon

import (
	"math"
)

func (e *Element) GetAngle() float64 {
	return e.Rotation * 180 / math.Pi
}

func (e *Element) SetAngle(deg float64) {
	e.Rotation = deg * math.Pi / 180
}

// AbsRotation は、親要素の回転を含めた画面上での絶対的な回転角（ラジアン）を返します。
func (e *Element) AbsRotation() float64 {
	if e.Parent == nil {
		return e.Rotation
	}
	// 親の絶対角度に、自身のローカル角度を足す
	return e.Parent.AbsRotation() + e.Rotation
}

// GetAbsAngle は画面上での絶対的な角度（度数法）を返します。
func (e *Element) GetAbsAngle() float64 {
	return e.AbsRotation() * 180 / math.Pi
}