package goasteroids

import (
	"go-asteroids/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	alienLaserSpeedPerSecond = 1000.0
)

type AlienLaser struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
	laserObj *resolv.ConvexPolygon
}

func NewAlienLaser(position Vector, rotation float64) *AlienLaser {
	sprite := assets.AlienLaserSprite

	halfW, halfH := HalfOfTheImage(sprite)

	position.X -= halfW
	position.Y -= halfH

	al := &AlienLaser{
		position: position,
		rotation: rotation,
		sprite: sprite,
		laserObj: resolv.NewRectangle(position.X, position.Y, float64(sprite.Bounds().Dx()), float64(sprite.Bounds().Dy())),
	}
	al.laserObj.SetPosition(position.X, position.Y)
	al.laserObj.Tags().Set(TagLaser)

	return al
}

func (al *AlienLaser) Update() {
	speed := alienLaserSpeedPerSecond / float64(ebiten.TPS())

	al.position.X += math.Sin(al.rotation) * speed
	al.position.Y += math.Cos(al.rotation) * -speed

	al.laserObj.SetPosition(al.position.X, al.position.Y)
}

func (al *AlienLaser) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(al.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(al.rotation)
	op.GeoM.Translate(al.position.X, al.position.Y)

	screen.DrawImage(al.sprite, op)
}
