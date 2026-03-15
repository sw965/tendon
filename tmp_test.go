package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sw965/tendon"
)

type FullRotationDemo struct {
	target *tendon.Element

	// Rotated系 (辺に吸着して回る緑チーム)
	rotatedL, rotatedR, rotatedA, rotatedB *tendon.Element

	// Bounds系 (外枠を避ける赤チーム)
	boundsL, boundsR, boundsA, boundsB *tendon.Element

	angle float64
}

func (g *FullRotationDemo) Update() error {
	g.angle += 0.5
	g.target.SetAngle(g.angle)

	// --- Rotated配置 ---
	g.rotatedL.PlaceLeftOfRotated(g.target, 20, tendon.AlignCenter)
	g.rotatedR.PlaceRightOfRotated(g.target, 20, tendon.AlignCenter)
	g.rotatedA.PlaceAboveRotated(g.target, 20, tendon.AlignCenter)
	g.rotatedB.PlaceBelowRotated(g.target, 20, tendon.AlignCenter)

	// --- Bounds配置 ---
	g.boundsL.PlaceLeftOfBounds(g.target, 20, tendon.AlignCenter)
	g.boundsR.PlaceRightOfBounds(g.target, 20, tendon.AlignCenter)
	g.boundsA.PlaceAboveBounds(g.target, 20, tendon.AlignCenter)
	g.boundsB.PlaceBelowBounds(g.target, 20, tendon.AlignCenter)

	return nil
}

func (g *FullRotationDemo) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{25, 25, 30, 255})

	g.target.Draw(screen)
	// 緑チーム描画
	g.rotatedL.Draw(screen)
	g.rotatedR.Draw(screen)
	g.rotatedA.Draw(screen)
	g.rotatedB.Draw(screen)
	// 赤チーム描画
	g.boundsL.Draw(screen)
	g.boundsR.Draw(screen)
	g.boundsA.Draw(screen)
	g.boundsB.Draw(screen)

	// ラベルの重なりを防ぐために位置を調整
	ebitenutil.DebugPrintAt(screen, "FULL ROTATION LAYOUT DEMO", 20, 20)
	ebitenutil.DebugPrintAt(screen, "[GREEN] Rotated: Sticks to edges and rotates together", 20, 50)
	ebitenutil.DebugPrintAt(screen, "[RED]   Bounds:  Stays upright, follows the bounding box", 20, 70)

	msg := fmt.Sprintf("Angle: %.1f", g.angle)
	ebitenutil.DebugPrintAt(screen, msg, 20, 110)
}

func (g *FullRotationDemo) Layout(w, h int) (int, int) { return 800, 600 }

func TestFullRotationDemo(t *testing.T) {
	// ターゲット(白)
	target := tendon.NewElement()
	target.Image = ebiten.NewImage(160, 80)
	target.Image.Fill(tendon.White)
	target.XRelativeToParent, target.YRelativeToParent = 320, 260
	target.AnchorX, target.AnchorY = 0.5, 0.5

	createMarker := func(c color.Color) *tendon.Element {
		e := tendon.NewElement()
		e.Image = ebiten.NewImage(24, 24)
		e.Image.Fill(c)
		e.AnchorX, e.AnchorY = 0.5, 0.5
		return e
	}

	game := &FullRotationDemo{
		target:   target,
		rotatedL: createMarker(tendon.Green),
		rotatedR: createMarker(tendon.Green),
		rotatedA: createMarker(tendon.Green),
		rotatedB: createMarker(tendon.Green),
		boundsL:  createMarker(tendon.Red),
		boundsR:  createMarker(tendon.Red),
		boundsA:  createMarker(tendon.Red),
		boundsB:  createMarker(tendon.Red),
	}

	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
