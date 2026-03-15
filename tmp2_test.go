package tendon_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sw965/tendon"
)

type CollisionTestGame struct {
	rectEl *tendon.Element
	cursor *tendon.Element
	angle  float64
	// 色の状態管理用フラグ
	lastContains bool
	lastOverlap  bool
}

func (g *CollisionTestGame) Update() error {
	// 1. 回転の更新
	g.angle += 0.5
	g.rectEl.SetAngle(g.angle)

	// 2. カーソル位置の更新
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	// 【修正点1】マウスの先端ではなく、カーソルの「中心」がマウス位置に来るように調整
	// これにより、見た目と判定の位置が完全に一致します
	g.cursor.XRelativeToParent = fmx - (g.cursor.BaseWidth() * g.cursor.AnchorX)
	g.cursor.YRelativeToParent = fmy - (g.cursor.BaseHeight() * g.cursor.AnchorY)

	// --- 3. Contains (点判定: 数学的に完璧) の検証 ---
	isContained := g.rectEl.Contains(fmx, fmy)
	if isContained != g.lastContains {
		if isContained {
			// マウスが完全に乗っている間は青白く光る
			g.rectEl.Image.Fill(color.RGBA{180, 200, 255, 255})
		} else {
			g.rectEl.Image.Fill(tendon.White)
		}
		g.lastContains = isContained
	}

	// --- 4. Overlaps (円近似判定: ぶつかり合い用) の検証 ---
	isOverlapped := g.cursor.Overlaps(g.rectEl)
	if isOverlapped != g.lastOverlap {
		// 【修正点2】色を切り替える前に古い画像を破棄して作り直す
		// AsCircleは既存の画像に色を上塗りするため、破棄しないと前の色が残ってしまいます
		if g.cursor.Image != nil {
			g.cursor.Image.Dispose()
		}
		g.cursor.Image = ebiten.NewImage(12, 12)

		if isOverlapped {
			// 円同士が触れたら赤
			g.cursor.AsCircle(tendon.Red, 0)
		} else {
			// 離れたら緑
			g.cursor.AsCircle(tendon.Green, 0)
		}
		g.lastOverlap = isOverlapped
	}

	return nil
}

func (g *CollisionTestGame) Draw(screen *ebiten.Image) {
	// 背景を少し暗く
	screen.Fill(color.RGBA{30, 30, 35, 255})

	// 1. メインの白い長方形を描画
	g.rectEl.Draw(screen)

	// 2. 補助線 (BoundingBox/AABB) を描画
	minX, minY, maxX, maxY := g.rectEl.BoundingBox()
	// 【修正点3】色の指定を color.RGBA{20, 20, 20, 20} に変更
	// Ebitenの乗算済みアルファの仕様により、RGB値がアルファ値を超えると白く光りすぎてしまいます。
	// ここを修正することで、中身の長方形が透けて見えるようになります。
	ebitenutil.DrawRect(screen, minX, minY, maxX-minX, maxY-minY, color.RGBA{20, 20, 20, 20})

	// 3. カーソル（小円）を一番上に描画
	g.cursor.Draw(screen)

	// ステータス表示
	msg := "PRECISION COLLISION TEST\n\n"
	msg += fmt.Sprintf("Contains (Point):  %v -> Highlights BLUE\n", g.lastContains)
	msg += fmt.Sprintf("Overlaps (Circle): %v -> Cursor turns RED\n\n", g.lastOverlap)
	msg += "[Settings]\n"
	msg += " - Approx: 12x4 Circles (Fine)\n"
	msg += " - Protrusion: 0.2 (Covers corners)\n"
	msg += " - Overlap: 0.1\n\n"
	msg += "Watch the corners: Blue (Contained) vs Red (Overlapped)"
	ebitenutil.DebugPrint(screen, msg)
}

func (g *CollisionTestGame) Layout(w, h int) (int, int) { return 800, 600 }

func TestCollisionPrecision(t *testing.T) {
	// ターゲットとなる長方形
	rectEl := tendon.NewElement()
	rectEl.Image = ebiten.NewImage(240, 80)
	rectEl.Image.Fill(tendon.White)
	rectEl.XRelativeToParent, rectEl.YRelativeToParent = 400, 300
	rectEl.AnchorX, rectEl.AnchorY = 0.5, 0.5

	// 円近似の精度設定
	rectEl.SetRectApprox(12, 4, 0.1, 0.2)

	// 判定用のカーソル
	cursor := tendon.NewElement()
	cursor.Image = ebiten.NewImage(12, 12)
	cursor.AsCircle(tendon.Green, 0)
	cursor.AnchorX, cursor.AnchorY = 0.5, 0.5

	game := &CollisionTestGame{
		rectEl: rectEl,
		cursor: cursor,
	}

	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}