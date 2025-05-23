package goasteroids

import (
	"go-asteroids/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxAcceleration   = 8.0
	rotationPerSecond = math.Pi
	ScreenWidth       = 1280
	ScreenHeight      = 720
)

var curAcceleration float64

type Player struct {
	game     *GameScene
	sprite   *ebiten.Image
	rotation float64
	position Vector
	velocity float64
}

func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSprite

	halfW, halfH := HalfOfTheImage(sprite)

	pos := Vector{
		X: ScreenWidth/2 - halfW,
		Y: ScreenHeight/2 - halfH,
	}

	p := &Player{
		sprite:   sprite,
		game:     game,
		position: pos,
	}

	return p
}

func (p *Player) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(p.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfH, -halfW)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	p.accelerate()
}

func (p *Player) accelerate() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.keepOnScreen()

		if curAcceleration < maxAcceleration {
			curAcceleration = p.velocity + 4
		}

		if curAcceleration >= maxAcceleration {
			curAcceleration = maxAcceleration
		}

		p.velocity = curAcceleration

		dx := math.Sin(p.rotation) * curAcceleration
		dy := math.Cos(p.rotation) * -curAcceleration

		p.position.X += dx
		p.position.Y += dy
	}
}

func (p *Player) keepOnScreen() {
	if p.position.X >= float64(ScreenWidth) {
		p.position.X = 0
	}
	if p.position.X < 0 {
		p.position.X = ScreenWidth
	}
	if p.position.Y >= float64(ScreenHeight) {
		p.position.Y = 0
	}
	if p.position.Y < 0 {
		p.position.Y = ScreenHeight
	}
}
