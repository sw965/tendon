package tendon

import (
	"math"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Grid struct {
	*Element
	rows int
	cols int
	// Element.Childrenの要素と同じポインターを持つ事を想定
	// matrix[r][c] の順でアクセスする二次元スライス
	matrix [][]*Element
}

// NewGrid は指定された行列で空のグリッドを生成します
func NewGrid(rows, cols int, cellW, cellH, gap float64) *Grid {
	container := NewElement()

	// グリッド全体のサイズを計算
	totalW := float64(cols)*cellW + float64(cols-1)*gap
	totalH := float64(rows)*cellH + float64(rows-1)*gap

	// 透明な背景画像をセットすることで、Elementとしてのサイズを持たせる
	if totalW > 0 && totalH > 0 {
		// 欠けを防ぐために、切り上げ
		container.Image = ebiten.NewImage(int(math.Ceil(totalW)), int(math.Ceil(totalH)))
	}

	// 二次元スライスの初期化 (行ベース)
	matrix := make([][]*Element, rows)
	for r := 0; r < rows; r++ {
		matrix[r] = make([]*Element, cols)
		for c := 0; c < cols; c++ {
			cell := NewElement()
			// 座標を計算してセット
			// 列(c)がX座標、行(r)がY座標に対応
			cell.XRelativeToParent = float64(c) * (cellW + gap)
			cell.YRelativeToParent = float64(r) * (cellH + gap)

			container.AppendChild(cell) // 親要素（コンテナ）に追加
			matrix[r][c] = cell
		}
	}

	return &Grid{
		Element: container,
		rows:    rows,
		cols:    cols,
		matrix:  matrix,
	}
}

// NewBorderGrid は、各セルにシンプルな枠線が描画されたグリッドを生成します。
func NewBorderGrid(rows, cols int, cellW, cellH, gap float64, borderColor color.Color, borderWidth float64) *Grid {
	// ベースとなるグリッドを生成
	g := NewGrid(rows, cols, cellW, cellH, gap)

	// 欠けを防ぐために切り上げ
	wi := int(math.Ceil(cellW))
	hi := int(math.Ceil(cellH))
	wf := float64(wi)
	hf := float64(hi)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cell := g.GetCell(r, c)

			// セルのサイズに合わせてImageを作成
			img := ebiten.NewImage(wi, hi)

			// 四辺に長方形を描画して枠線とする
			ebitenutil.DrawRect(img, 0, 0, wf, borderWidth, borderColor)
			ebitenutil.DrawRect(img, 0, hf-borderWidth, wf, borderWidth, borderColor)
			ebitenutil.DrawRect(img, 0, 0, borderWidth, hf, borderColor)
			ebitenutil.DrawRect(img, wf-borderWidth, 0, borderWidth, hf, borderColor)

			cell.Image = img
		}
	}
	return g
}

// GetCell は指定した行(r)と列(c)のセルを返します
func (g *Grid) GetCell(r, c int) *Element {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return nil
	}
	return g.matrix[r][c]
}

// SetCell は特定のセルの Element 自体を入れ替えます
func (g *Grid) SetCell(r, c int, newElem *Element) {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return
	}

	oldElem := g.matrix[r][c]

	// 1. 座標の継承（入れ替えるなら位置を合わせるのが親切）
	newElem.XRelativeToParent = oldElem.XRelativeToParent
	newElem.YRelativeToParent = oldElem.YRelativeToParent

	// 2. 内部マトリックスの更新 (行ベース)
	g.matrix[r][c] = newElem

	// 3. 親子関係の更新
	for i, child := range g.Children {
		if child == oldElem {
			g.Children[i] = newElem
			newElem.Parent = g.Element
			break
		}
	}
}
