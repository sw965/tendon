package tendon_test

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sw965/tendon"
)

const (
	MonsterCardPath = "C:/Users/rayze/Desktop/img/Dark Magician Girl.jpg"
)

type ScaleTestGame struct {
	elements tendon.Elements
	hitBuf   tendon.Elements
	dragBuf  tendon.Elements
}

func (g *ScaleTestGame) Update() error {
	// --- スケールのテスト操作 ---
	// Ctrl + ホイールで「全ての要素の親」である grid のスケールを操作
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight) {
		_, dy := ebiten.Wheel()
		if dy != 0 && len(g.elements) > 0 {
			root := g.elements[0]
			newScale := root.WidthScale + dy*0.1
			if newScale < 0.1 {
				newScale = 0.1
			}
			root.SetScale(newScale)
			fmt.Printf("Parent Scale: %.2f\n", newScale)
		}
	}

	// 各要素の状態更新
	for _, e := range g.elements {
		e.Update()
	}

	// 当たり判定
	g.hitBuf = g.hitBuf[:0]
	cx, cy := ebiten.CursorPosition()
	g.elements.FindAllHitTest(float64(cx), float64(cy), &g.hitBuf)

	// ホバー状態の更新
	g.elements.UpdateHover(g.hitBuf)

	// ドラッグ移動の更新
	g.dragBuf.UpdateDragMove()

	// 左クリック：ドラッグ開始
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && len(g.hitBuf) > 0 {
		target := g.hitBuf[0]
		if target.Draggable {
			target.StartDrag()
			target.Z = 999
			g.dragBuf = append(g.dragBuf, target)
		}
	}

	// 左クリック離した：ドラッグ終了
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.elements.StopAllDrag()
		g.dragBuf = g.dragBuf[:0]
	}

	return nil
}

func (g *ScaleTestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.elements.Draw(screen)

	// 操作説明
	ebitenutil.DebugPrintAt(screen, "Ctrl + Mouse Wheel: Scale Parent (Grid)", 10, 10)
	ebitenutil.DebugPrintAt(screen, "Drag Card: Check position consistency", 10, 30)
}

func (g *ScaleTestGame) Layout(w, h int) (int, int) {
	return 1280, 720
}

func TestScaleConsistency(t *testing.T) {
	// 1. 親となるグリッドの作成
	grid := tendon.NewBorderGrid(2, 5, 120, 170, 15, color.RGBA{0, 255, 0, 255}, 2)
	grid.XRelativeToParent = 100
	grid.YRelativeToParent = 100
	grid.Image.Fill(color.RGBA{255, 0, 0, 40}) // 親の領域を赤く可視化

	img, _, err := ebitenutil.NewImageFromFile(MonsterCardPath)
	if err != nil {
		t.Fatalf("画像の読み込み失敗: %v", err)
	}

	// 2. 子要素（カード）の作成
	card := tendon.NewElement()
	card.Id = 1
	card.Image = img
	card.SetScale(0.1)
	card.Draggable = true

	// 3. 配置
	if cell := grid.GetCell(0, 0); cell != nil {
		cell.AppendChild(card)
		card.PlaceCenterOf(cell)
	}

	game := &ScaleTestGame{
		elements: tendon.Elements{grid.Element},
		hitBuf:   make(tendon.Elements, 0, 10),
		dragBuf:  make(tendon.Elements, 0, 5),
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Tendon Scale Consistency Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}