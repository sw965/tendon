package tendon_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type ColorTestGame struct {
	root tendon.Component
}

func (g *ColorTestGame) Update() error {
	return nil
}

func (g *ColorTestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{45, 45, 50, 255})
	g.root.Draw(screen)
}

func (g *ColorTestGame) Layout(w, h int) (int, int) {
	return 850, 500
}

func TestColorPaletteExtended(t *testing.T) {
	return
	type colorInfo struct {
		name string
		c    color.RGBA
	}

	palette := []colorInfo{
		{"Red", tendon.Red}, {"Blue", tendon.Blue}, {"Yellow", tendon.Yellow}, {"Green", tendon.Green},
		{"Orange", tendon.Orange}, {"Brown", tendon.Brown}, {"Pink", tendon.Pink}, {"LightBlue", tendon.LightBlue},
		{"YellowGreen", tendon.YellowGreen}, {"Purple", tendon.Purple}, {"PaleOrange", tendon.PaleOrange}, {"Ocher", tendon.Ocher},
		{"Vermilion", tendon.Vermilion}, {"Gray", tendon.Gray}, {"Black", tendon.Black}, {"White", tendon.White},
	}

	cellW, cellH := 180.0, 80.0
	grid := tendon.NewGrid(4, 4, cellW, cellH, 20)
	grid.XRelativeToParent = 40
	grid.YRelativeToParent = 40

	for i, info := range palette {
		col := i % 4
		row := i / 4
		cell := grid.GetCell(col, row)
		if cell == nil {
			continue
		}

		cell.Image = ebiten.NewImage(int(cellW), int(cellH))
		cell.Image.Fill(info.c)

		l, err := tendon.NewLabel(info.name, 20)
		if err != nil {
			t.Fatal(err)
		}

		cell.AppendChild(l)
		l.PlaceCenterOf(cell)
	}

	game := &ColorTestGame{
		root: grid,
	}

	ebiten.SetWindowSize(850, 500)
	ebiten.SetWindowTitle("Tendon Extended Color Palette")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
