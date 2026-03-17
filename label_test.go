package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sw965/tendon"
)

type LabelTestGame struct {
	root    *tendon.Element
	counter int
	dynamic *tendon.Label
}

func (g *LabelTestGame) Update() error {
	g.counter++

	if g.counter%60 == 0 {
		newSize := float64(16 + ((g.counter/60)%3)*16)
		newText := fmt.Sprintf("Size: %.0fpx / Frame: %d", newSize, g.counter)

		g.dynamic.SetText(newText, g.dynamic.Font().Source, newSize)
	}
	return nil
}

func (g *LabelTestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.root.Draw(screen)
}

func (g *LabelTestGame) Layout(w, h int) (int, int) {
	return 640, 480
}

func TestLabelRefactored(t *testing.T) {
	return
	staticLabel, err := tendon.NewLabel("Static Label (Refactored)", 24)
	if err != nil {
		t.Fatal(err)
	}
	staticLabel.XRelativeToParent = 50
	staticLabel.YRelativeToParent = 50

	dynamicLabel, err := tendon.NewLabel("Dynamic Label", 16)
	if err != nil {
		t.Fatal(err)
	}
	dynamicLabel.XRelativeToParent = 50
	dynamicLabel.YRelativeToParent = 150

	scaledLabel, err := tendon.NewLabel("Scaled 16px (Should be Blurry)", 16)
	if err != nil {
		t.Fatal(err)
	}
	scaledLabel.XRelativeToParent = 50
	scaledLabel.YRelativeToParent = 250
	scaledLabel.SetScale(2.5)

	root := tendon.NewElement()
	root.AppendChild(staticLabel)
	root.AppendChild(dynamicLabel)
	root.AppendChild(scaledLabel)

	game := &LabelTestGame{
		root:    root,
		dynamic: dynamicLabel,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon Refactored Label Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
