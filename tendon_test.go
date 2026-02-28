package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon" // 自身のパスに合わせて調整してください
)

func TestElement_OnUpdate(t *testing.T) {
	isCalled := false
	el := &tendon.Element{
		Visible: true,
		OnUpdate: func(e *tendon.Element) {
			isCalled = true
		},
	}
	el.Update(0, 0)
	if !isCalled {
		t.Errorf("Expected OnUpdate to be called")
	}
}

type TestGame struct {
	root *tendon.Element
}

func (g *TestGame) Update() error {
	g.root.Update(0, 0)
	return nil
}

func (g *TestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	g.root.Draw(screen)
}

func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 400
}

func TestInteractive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping interactive test")
	}

	root := &tendon.Element{Visible: true}

	// ドラッグ可能な親パネル
	panel := tendon.NewButton(50, 50, 300, 300, "Parent Panel", color.RGBA{60, 60, 60, 255})
	panel.Draggable = true

	// ボタン1: 押した瞬間 (JustPressed)
	btn1 := tendon.NewButton(20, 40, 150, 40, "Press Event", color.RGBA{180, 50, 50, 255})
	btn1.OnLeftClick = func(e *tendon.Element) {
		fmt.Println(">> Button 1: Pressed (Down)!")
	}

	// ボタン2: 離した瞬間 (Released)
	btn2 := tendon.NewButton(20, 100, 150, 40, "Release Event", color.RGBA{50, 100, 200, 255})
	btn2.OnLeftReleased = func(e *tendon.Element) {
		fmt.Println(">> Button 2: Released (Up)!")
	}

	// ボタン3: 押しっぱなし & 離した時の合わせ技
	btn3 := tendon.NewButton(20, 160, 150, 40, "Hold & Release", color.RGBA{50, 150, 50, 255})
	btn3.OnLeftClick = func(e *tendon.Element) {
		fmt.Println(">> Button 3: Charge Start!")
	}
	btn3.OnLeftPressed = func(e *tendon.Element) {
		// 押し続けている間ドットを表示
		fmt.Print(".")
	}
	btn3.OnLeftReleased = func(e *tendon.Element) {
		fmt.Println("\n>> Button 3: Fire!")
	}

	panel.Children = append(panel.Children, btn1, btn2, btn3)
	root.Children = append(root.Children, panel)

	game := &TestGame{root: root}
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}