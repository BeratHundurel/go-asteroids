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

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	cm := colorm.ColorM{}
	cm.Scale(1, 1, 1, 0.2)
	op.GeoM.Translate(l.position.X, l.position.Y)

	colorm.DrawImage(screen, l.sprite, cm, op)
}

func (l *LifeIndicator) Update() {
	// Update logic for the life indicator can be added here if needed
}
