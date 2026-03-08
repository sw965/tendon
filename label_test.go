package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type LabelHighResGame struct {
	elements tendon.Elements
	counter  int
	dynamic  *tendon.Label
}

func (g *LabelHighResGame) Update() error {
	g.counter++
	// 約2秒ごとにサイズを 16px -> 32px -> 48px と切り替えるテスト
	if g.counter%2 == 0 {
		newSize := float64(16 + ((g.counter/120)%3)*16)
		g.dynamic.SetSize(newSize) // サイズを直接変更
		g.dynamic.SetText(fmt.Sprintf("SetSize(%.0fpx)", newSize))
	}
	return nil
}

func (g *LabelHighResGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.elements.Draw(screen)
}

func (g *LabelHighResGame) Layout(w, h int) (int, int) {
	return 640, 480
}

func TestLabelHighResolution(t *testing.T) {
	// 1. ネイティブ解像度で大きく描画（クッキリ）
	highRes := tendon.NewLabel("Native 32px (Crystal Clear)", 32)
	highRes.XRelativeToParent = 50
	highRes.YRelativeToParent = 80

	// 2. 16px で描いたものを 2倍に拡大（ボヤける）
	lowRes := tendon.NewLabel("Scaled 16px (Blurry)", 16)
	lowRes.XRelativeToParent = 50
	lowRes.YRelativeToParent = 180
	lowRes.SetScale(2.0) // 座標系で拡大

	// 3. 動的にサイズが変わるラベル
	dynamic := tendon.NewLabel("SetSize(16px)", 16)
	dynamic.XRelativeToParent = 50
	dynamic.YRelativeToParent = 300

	game := &LabelHighResGame{
		elements: tendon.Elements{
			highRes.Element,
			lowRes.Element,
			dynamic.Element,
		},
		dynamic: dynamic,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon High-Res Label Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}