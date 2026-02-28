package tendon_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type TestGame struct {
	elements []*tendon.Element
}

func (g *TestGame) Update() error {
	for _, e := range g.elements {
		e.Update(0, 0)
	}
	return nil
}

func (g *TestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	for _, e := range g.elements {
		e.Draw(screen)
	}
}

func (g *TestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 400
}

func TestInteractive(t *testing.T) {
	roots := []*tendon.Element{
		tendon.NewButton(50, 50, 100, 100, "OK", color.RGBA{60, 60, 60, 255}),
		tendon.NewButton(25, 25, 50, 50, "NO", color.RGBA{150, 20, 20, 255}),
		tendon.NewButton(25, 25, 50, 50, "WTF", color.RGBA{150, 20, 20, 255}),
	}

	for i := 1; i < len(roots); i++ {
		roots[i].PlaceRightOf(roots[i-1], 10.0, tendon.AlignCenter)
	}

	game := &TestGame{elements:roots}
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}