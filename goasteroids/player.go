package goasteroids

import (
	"go-asteroids/assets"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	maxAcceleration   = 8.0
	rotationPerSecond = math.Pi
	ScreenWidth       = 1280
	ScreenHeight      = 720
	shootCoolDown     = time.Millisecond * 150
	burstCoolDown     = time.Millisecond * 500
	laserSpawnOffset  = 50.0
	maxShotsPerBurst  = 3
)

var curAcceleration float64
var shotsFired int = 0

type Player struct {
	game          *GameScene
	sprite        *ebiten.Image
	rotation      float64
	position      Vector
	velocity      float64
	playerObj     *resolv.Circle
	shootCoolDown *Timer
	burstCoolDown *Timer
}

func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSprite

	halfW, halfH := HalfOfTheImage(sprite)

	pos := Vector{
		X: ScreenWidth/2 - halfW,
		Y: ScreenHeight/2 - halfH,
	}

	playerObj := resolv.NewCircle(pos.X, pos.Y, float64(sprite.Bounds().Dx())/2)

	p := &Player{
		sprite:        sprite,
		game:          game,
		position:      pos,
		playerObj:     playerObj,
		shootCoolDown: NewTimer(shootCoolDown),
		burstCoolDown: NewTimer(burstCoolDown),
	}

	p.playerObj.SetPosition(pos.X, pos.Y)
	p.playerObj.Tags().Set(TagPlayer)

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
	
	p.playerObj.SetPosition(p.position.X, p.position.Y)
	
	p.burstCoolDown.Update()
	
	p.shootCoolDown.Update()
	
	p.fireLasers()
}

func (p *Player) fireLasers() {
	if p.burstCoolDown.IsReady() {
		if p.shootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.shootCoolDown.Reset()
			shotsFired++
			if shotsFired <= maxShotsPerBurst {
				halfW, halfH := HalfOfTheImage(p.sprite)
				
				spawnPos := Vector{
					p.position.X + halfW + math.Sin(p.rotation)*laserSpawnOffset,
					p.position.Y + halfH + math.Cos(p.rotation)*-laserSpawnOffset,
				}
				
				p.game.laserCount++
				laser := NewLaser(spawnPos, p.rotation, p.game.laserCount, p.game)
				p.game.lasers[p.game.laserCount] = laser
				p.game.space.Add(laser.laserObj)
			} else {
				p.burstCoolDown.Reset()
				shotsFired = 0
			}
		}
	}
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
		p.playerObj.SetPosition(0, p.position.Y)
	}
	if p.position.X < 0 {
		p.position.X = ScreenWidth
		p.playerObj.SetPosition(ScreenWidth, p.position.Y)
	}
	if p.position.Y >= float64(ScreenHeight) {
		p.position.Y = 0
		p.playerObj.SetPosition(p.position.X, 0)
	}
	if p.position.Y < 0 {
		p.position.Y = ScreenHeight
		p.playerObj.SetPosition(p.position.X, ScreenHeight)
	}
}
