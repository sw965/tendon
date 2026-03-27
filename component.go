package tendon

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	// TODO 戻り値にerrorを返すようにする？
	Update()
	// TODO 戻り値にerrorを返すようにする？
	Draw(screen *ebiten.Image)
	BaseElement() *Element
}

type Components []Component
