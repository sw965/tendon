package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
)

// Element はUIツリーを構成する基本コンポーネントです。
// 相対座標の計算、親子関係の構築、描画、およびマウス操作の判定を担います。
//
// 【3つの主要フラグの違いと、子要素への影響】
//
// 以下のフラグを組み合わせることで、様々なUIの要件（コンテナ、装飾、無効化など）を実現します。
//
//		フラグ名    | 意味           | 描画 | 当たり判定 | 状態更新 | 子要素への影響
//		------------|----------------|------|------------|----------|----------------------------------------
//		Visible     | 可視性         |  ❌  |     ❌     |    ❌    | falseの場合、子要素も全て非表示・判定外になる
//		Enabled     | 有効/無効      |  ⭕️  |     ⭕️     |    ❌    | falseにしても、子要素のEnabledには影響しない
//		PassThrough | 当たり判定透過 |  ⭕️  |     ❌     |    ❌    | trueにしても、子要素の当たり判定は通常通り行われる
//
//	  - Visible:
//	    要素を完全にツリーから除外したい場合に使用します（例：非表示になったウィンドウ）。
//
//	  - Enabled:
//	    クリック等の対象として「重なり」は検知させたいが、操作は弾きたい場合に使用します。
//	    （例：グレーアウトして押せなくなったボタン）。
//
//	  - PassThrough:
//	    見た目だけで当たり判定を下に貫通させたい場合や、当たり判定計算の負荷を下げたい場合に使用します。
//	    （例：ボタンの上に乗っているテキスト、複数のボタンを並べるための透明なレイアウト用親コンテナ）。
type Element struct {
	Id                int
	XRelativeToParent float64
	YRelativeToParent float64

	Image        *ebiten.Image
	widthScale   float64
	heightScale  float64
	isScaleDirty bool
	Filter       ebiten.Filter

	Visible     bool
	Enabled     bool
	PassThrough bool
	z           int

	// ドラッグ関連
	Draggable   bool
	ManualDrag  bool
	isDragging  bool
	dragOffsetX float64
	dragOffsetY float64
	DragDeltaX  float64
	DragDeltaY  float64

	isHovered      bool
	isJustHoverIn  bool
	isJustHoverOut bool

	// 移動アニメーション関連
	toX                float64
	toY                float64
	isMoving           bool
	isJustMoveFinished bool
	EasingFunc         EasingFunc

	// TODO カプセル化の検討
	Children           Components
	Parent             *Element
	childrenOrderDirty bool

	Rotation float64
	// パーセンテージ
	AnchorX         float64
	AnchorY         float64
	Shape           Shape
	CircleColliders []CircleCollider
	rebuildCollider func()
}

func NewElement() *Element {
	e := &Element{
		Visible:            true,
		Enabled:            true,
		EasingFunc:         func(current, target float64) float64 { return target },
		AnchorX:            0.5,
		AnchorY:            0.5,
		childrenOrderDirty: true,
		isScaleDirty:       false,
	}
	e.SetScale(1.0)
	return e
}

func (e *Element) Update() {
	if !e.Visible {
		return
	}

	e.sortChildren()
	// TODO ここでコピーする？

	e.isJustMoveFinished = false
	if e.isMoving {
		if e.EasingFunc == nil {
			e.XRelativeToParent = e.toX
			e.YRelativeToParent = e.toY
		} else {
			e.XRelativeToParent = e.EasingFunc(e.XRelativeToParent, e.toX)
			e.YRelativeToParent = e.EasingFunc(e.YRelativeToParent, e.toY)
		}

		if e.XRelativeToParent == e.toX && e.YRelativeToParent == e.toY {
			e.isMoving = false
			e.isJustMoveFinished = true
		}
	}

	// イテレーションを壊さないためのコピー
	// 例えば、child.Updateでe.Childrenを削除したとき、その後のループが壊れる
	// コピーすれば、e.Children = nil としてしまっても、ccはnilにならない
	// ただし、浅いコピーであるから、childの状態を書き換えることは出来る
	cc := make(Components, len(e.Children))
	copy(cc, e.Children)
	for _, child := range cc {
		child.Update()
	}
}

func (e *Element) Draw(screen *ebiten.Image) {
	if !e.Visible {
		return
	}

	if e.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM = e.AbsGeoM()
		op.Filter = e.Filter
		screen.DrawImage(e.Image, op)
	}

	// ここのコードの意図は、Element.Updateを参照
	e.sortChildren()
	cc := make(Components, len(e.Children))
	copy(cc, e.Children)
	for _, child := range cc {
		child.Draw(screen)
	}
}

func (e *Element) BaseElement() *Element {
	return e
}

func (e *Element) Z() int {
	return e.z
}

func (e *Element) SetZ(z int) {
	if e.z == z {
		return
	}
	e.z = z
	// 自身のZが変わったら、親にソートし直しを要求する
	if e.Parent != nil {
		e.Parent.childrenOrderDirty = true
	}
}

func (e *Element) sortChildren() {
	if e.childrenOrderDirty {
		slices.SortStableFunc(e.Children, func(a, b Component) int {
			return a.BaseElement().z - b.BaseElement().z
		})
		e.childrenOrderDirty = false
	}
}

func (e *Element) AppendChild(child Component) {
	child.BaseElement().Parent = e
	e.Children = append(e.Children, child)
	e.childrenOrderDirty = true
}

func (e *Element) RemoveChild(target Component) bool {
	t := target.BaseElement()

	index := slices.IndexFunc(e.Children, func(c Component) bool {
		return c.BaseElement() == t
	})

	if index == -1 {
		return false
	}

	t.Parent = nil
	e.Children = slices.Delete(e.Children, index, index+1)
	e.childrenOrderDirty = true
	return true
}

func (e *Element) Dispose() {
	if e.Image != nil {
		e.Image.Dispose()
		e.Image = nil
	}
	for _, child := range e.Children {
		child.BaseElement().Dispose()
	}
}

func (e *Element) RelGeoM() ebiten.GeoM {
	var m ebiten.GeoM
	w, h := e.BaseWidth(), e.BaseHeight()
	// 画像のアンカー位置を原点(0,0)とみなす
	m.Translate(-w*e.AnchorX, -h*e.AnchorY)
	// 原点 を中心としたスケール(縮小・拡大)の適用
	m.Scale(e.widthScale, e.heightScale)
	// 原点 を中心とした回転の適用
	m.Rotate(e.Rotation)

	m.Translate(w*e.AnchorX*e.widthScale, h*e.AnchorY*e.heightScale)
	// 親から見た相対位置に移動する
	m.Translate(e.XRelativeToParent, e.YRelativeToParent)
	return m
}

func (e *Element) AbsGeoM() ebiten.GeoM {
	m := e.RelGeoM()
	if e.Parent != nil {
		m.Concat(e.Parent.AbsGeoM())
	}
	return m
}

func (e *Element) resolveDirtyScale() {
	if !e.isScaleDirty {
		return
	}
	if e.rebuildCollider != nil {
		e.rebuildCollider()
	}
	e.isScaleDirty = false
}