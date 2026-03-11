package tendon_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type ButtonTestGame struct {
	elements tendon.Components // ★ Elements から変更
}

func (g *ButtonTestGame) Update() error {
	return nil
}

func (g *ButtonTestGame) Draw(screen *ebiten.Image) {
	// 各ボタンの色が分かりやすいように、背景は暗めのグレーにする
	screen.Fill(color.RGBA{40, 40, 45, 255})
	g.elements.Draw(screen)
}

func (g *ButtonTestGame) Layout(w, h int) (int, int) {
	return 640, 480
}

func TestButtonAutoFit(t *testing.T) {
	return
	// ボタンを縦に綺麗に並べるための Box を作成
	box := tendon.NewBox(640, 480, 20)
	box.MainAlignment = tendon.AlignCenter
	box.CrossAlignment = tendon.AlignCenter

	// 1. 短いテキスト（縮小されない）
	btn1, err := tendon.NewButton(200, 60, "OK", tendon.Blue)
	if err != nil {
		t.Fatal(err)
	}

	// 2. 中くらいのテキスト
	btn2, err := tendon.NewButton(200, 60, "Start Game", tendon.Orange)
	if err != nil {
		t.Fatal(err)
	}

	// 3. 以前はみ出していた長さのテキスト
	btn3, err := tendon.NewButton(200, 60, "YellowGreen", tendon.YellowGreen)
	if err != nil {
		t.Fatal(err)
	}

	// 4. 極端に長いテキスト（かなり小さく縮小されて枠に収まるはず）
	btn4, err := tendon.NewButton(200, 60, "This is a very very long text!", tendon.Red)
	if err != nil {
		t.Fatal(err)
	}

	// ★ .Element を付けずにそのまま追加できる！
	box.AppendChild(btn1)
	box.AppendChild(btn2)
	box.AppendChild(btn3)
	box.AppendChild(btn4)
	
	// ★ 方向をプロパティで指定してから引数なしでUpdate
	box.Orientation = tendon.Vertical
	box.Update()

	game := &ButtonTestGame{
		elements: tendon.Components{box}, // ★ .Element を削除
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon Button Auto-Fit Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}