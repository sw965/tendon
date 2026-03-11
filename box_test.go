package tendon_test

import (
	"image/color"
	_ "image/jpeg"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sw965/tendon"
)

var dmgImg *ebiten.Image

func init() {
	return
	img, _, err := ebitenutil.NewImageFromFile("")
	if err != nil {
		panic(err)
	}
	dmgImg = img
}

type DynamicBoxGame struct {
	box *tendon.Box
}

func (g *DynamicBoxGame) Update() error {
	// 左クリック：カードの削除
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		var hits tendon.Components // ★ Elements から変更

		tendon.Components{g.box}.FindAllFromPoint(float64(cx), float64(cy), &hits)

		for _, e := range hits {
			// ★ eはインターフェースなので、ベースの実体(BaseElement)同士で比較する
			if e.BaseElement() != g.box.BaseElement() {
				g.box.RemoveChild(e)
				g.box.Update() // ★ 引数削除
				break
			}
		}
	}

	// 右クリック：カードの追加
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		card := tendon.NewElement()
		card.Image = dmgImg
		card.SetScale(0.1)
		g.box.AppendChild(card)
		g.box.Update() // ★ 引数削除
	}

	return nil
}

func (g *DynamicBoxGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.box.Draw(screen)

	ebitenutil.DebugPrint(screen, "Left Click: Remove Card / Right Click: Add Card")
}

func (g *DynamicBoxGame) Layout(w, h int) (int, int) {
	return 1000, 600
}

func TestBoxDynamic(t *testing.T) {
	return
	hBox := tendon.NewBox(900, 400, 10)
	hBox.XRelativeToParent = 50
	hBox.YRelativeToParent = 100
	hBox.MainAlignment = tendon.AlignCenter
	hBox.Image.Fill(color.RGBA{50, 50, 55, 255})

	for i := 0; i < 5; i++ {
		card := tendon.NewElement()
		card.Image = dmgImg
		card.SetScale(0.1)
		hBox.AppendChild(card)
	}
	hBox.Orientation = tendon.Horizontal // ★ 念のため指定
	hBox.Update()                        // ★ 引数削除

	game := &DynamicBoxGame{
		box: hBox,
	}

	ebiten.SetWindowSize(1000, 600)
	ebiten.SetWindowTitle("Tendon Dynamic Box Demo")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}