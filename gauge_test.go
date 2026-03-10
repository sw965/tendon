package tendon_test

import (
	"fmt"
	"image/color"
	"math/rand"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sw965/tendon"
)

// --- 1. 開発者が独自に作成したカスタム要素 ---

type Particle struct {
	x, y   float64
	vx, vy float64
	life   float64
}

type ExplosiveGauge struct {
	*tendon.Gauge
	baseW     float64
	baseH     float64
	particles []*Particle
}

func NewExplosiveGauge(w, h, max float64, barColor color.Color) (*ExplosiveGauge, error) {
	g, err := tendon.NewLerpGauge(w, h, max, barColor) // 最新の引数に対応
	if err != nil {
		return nil, err
	}
	return &ExplosiveGauge{
		Gauge: g,
		baseW: w,
		baseH: h,
	}, nil
}

func (eg *ExplosiveGauge) Update() {
	oldScale := eg.Bar.WidthScale

	eg.Gauge.Update()

	newScale := eg.Bar.WidthScale

	if newScale < oldScale {
		for i := 0; i < 3; i++ {
			p := &Particle{
				x:    newScale * eg.baseW,
				y:    rand.Float64() * eg.baseH,
				vx:   (rand.Float64()*2 - 1.0) * 4,
				vy:   (rand.Float64()*2 - 1.0) * 4,
				life: 1.0,
			}
			eg.particles = append(eg.particles, p)
		}
	}

	var active []*Particle
	for _, p := range eg.particles {
		p.x += p.vx
		p.y += p.vy
		p.life -= 0.04
		if p.life > 0 {
			active = append(active, p)
		}
	}
	eg.particles = active
}

func (eg *ExplosiveGauge) Draw(screen *ebiten.Image) {
	eg.Gauge.Draw(screen)

	absX, absY := eg.AbsPos()
	absWsc, absHsc := eg.AbsWidthScale(), eg.AbsHeightScale()

	for _, p := range eg.particles {
		c := color.RGBA{255, 120, 0, uint8(255 * p.life)}
		px := absX + (p.x * absWsc)
		py := absY + (p.y * absHsc)
		ebitenutil.DrawRect(screen, px, py, 3, 3, c)
	}
}

// --- 2. テスト用のゲーム ---

type GaugeTestGame struct {
	hpBar   *ExplosiveGauge
	manaBar *tendon.Gauge
	lpText  *tendon.Counter
	value   float64
}

func (g *GaugeTestGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.value -= 20
		if g.value < 0 {
			g.value = 0
		}
		// SetValue ではなく、Target フィールドに直接代入する
		g.hpBar.Target = g.value
		g.manaBar.Target = g.value * 2
		g.lpText.Target = g.value * 80
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.value = 100
		g.hpBar.Target = g.value
		g.manaBar.Target = 200
		g.lpText.Target = 8000
	}

	g.hpBar.Update()
	g.manaBar.Update()
	g.lpText.Update()

	return nil
}

func (g *GaugeTestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 35, 255})

	g.hpBar.Draw(screen)
	g.manaBar.Draw(screen)
	g.lpText.Draw(screen)

	ebitenutil.DebugPrint(screen, "\n\n  Space: Damage / Enter: Heal\n  Red: Lerp (Spark) / Blue: Linear / Text: Lerp Counter")
}

func (g *GaugeTestGame) Layout(w, h int) (int, int) {
	return 640, 480
}

func TestGaugeAnimation(t *testing.T) {
	return
	hp, err := NewExplosiveGauge(300, 30, 100, tendon.Red)
	if err != nil {
		t.Fatal(err)
	}
	hp.XRelativeToParent = 170
	hp.YRelativeToParent = 100

	mana, err := tendon.NewLinearGauge(300, 20, 200, tendon.Blue)
	if err != nil {
		t.Fatal(err)
	}
	mana.XRelativeToParent = 170
	mana.YRelativeToParent = 150

	// 第3引数の 0.1 は不要になったため削除 (内部定数が使われる)
	lp, err := tendon.NewLerpCounter(48, 8000)
	if err != nil {
		t.Fatal(err)
	}
	lp.XRelativeToParent = 170
	lp.YRelativeToParent = 250
	lp.FormatFunc = func(v float64) string {
		return fmt.Sprintf("LP: %.0f", v)
	}

	game := &GaugeTestGame{
		hpBar:   hp,
		manaBar: mana,
		lpText:  lp,
		value:   100,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tendon Gauge & Counter Test")

	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
}
