package goasteroids

import (
	"go-asteroids/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type LifeIndicator struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewLifeIndicator(position Vector, rotation float64) *LifeIndicator {
	sprite := assets.LifeIndicator
	return &LifeIndicator{
		position: position,
		rotation: rotation,
		sprite:   sprite,
	}
}

func (l *LifeIndicator) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(l.sprite)

	cm := colorm.ColorM{}
	cm.Scale(1, 1, 1, 0.2)

	sharedColormDrawOp.GeoM.Reset()
	sharedColormDrawOp.GeoM.Translate(halfW, halfH)
	sharedColormDrawOp.GeoM.Translate(l.position.X, l.position.Y)

	colorm.DrawImage(screen, l.sprite, cm, sharedColormDrawOp)
}

func (l *LifeIndicator) Update() {
	// Requirement for the interface
}
