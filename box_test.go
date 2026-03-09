package tendon_test

import (
	"image/color"
	_ "image/jpeg"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil" // マウス判定に必要
	"github.com/sw965/tendon"
)

var dmgImg *ebiten.Image

func init() {
	return
	// 指定されたパスから画像を読み込み
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
		var hits tendon.Elements
		// クリック地点にある要素をすべて取得
		g.box.FindAllFromPoint(float64(cx), float64(cy), &hits)

		// ヒットした要素の中から、Box自身（背景）以外のカードを探して削除
		for _, e := range hits {
			if e != g.box.Element {
				g.box.RemoveChild(e)            // ターゲットをポインタで指定
				g.box.Update(tendon.Horizontal) // 削除後に再整列
				break
			}
		}
	}

	// 右クリック：カードの追加
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		card := tendon.NewElement()
		card.Image = dmgImg
		card.SetScale(0.1)
		g.box.AppendChild(card)         // 新しいカードを追加
		g.box.Update(tendon.Horizontal) // 追加後に再整列
	}

	return nil
}

func (g *DynamicBoxGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})
	g.box.Draw(screen) //

	ebitenutil.DebugPrint(screen, "Left Click: Remove Card / Right Click: Add Card")
}

func (g *DynamicBoxGame) Layout(w, h int) (int, int) {
	return 1000, 600
}

func TestBoxDynamic(t *testing.T) {
	return
	// 1000x400 のゆったりした Box を作成
	hBox := tendon.NewBox(900, 400, 10)
	hBox.XRelativeToParent = 50
	hBox.YRelativeToParent = 100
	hBox.MainAlignment = tendon.AlignCenter
	hBox.Image.Fill(color.RGBA{50, 50, 55, 255}) // 背景色

	// 初期状態で 5枚追加
	for i := 0; i < 5; i++ {
		card := tendon.NewElement()
		card.Image = dmgImg
		card.SetScale(0.1)
		hBox.AppendChild(card)
	}
	hBox.Update(tendon.Horizontal)

	game := &DynamicBoxGame{
		box: hBox,
	}

	ebiten.SetWindowSize(1000, 600)
	ebiten.SetWindowTitle("Tendon Dynamic Box Demo")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
