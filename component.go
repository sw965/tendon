package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Update()
	Draw(screen *ebiten.Image)
	BaseElement() *Element
}

type Components []Component