package tendon_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type ColorTestGame struct {
	elements tendon.Elements
}

func (g *ColorTestGame) Update() error {
	return nil
}

func (g *ColorTestGame) Draw(screen *ebiten.Image) {
	// 背景を少し暗めのグレーにして、各色を際立たせる
	screen.Fill(color.RGBA{45, 45, 50, 255})
	g.elements.Draw(screen)
}

func (g *ColorTestGame) Layout(w, h int) (int, int) {
	return 850, 500
}

func TestColorPaletteExtended(t *testing.T) {
	type colorInfo struct {
		name string
		c    color.RGBA
	}

	// 全16色のリスト
	palette := []colorInfo{
		{"Red", tendon.Red},
		{"Blue", tendon.Blue},
		{"Yellow", tendon.Yellow},
		{"Green", tendon.Green},
		{"Orange", tendon.Orange},
		{"Brown", tendon.Brown},
		{"Pink", tendon.Pink},
		{"LightBlue", tendon.LightBlue},
		{"YellowGreen", tendon.YellowGreen},
		{"Purple", tendon.Purple},
		{"PaleOrange", tendon.PaleOrange},
		{"Ocher", tendon.Ocher},
		{"Vermilion", tendon.Vermilion},
		{"Gray", tendon.Gray},
		{"Black", tendon.Black},
		{"White", tendon.White},
	}

	// 4列 × 4行 のグリッドを作成
	grid := tendon.NewGrid(4, 4, 180, 80, 20)
	grid.XRelativeToParent = 40
	grid.YRelativeToParent = 40

	for i, info := range palette {
		row := i / 4
		col := i % 4
		cell := grid.GetCell(row, col)
		if cell == nil {
			continue
		}

		// ボタンとして配置（ラベル付きの矩形として利用）
		btn := tendon.NewButton(0, 0, 180, 80, info.name, info.c)
		cell.AppendChild(btn)
	}

	game := &ColorTestGame{
		elements: tendon.Elements{grid.Element},
	}

	ebiten.SetWindowSize(850, 500)
	ebiten.SetWindowTitle("Tendon Extended Color Palette")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
