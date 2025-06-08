package goasteroids

import (
	"go-asteroids/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type HyperSpaceIndicator struct {
	position Vector
	sprite   *ebiten.Image
	rotation float64
}

func NewHyperSpaceIndicator(pos Vector) *HyperSpaceIndicator {
	return &HyperSpaceIndicator{
		position: pos,
		sprite:   assets.HyperSpaceIndicator,
	}
}

func (hsi *HyperSpaceIndicator) Update() {
	// HyperSpaceIndicator does not need to update anything
}

func (hsi *HyperSpaceIndicator) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(hsi.sprite)

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)
	op.GeoM.Translate(hsi.position.X, hsi.position.Y)
	colorm.DrawImage(screen, hsi.sprite, cm, op)
}
