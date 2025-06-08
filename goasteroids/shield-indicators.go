package goasteroids

import (
	"go-asteroids/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type ShieldIndicator struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewShieldIndicator(pos Vector) *ShieldIndicator {
	return &ShieldIndicator{
		position: pos,
		sprite:   assets.ShieldIndicator,
	}
}

func (s *ShieldIndicator) Update() {}

func (s *ShieldIndicator) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(s.sprite)

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)
	op.GeoM.Translate(s.position.X, s.position.Y)
	colorm.DrawImage(screen, s.sprite, cm, op)
}
