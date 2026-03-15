package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

func CreateCircularImage(baseImg *ebiten.Image, borderColor color.Color, borderWidth float32) *ebiten.Image {
	bounds := baseImg.Bounds()
	// 画像の幅と高さを取得する。
	width, height := bounds.Dx(), bounds.Dy()

	// baseImgと同じ幅と高さの完全に透明な空の画像を生成する
	canvas := ebiten.NewImage(width, height)
	centerX := float32(width) / 2
	centerY := float32(height) / 2
	// TODO 後でコメントを書く？
	radius := float32(math.Min(float64(width), float64(height))) / 2

	// 四角形の画像を、centerX, centerYを中心として、半径 = w/2 の円の画像に変換にする
	// これによりで、四角形のbaseImgにすっぽりと収まるような円の画像になる
	// またこの時点で指定したborderColorに染まった円になる
	// 最後の引数がtrueであれば、円のフチのギザギザを滑らかにする処理が入る
	vector.DrawFilledCircle(canvas, centerX, centerY, radius, borderColor, true)

	// 再度、baseImgと同じ幅と高さの完全に透明な画像を生成する
	innerCircleImg := ebiten.NewImage(width, height)
	// canvasの円よりも小さめに作るために半径を小さくする
	innerRadius := radius - borderWidth
	// canvasの円よりも小さい円の白の染まった画像を生成する
	vector.DrawFilledCircle(innerCircleImg, centerX, centerY, innerRadius, color.White, true)

	op := &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeSourceIn
	// innerCircleImgに元の画像を貼り付ける (元の画像を円形にする)
	innerCircleImg.DrawImage(baseImg, op)
	// borderColorに染まった円の画像に、innerCircleImgを貼り付ける
	canvas.DrawImage(innerCircleImg, nil)
	return canvas
}
