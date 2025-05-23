package goasteroids

import (
	"go-asteroids/assets"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationSpeedMin                    = -0.02
	rotationSpeedMax                    = 0.02
	numberOfSmallMeteorsFromLargeMeteor = 4
)

type Meteor struct {
	game          *GameScene
	position      Vector
	movement      Vector
	rotation      float64
	angle         float64
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewMeteor(baseVelocity float64, g *GameScene, index int) *Meteor {
	target := Vector{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	angle := rand.Float64() * 2 * math.Pi

	r := ScreenWidth/2.0 + 500

	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	velocity := baseVelocity + rand.Float64()*1.5

	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

	m := &Meteor{
		game:          g,
		position:      pos,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		angle:         angle,
		sprite:        sprite,
	}

	return m
}

func (m *Meteor) Update() {
	dx := m.movement.X
	dy := m.movement.Y

	m.position.X += dx
	m.position.Y += dy
	m.rotation += m.rotationSpeed
	
	m.keepOnScreen()

}

func (m *Meteor) Draw(screen *ebiten.Image) {
	halW, halfH := HalfOfTheImage(m.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfH, -halW)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfH, halW)
	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)

}

func (m *Meteor) keepOnScreen() {
	if m.position.X >= ScreenWidth{
		m.position.X = 0
	}
	if m.position.X < 0 {
		m.position.X = ScreenWidth
	}
	if m.position.Y >= ScreenHeight {
		m.position.Y = 0
	}
	if m.position.Y < 0 {
		m.position.Y = ScreenHeight
	}
}