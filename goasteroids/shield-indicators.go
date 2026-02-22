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

	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)
	
	sharedColormDrawOp.GeoM.Reset()
	sharedColormDrawOp.GeoM.Translate(halfW, halfH)
	sharedColormDrawOp.GeoM.Translate(s.position.X, s.position.Y)
	
	colorm.DrawImage(screen, s.sprite, cm, sharedColormDrawOp)
}
