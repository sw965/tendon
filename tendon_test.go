package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type TestGame struct {
	elements tendon.Elements
}

func (g *TestGame) Update() error {
	g.elements.Update(0, 0)
	return nil
}

func (g *TestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	g.elements.Draw(screen)
}

func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func TestInteractive(t *testing.T) {
	// 1. 大きな「親パネル」
	panel := tendon.NewButton(50, 50, 300, 300, "Parent Panel", color.RGBA{80, 80, 150, 255})
	panel.Z = 1
	panel.Draggable = true

	// ★ 追加: 親パネルのホバー判定を可視化
	panel.OnMouseEnter = func(e *tendon.Element) {
		fmt.Println("🟦 親パネルにマウスが【入りました】")
	}
	panel.OnMouseLeave = func(e *tendon.Element) {
		fmt.Println("🟦 親パネルからマウスが【出ました】")
	}

	// 2. パネルの中に入れる「子ボタン」
	childBtn := tendon.NewButton(100, 150, 100, 50, "Child Btn", color.RGBA{200, 80, 80, 255})

	// ★ 追加: 子ボタンのホバー判定を可視化
	childBtn.OnMouseEnter = func(e *tendon.Element) {
		fmt.Println("  🟥 子ボタンにマウスが【入りました】")
	}
	childBtn.OnMouseLeave = func(e *tendon.Element) {
		fmt.Println("  🟥 子ボタンからマウスが【出ました】")
	}
	childBtn.OnLeftClick = func(e *tendon.Element) {
		fmt.Println("【確認】子ボタンがクリックされました！")
	}

	// 3. パネルの子要素として登録
	panel.Children = append(panel.Children, childBtn)

	otherElem := tendon.NewButton(400, 50, 100, 100, "Other", color.RGBA{100, 100, 100, 255})
	otherElem.Z = 2

	game := &TestGame{
		elements: tendon.Elements{panel, otherElem},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon Hover Bug Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
