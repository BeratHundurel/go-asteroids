package goasteroids

import (
	"go-asteroids/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Shield struct {
	position  Vector
	rotation  float64
	sprite    *ebiten.Image
	shiledObj *resolv.Circle
	game      *GameScene
}

func NewShield(game *GameScene, position Vector, rotation float64) *Shield {
	sprite := assets.ShieldSprite
	halfW, halfH := HalfOfTheImage(sprite)

	position.X -= halfW
	position.Y -= halfH

	shieldObj := resolv.NewCircle(0, 0, halfW)

	s := &Shield{
		position:  position,
		rotation:  rotation,
		sprite:    sprite,
		shiledObj: shieldObj,
		game:      game,
	}

	s.game.space.Add(s.shiledObj)

	return s
}

func (s *Shield) Update() {
	diffX := float64(s.sprite.Bounds().Dx()-s.game.player.sprite.Bounds().Dx()) * 0.5
	diffY := float64(s.sprite.Bounds().Dy()-s.game.player.sprite.Bounds().Dy()) * 0.5

	pos := Vector{
		X: s.game.player.position.X - diffX,
		Y: s.game.player.position.Y - diffY,
	}

	s.position = pos
	s.rotation = s.game.player.rotation
	s.shiledObj.Move(pos.X, pos.Y)
}

func (s *Shield) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(s.sprite)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(s.rotation)
	op.GeoM.Translate(halfW, halfH)
	
	op.GeoM.Translate(s.position.X, s.position.Y)
	
	screen.DrawImage(s.sprite, op)
}
