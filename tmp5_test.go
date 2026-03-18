package tendon_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sw965/tendon"
)

type IsometricTestGame struct {
	root     *tendon.Element
	boardCnt *tendon.Element
	grid     *tendon.Grid
	pieces   []*tendon.Element
	angle    float64
}

func (g *IsometricTestGame) Update() error {
	// 1. 盤面全体をゆっくり回転させる
	g.angle += 0.2
	g.boardCnt.SetAngle(g.angle)

	// 盤面のレイアウト計算を最新にするために一度Update
	g.boardCnt.Update()

	// 2. 駒をチェス盤の特定のマス目に「吸着」させる
	positions := [][2]int{
		{0, 0}, {3, 3}, {7, 7}, {2, 5}, {6, 1},
	}

	for i, piece := range g.pieces {
		if i >= len(positions) {
			break
		}
		col, row := positions[i][0], positions[i][1]

		// マス目の要素を取得
		cell := g.grid.GetCell(col, row).BaseElement()

		// ★ ここが要！：マス目の「中心(ローカル 0,0)」の画面上の絶対座標を取得
		cx, cy := cell.LocalPosToAbsPos(0, 0)

		// 駒をその絶対座標に移動させる
		piece.SetAbsPos(cx, cy)
	}

	g.root.Update()
	return nil
}

func (g *IsometricTestGame) Draw(screen *ebiten.Image) {
	// 背景色
	screen.Fill(color.RGBA{30, 30, 35, 255})

	// Rootを描画すれば、盤面も駒も全て描画される
	g.root.Draw(screen)

	msg := "ISOMETRIC (2.5D) BOARD DEMO\n\n"
	msg += "- Board : Rotated & Height-Scaled (0.5)\n"
	msg += "- Pieces: Root child, tracking cell's AbsPos\n\n"
	msg += "Look how the pieces stand upright\nwhile the board rotates!"
	ebitenutil.DebugPrint(screen, msg)
}

func (g *IsometricTestGame) Layout(w, h int) (int, int) { return 800, 600 }

func TestIsometricBoard(t *testing.T) {
	root := tendon.NewElement()

	// --- 1. 盤面コンテナの作成 ---
	boardCnt := tendon.NewElement()
	boardCnt.XRelativeToParent = 400
	boardCnt.YRelativeToParent = 300
	boardCnt.heightScale = 0.5 // ★ ここで縦を半分に潰して「斜め視点」を作る！
	boardCnt.AnchorX, boardCnt.AnchorY = 0.5, 0.5
	root.AppendChild(boardCnt)

	// --- 2. チェス盤 (Grid) の作成 ---
	cols, rows := 8, 8
	cellSize := 40.0
	grid := tendon.NewGrid(cols, rows, cellSize, cellSize, 0)
	// Grid全体のアンカーを中心にする
	grid.AnchorX, grid.AnchorY = 0.5, 0.5
	boardCnt.AppendChild(grid)

	// 市松模様に色を塗る
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cell := grid.GetCell(c, r).BaseElement()
			if (c+r)%2 == 0 {
				cell.Image.Fill(color.RGBA{200, 200, 200, 255}) // 白マス
			} else {
				cell.Image.Fill(color.RGBA{80, 80, 80, 255}) // 黒マス
			}
		}
	}

	// --- 3. 駒の作成 ---
	var pieces []*tendon.Element
	for i := 0; i < 5; i++ {
		piece := tendon.NewElement()
		piece.Image = ebiten.NewImage(20, 40)

		if i == 0 {
			piece.Image.Fill(tendon.Red) // 1つだけ目立たせる
		} else {
			piece.Image.Fill(tendon.Blue)
		}

		// ★ ここがポイント：アンカーを「足元中央」にする
		piece.AnchorX = 0.5
		piece.AnchorY = 1.0

		// 駒は盤面の子ではなく、Rootの直接の子にする（盤面の変形を受けないようにするため）
		root.AppendChild(piece)
		pieces = append(pieces, piece)
	}

	game := &IsometricTestGame{
		root:     root,
		boardCnt: boardCnt,
		grid:     grid,
		pieces:   pieces,
		angle:    45, // 初期角度
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Tendon Isometric Board")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
