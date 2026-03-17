package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sw965/tendon"
)

type DragAndRotateGame struct {
	parent *tendon.Element
	child  *tendon.Element
}

func (g *DragAndRotateGame) Update() error {
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	// 1. 当たり判定（再帰的に全階層を探す）
	var hits tendon.Components
	g.parent.FindAllFromPoint(fmx, fmy, &hits)
	g.parent.UpdateHover(hits)

	// 2. 左クリックでドラッグ開始
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if len(hits) > 0 {
			// 一番上に重なっている要素をドラッグ開始
			hits[0].BaseElement().StartDrag()
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// 全てのドラッグを停止
		g.parent.StopAllDrag()
	}

	// 3. 右クリックで回転（15度ずつ）
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if len(hits) > 0 {
			e := hits[0].BaseElement()
			e.SetAngle(e.GetAngle() + 15)
		}
	}

	// 4. ドラッグ移動の実行（再帰的に処理される）
	g.parent.UpdateDragMove()

	return nil
}

func (g *DragAndRotateGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})

	// 親と子を描画（親子関係があるので parent.Draw だけで子も描画される）
	g.parent.Draw(screen)

	// デバッグ情報の表示
	msg := "SCENE GRAPH TEST\n"
	msg += "--------------------------\n"
	msg += "Left Click  : Drag Move\n"
	msg += "Right Click : Rotate (+15deg)\n\n"

	msg += fmt.Sprintf("[Parent]\n Angle: %.0f (Abs: %.0f)\n",
		g.parent.GetAngle(), g.parent.GetAbsAngle())
	msg += fmt.Sprintf("[Child]\n Angle: %.0f (Abs: %.0f)\n",
		g.child.GetAngle(), g.child.GetAbsAngle())

	ebitenutil.DebugPrint(screen, msg)
}

func (g *DragAndRotateGame) Layout(w, h int) (int, int) { return 800, 600 }

func TestDragAndRotateHierarchy(t *testing.T) {
	// 親要素：グレーの大きな板
	parent := tendon.NewElement()
	parent.Image = ebiten.NewImage(300, 300)
	parent.Image.Fill(color.RGBA{60, 60, 65, 255})
	parent.XRelativeToParent, parent.YRelativeToParent = 400, 300
	parent.AnchorX, parent.AnchorY = 0.5, 0.5
	parent.Draggable = true

	// 子要素：赤い小さな正方形
	child := tendon.NewElement()
	child.Image = ebiten.NewImage(80, 80)
	child.Image.Fill(tendon.Red)
	// 親の中心から少し右上に配置
	child.XRelativeToParent, child.YRelativeToParent = 100, -100
	child.AnchorX, child.AnchorY = 0.5, 0.5
	child.Draggable = true

	// シーングラフの構築
	parent.AppendChild(child)

	game := &DragAndRotateGame{
		parent: parent,
		child:  child,
	}

	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
