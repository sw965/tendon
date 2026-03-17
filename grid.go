package tendon

import (
	"image/color"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO gridの設計を考えておく。枠と中身を切り分ける？
type Grid struct {
	*Element
	Rows           int
	Cols           int
	CellW          float64
	CellH          float64
	Gap            float64
	// TODO 二次元にするか検討する
	LayoutChildren Components
}

func NewGrid(cols, rows int, cellW, cellH, gap float64) *Grid {
	base := NewElement()

	// グリッド全体のサイズ
	totalW := float64(cols)*cellW + float64(cols-1)*gap
	totalH := float64(rows)*cellH + float64(rows-1)*gap

	if totalW > 0 && totalH > 0 {
		// baseに透明な画像をセットする
		base.Image = ebiten.NewImage(int(math.Ceil(totalW)), int(math.Ceil(totalH)))
	}

	g := &Grid{
		Element:        base,
		Cols:           cols,
		Rows:           rows,
		CellW:          cellW,
		CellH:          cellH,
		Gap:            gap,
		LayoutChildren: Components{},
	}

	for i := 0; i < cols*rows; i++ {
		cell := NewElement()
		if cellW > 0 && cellH > 0 {
			cell.Image = ebiten.NewImage(int(math.Ceil(cellW)), int(math.Ceil(cellH)))
		}
		g.AppendChild(cell)
	}

	g.Reflow()
	return g
}

func NewBorderGrid(cols, rows int, cellW, cellH, gap float64, borderColor color.Color, borderWidth float64) *Grid {
	g := NewGrid(cols, rows, cellW, cellH, gap)

	wi := int(math.Ceil(cellW))
	hi := int(math.Ceil(cellH))
	wf := float64(wi)
	hf := float64(hi)

	// 最初に追加された空セルに枠線を描画する
	for i, c := range g.LayoutChildren {
		if i >= cols*rows {
			break
		}
		cell := c.BaseElement()
		img := ebiten.NewImage(wi, hi)
		ebitenutil.DrawRect(img, 0, 0, wf, borderWidth, borderColor)
		ebitenutil.DrawRect(img, 0, hf-borderWidth, wf, borderWidth, borderColor)
		ebitenutil.DrawRect(img, 0, 0, borderWidth, hf, borderColor)
		ebitenutil.DrawRect(img, wf-borderWidth, 0, borderWidth, hf, borderColor)
		cell.Image = img
	}
	return g
}

// AppendChild は Box と同様に本体と名簿の両方に追加します。
func (g *Grid) AppendChild(child Component) {
	g.Element.AppendChild(child)
	g.LayoutChildren = append(g.LayoutChildren, child)
}

// RemoveChild は Box と同様に両方から削除します。
func (g *Grid) RemoveChild(target Component) bool {
	removed := g.Element.RemoveChild(target)
	if !removed {
		return false
	}

	t := target.BaseElement()
	index := slices.IndexFunc(g.LayoutChildren, func(c Component) bool {
		return c.BaseElement() == t
	})

	if index != -1 {
		g.LayoutChildren = slices.Delete(g.LayoutChildren, index, index+1)
	}
	return true
}

func (g *Grid) GetCell(c, r int) Component {
	if c < 0 || c >= g.Cols || r < 0 || r >= g.Rows {
		return nil
	}
	index := r*g.Cols + c
	if index < len(g.LayoutChildren) {
		return g.LayoutChildren[index]
	}
	return nil
}

// SetCell は特定のセルのコンポーネントを安全に入れ替えます。
// ★ 引数を Component に変更し、古い要素のメモリリークを塞ぎました。
func (g *Grid) SetCell(c, r int, newElem Component) {
	if c < 0 || c >= g.Cols || r < 0 || r >= g.Rows {
		return
	}
	index := r*g.Cols + c
	if index >= len(g.LayoutChildren) {
		return
	}

	oldElem := g.LayoutChildren[index]
	oldBase := oldElem.BaseElement()

	// ★ 重要：古い要素の親参照を切り、ガベージコレクション(メモリ解放)の対象にする
	oldBase.Parent = nil

	// 新しい要素の親を設定
	newBase := newElem.BaseElement()
	newBase.Parent = g.Element

	// 名簿の入れ替え
	g.LayoutChildren[index] = newElem

	// 本体の描画ツリー(Children)の入れ替え
	for i, child := range g.Children {
		if child.BaseElement() == oldBase {
			g.Children[i] = newElem
			g.childrenOrderDirty = true
			break
		}
	}

	// TODO 必要か検討する
	g.Reflow()
}

// Reflow は LayoutChildren の並び順に従って、自動的に折り返してグリッド状に整列させます。
func (g *Grid) Reflow() {
	for i, c := range g.LayoutChildren {
		child := c.BaseElement()

		// インデックスから 列(col) と 行(row) を計算
		col := i % g.Cols
		row := i / g.Cols

		child.XRelativeToParent = float64(col) * (g.CellW + g.Gap)
		child.YRelativeToParent = float64(row) * (g.CellH + g.Gap)
	}
}
