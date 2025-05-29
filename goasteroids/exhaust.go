package goasteroids

import (
	"go-asteroids/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const exhaustSpawnOffset = -50.0

type Exhaust struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewExhaust(position Vector, rotation float64) *Exhaust {
	sprite := assets.ExhaustSprite
	halfW, halfH := HalfOfTheImage(sprite)
	position.X -= halfW
	position.Y -= halfH

	return &Exhaust{
		position: position,
		rotation: rotation,
		sprite: sprite,
	}
}

func (e *Exhaust) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(assets.ExhaustSprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(e.rotation)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(e.position.X, e.position.Y)

	screen.DrawImage(e.sprite, op)
}

func (e *Exhaust) Update() {
	speed := maxAcceleration / float64(ebiten.TPS())
	e.position.X += math.Sin(e.rotation) * speed
	e.position.Y += math.Cos(e.rotation) * -speed
}


