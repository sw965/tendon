package tendon_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sw965/tendon"
)

type BoxDragGame struct {
	box *tendon.Box
}

func (g *BoxDragGame) Update() error {
	// 1. 基本更新
	g.box.Update()

	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	// 2. ドラッグ開始判定
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var hits tendon.Components
		g.box.FindAllFromPoint(fmx, fmy, &hits)
		for _, h := range hits {
			el := h.BaseElement()
			if el != g.box.BaseElement() {
				el.SetZ(100) // ドラッグ中は手前に表示
				el.StartDrag()
				break
			}
		}
	}

	// 3. ドラッグ終了判定
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, child := range g.box.LayoutChildren {
			el := child.BaseElement()
			if el.IsDragging() {
				el.SetZ(0) // Z値を戻す
				el.StopDrag()
				// ★ 手を離した瞬間に整列を呼び出すことで、
				// ドラッグされていた要素が元のスロットに戻る
				g.box.Reflow()
			}
		}
	}

	// 4. ドラッグ移動の実行
	g.box.UpdateDragMove()

	return nil
}

func (g *BoxDragGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.box.Draw(screen)
	ebitenutil.DebugPrint(screen, "Drag a card and release it to see it snap back.")
}

func (g *BoxDragGame) Layout(w, h int) (int, int) { return 800, 600 }

func TestBoxDragReleaseSnap(t *testing.T) {
	box := tendon.NewBox(700, 200, 20)
	box.XRelativeToParent, box.YRelativeToParent = 50, 200
	box.Image = ebiten.NewImage(700, 200)
	box.Image.Fill(color.RGBA{50, 50, 55, 255})

	colors := []color.RGBA{tendon.Red, tendon.Green, tendon.Blue, tendon.Orange}
	for _, c := range colors {
		btn, _ := tendon.NewButton(120, 150, "Card", c)
		btn.Draggable = true
		box.AppendChild(btn)
	}

	box.Reflow()

	game := &BoxDragGame{box: box}
	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
