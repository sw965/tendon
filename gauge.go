package tendon

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// この辺りの命名変更の余地あり
const (
	defaultLinearFrames       = 60.0  // 等速アニメーションが完了するまでのデフォルトフレーム数
	defaultLerpRatio          = 0.1   // 毎フレーム目標に近づく割合 (10%)
	defaultLerpSnapDistFactor = 0.001 // 最大値に対する吸着距離の係数 (0.1%)
	defaultLerpMinSnapDist    = 0.1   // 吸着距離の最小値
	defaultFontSizeFactor     = 0.7   // ゲージの高さに対するフォントサイズの割合
)

var (
	// デフォルトの背景色
	defaultGaugeBgColor = color.RGBA{40, 40, 45, 255}
)

type EasingFunc func(current, target float64) float64

// (ratio*100)ずつ目標値に近づき、snapDictの距離で吸着
func NewLerpEasingFunc(ratio, snapDist float64) EasingFunc {
	return func(current, target float64) float64 {
		if math.Abs(current-target) < snapDist {
			return target
		}
		return current + (target-current)*ratio
	}
}

func NewLinearEasingFunc(speed float64) EasingFunc {
	return func(current, target float64) float64 {
		diff := target - current
		if math.Abs(diff) <= speed {
			return target
		}
		if diff > 0 {
			return current + speed
		}
		return current - speed
	}
}

type Counter struct {
	*Element
	Label      *Label
	Target     float64
	Current    float64
	EasingFunc EasingFunc
	FormatFunc func(float64) string
}

func newCounter(size, init float64) (*Counter, error) {
	l, err := NewLabel("", size)
	if err != nil {
		return nil, err
	}
	c := &Counter{
		Element: NewElement(),
		Label:   l,
		Target:  init,
		Current: init,
		FormatFunc: func(current float64) string {
			return fmt.Sprintf("%.0f", math.Ceil(current))
		},
	}
	c.AppendChild(l.Element)
	c.updateLabel()
	return c, nil
}

func NewLinearCounter(size, init float64) (*Counter, error) {
	c, err := newCounter(size, init)
	if err != nil {
		return nil, err
	}
	speed := init / defaultLinearFrames
	c.EasingFunc = NewLinearEasingFunc(speed)
	return c, nil
}

func NewLerpCounter(size, init float64) (*Counter, error) {
	c, err := newCounter(size, init)
	if err != nil {
		return nil, err
	}
	c.EasingFunc = NewLerpEasingFunc(defaultLerpRatio, math.Max(defaultLerpMinSnapDist, init*defaultLerpSnapDistFactor))
	return c, nil
}

func (c *Counter) Update() {
	c.Element.Update()
	if c.Current == c.Target {
		return
	}
	c.Current = c.EasingFunc(c.Current, c.Target)
	c.updateLabel()
}

func (c *Counter) updateLabel() {
	txt := c.FormatFunc(c.Current)
	c.Label.Update(txt, c.Label.Font().Source, c.Label.Font().Size)
}

type Gauge struct {
	*Counter
	Bar    *Element
	Max    float64
	OnSync func(g *Gauge)
}

func newGauge(w, h, max float64, barColor color.Color) (*Gauge, error) {
	base := NewElement()
	base.Image = ebiten.NewImage(int(w), int(h))
	base.Image.Fill(defaultGaugeBgColor)

	bar := NewElement()
	bar.Image = ebiten.NewImage(int(w), int(h))
	bar.Image.Fill(barColor)
	base.AppendChild(bar)

	fontSize := h * defaultFontSizeFactor
	counter, err := newCounter(fontSize, max)
	if err != nil {
		return nil, err
	}

	g := &Gauge{
		Counter: counter,
		Bar:     bar,
		Max:     max,
		OnSync: func(g *Gauge) {
			g.Label.PlaceCenterOf(g.Element)
		},
	}

	g.Counter.Element = base
	g.AppendChild(counter.Label.Element)
	g.FormatFunc = func(current float64) string {
		return fmt.Sprintf("%.0f / %.0f", math.Ceil(current), g.Max)
	}
	return g, nil
}

func NewLinearGauge(w, h, max float64, barColor color.Color) (*Gauge, error) {
	g, err := newGauge(w, h, max, barColor)
	if err != nil {
		return nil, err
	}
	speed := math.Max(0.1, max/defaultLinearFrames)
	g.EasingFunc = NewLinearEasingFunc(speed)
	g.Refresh()
	return g, nil
}

func NewLerpGauge(w, h, max float64, barColor color.Color) (*Gauge, error) {
	g, err := newGauge(w, h, max, barColor)
	if err != nil {
		return nil, err
	}
	g.EasingFunc = NewLerpEasingFunc(defaultLerpRatio, math.Max(defaultLerpMinSnapDist, max*defaultLerpSnapDistFactor))
	g.Refresh()
	return g, nil
}

func (g *Gauge) Update() {
	old := g.Current
	g.Counter.Update()
	if old != g.Current {
		g.Refresh()
	}
}

func (g *Gauge) Refresh() {
	percent := g.Current / g.Max
	g.Bar.WidthScale = math.Max(0, math.Min(percent, 1.0))
	if g.OnSync != nil {
		g.OnSync(g)
	}
}
