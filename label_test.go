package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type LabelTestGame struct {
	elements tendon.Elements
	counter  int
	dynamic  *tendon.Label
}

func (g *LabelTestGame) Update() error {
	g.counter++

	// 約1秒ごとにテキストとサイズを更新するテスト
	if g.counter%60 == 0 {
		// 現在のフォントソースを維持しつつ、サイズだけ計算
		newSize := float64(16 + ((g.counter/60)%3)*16)
		newText := fmt.Sprintf("Size: %.0fpx / Frame: %d", newSize, g.counter)

		// 【ポイント】現在の Font().Source を取得して Update に渡す
		// これにより、差分計算が走り、必要な場合のみ再描画されます
		g.dynamic.Update(newText, g.dynamic.Font().Source, newSize)
	}
	return nil
}

func (g *LabelTestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.elements.Draw(screen)
}

func (g *LabelTestGame) Layout(w, h int) (int, int) {
	return 640, 480
}

func TestLabelRefactored(t *testing.T) {
	return
	// 1. 静的なラベルの作成（エラーハンドリング付き）
	staticLabel, err := tendon.NewLabel("Static Label (Refactored)", 24)
	if err != nil {
		t.Fatal(err)
	}
	staticLabel.XRelativeToParent = 50
	staticLabel.YRelativeToParent = 50

	// 2. 動的なラベルの作成
	dynamicLabel, err := tendon.NewLabel("Dynamic Label", 16)
	if err != nil {
		t.Fatal(err)
	}
	dynamicLabel.XRelativeToParent = 50
	dynamicLabel.YRelativeToParent = 150

	// 3. スケール変更のテスト（ボヤけの確認）
	scaledLabel, err := tendon.NewLabel("Scaled 16px (Should be Blurry)", 16)
	if err != nil {
		t.Fatal(err)
	}
	scaledLabel.XRelativeToParent = 50
	scaledLabel.YRelativeToParent = 250
	scaledLabel.SetScale(2.5) // Element の機能で拡大

	game := &LabelTestGame{
		elements: tendon.Elements{
			staticLabel.Element,
			dynamicLabel.Element,
			scaledLabel.Element,
		},
		dynamic: dynamicLabel,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon Refactored Label Test")

	// 実際に動かして、Update による再描画と Dirty Check の挙動を確認
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
