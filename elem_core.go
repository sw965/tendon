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

	Image       *ebiten.Image
	WidthScale  float64
	HeightScale float64
	Filter      ebiten.Filter

	Visible     bool
	Enabled     bool
	PassThrough bool
	Z           int

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

	// カプセル化の検討
	Children        Components
	Parent          *Element
	childOrderDirty bool

	Rotation float64
	// パーセンテージ
	AnchorX         float64
	AnchorY         float64
	Shape           Shape
	CircleColliders []CircleCollider
}

func NewElement() *Element {
	e := &Element{
		Visible:         true,
		Enabled:         true,
		EasingFunc:      func(current, target float64) float64 { return target },
		AnchorX:         0.5,
		AnchorY:         0.5,
		childOrderDirty: true,
	}
	e.SetScale(1.0)
	return e
}

func (e *Element) Update() {
	if !e.Visible {
		return
	}

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

	for _, child := range e.Children {
		child.Update()
	}
}

func (e *Element) Draw(screen *ebiten.Image) {
	if !e.Visible {
		return
	}

	if e.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM = e.TransformMatrix()
		op.Filter = e.Filter
		screen.DrawImage(e.Image, op)
	}

	if e.childOrderDirty {
		e.Children.SortByZAsc()
		e.childOrderDirty = false
	}

	for _, child := range e.Children {
		child.Draw(screen)
	}
}

func (e *Element) BaseElement() *Element {
	return e
}

func (e *Element) AppendChild(child Component) {
	child.BaseElement().Parent = e
	e.Children = append(e.Children, child)
	e.childOrderDirty = true
}

func (e *Element) RemoveChild(target Component) bool {
	index := slices.Index(e.Children, target)
	if index == -1 {
		return false
	}

	// 子要素の親参照を解除
	target.BaseElement().Parent = nil

	// スライスから削除
	e.Children = slices.Delete(e.Children, index, index+1)
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

// LocalMatrix は「親要素から見た」この要素自身の相対的な変換行列を作成します。
func (e *Element) LocalMatrix() ebiten.GeoM {
	var m ebiten.GeoM
	w, h := e.BaseWidth(), e.BaseHeight()
	// 画像のアンカー位置を原点(0,0)に持ってくる
	m.Translate(-w*e.AnchorX, -h*e.AnchorY)
	// ローカルのスケールと回転を適用
	m.Scale(e.WidthScale, e.HeightScale)
	m.Rotate(e.Rotation)
	// 親要素の空間内での配置位置へ移動
	m.Translate(e.XRelativeToParent+w*e.AnchorX*e.WidthScale, e.YRelativeToParent+h*e.AnchorY*e.HeightScale)
	return m
}

// TransformMatrix は、この要素の画面上での最終的な変換行列（ワールド行列）を返します。
func (e *Element) TransformMatrix() ebiten.GeoM {
	m := e.LocalMatrix()
	// 親要素が存在する場合、親の変換行列を「掛け合わせる(Concat)」ことで、
	// 親の移動・回転・スケールに完全に追従する（シーングラフ）ようになります。
	if e.Parent != nil {
		m.Concat(e.Parent.TransformMatrix())
	}
	return m
}