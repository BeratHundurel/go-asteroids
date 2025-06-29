package goasteroids

import (
	"go-asteroids/assets"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

const (
	maxAcceleration      = 8.0
	rotationPerSecond    = math.Pi
	ScreenWidth          = 1280
	ScreenHeight         = 720
	shootCoolDown        = time.Millisecond * 150
	burstCoolDown        = time.Millisecond * 500
	laserSpawnOffset     = 50.0
	maxShotsPerBurst     = 3
	dyingAnimationAmount = 50 * time.Millisecond
	numberOfLives        = 3
	numberOfShields      = 3
	shieldDuration       = time.Second * 6
	hyperSpaceCooldown   = time.Second * 10
	driftTime            = time.Second * 30
)

var curAcceleration float64
var shotsFired = 0

type Player struct {
	game                *GameScene
	sprite              *ebiten.Image
	rotation            float64
	position            Vector
	velocity            float64
	playerObj           *resolv.Circle
	shootCoolDown       *Timer
	burstCoolDown       *Timer
	isShielded          bool
	isDying             bool
	isDead              bool
	dyingTimer          *Timer
	dyingCounter        int
	livesRemaining      int
	lifeIndicators      []*LifeIndicator
	shieldTimer         *Timer
	shieldRemaining     int
	shieldIndicators    []*ShieldIndicator
	hyperSpaceIndicator *HyperSpaceIndicator
	hyperSpaceTimer     *Timer
	driftTimer          *Timer
	driftAngle          float64
}

func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSprite

	halfW, halfH := HalfOfTheImage(sprite)

	pos := Vector{
		X: ScreenWidth/2 - halfW,
		Y: ScreenHeight/2 - halfH,
	}

	playerObj := resolv.NewCircle(pos.X, pos.Y, float64(sprite.Bounds().Dx())/2)

	var lifeIndicators []*LifeIndicator
	var xPosition = 20.0
	for range numberOfLives {
		li := NewLifeIndicator(Vector{X: xPosition, Y: 20}, 0)
		lifeIndicators = append(lifeIndicators, li)
		xPosition += 50.0
	}

	var shieldIndicators []*ShieldIndicator
	xPosition = 45.0
	for range numberOfShields {
		si := NewShieldIndicator(Vector{X: xPosition, Y: 60})
		shieldIndicators = append(shieldIndicators, si)
		xPosition += 50.0
	}

	p := &Player{
		sprite:              sprite,
		game:                game,
		position:            pos,
		playerObj:           playerObj,
		shootCoolDown:       NewTimer(shootCoolDown),
		burstCoolDown:       NewTimer(burstCoolDown),
		isShielded:          false,
		isDying:             false,
		isDead:              false,
		dyingTimer:          NewTimer(dyingAnimationAmount),
		dyingCounter:        0,
		livesRemaining:      3,
		lifeIndicators:      lifeIndicators,
		shieldRemaining:     numberOfShields,
		shieldIndicators:    shieldIndicators,
		hyperSpaceIndicator: NewHyperSpaceIndicator(Vector{X: 37.0, Y: 95.0}),
		hyperSpaceTimer:     nil,
		driftTimer:          nil,
	}

	p.playerObj.SetPosition(pos.X, pos.Y)
	p.playerObj.Tags().Set(TagPlayer)

	return p
}

func (p *Player) Draw(screen *ebiten.Image) {
	halfW, halfH := HalfOfTheImage(p.sprite)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	p.isPlayerDead()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	p.accelerate()

	p.useShield()

	p.isDoneAccelerating()

	p.reverse()

	p.isDoneReversing()

	p.isPlayerDrifting()

	p.isDriftFinished()

	p.updateExhaustSprite()

	p.playerObj.SetPosition(p.position.X, p.position.Y)

	p.burstCoolDown.Update()

	p.shootCoolDown.Update()

	p.fireLasers()

	p.hyperSpace()

	if p.hyperSpaceTimer != nil {
		p.hyperSpaceTimer.Update()
	}
}

func (p *Player) isPlayerDrifting() {
	if p.driftTimer != nil {
		p.keepOnScreen()

		p.driftTimer.Update()

		decelerationSpeed := p.velocity / float64(ebiten.TPS()) * 4

		p.position.X += math.Sin(p.driftAngle) * decelerationSpeed
		p.position.Y += math.Cos(p.driftAngle) * -decelerationSpeed
		p.playerObj.SetPosition(p.position.X, p.position.Y)
	}
}

func (p *Player) isDriftFinished() {
	if p.driftTimer != nil && p.driftTimer.IsReady() {
		p.driftTimer = nil
		p.velocity = 0
	}
}

func (p *Player) hyperSpace() {
	if ebiten.IsKeyPressed(ebiten.KeyH) && (p.hyperSpaceTimer == nil || p.hyperSpaceTimer.IsReady()) {
		var randX, randY int
		for {
			randX = rand.Intn(ScreenWidth)
			randY = rand.Intn(ScreenHeight)

			collision := p.game.checkCollision(p.playerObj, nil)
			if !collision {
				break
			}
		}
		p.position.X = float64(randX)
		p.position.Y = float64(randY)

		if p.hyperSpaceTimer == nil {
			p.hyperSpaceTimer = NewTimer(hyperSpaceCooldown)
		}
		p.hyperSpaceTimer.Reset()
	}
}

func (p *Player) useShield() {
	if ebiten.IsKeyPressed(ebiten.KeyS) && p.shieldRemaining > 0 && !p.isShielded {
		if !p.game.shieldsUpPlayer.IsPlaying() {
			p.game.shieldsUpPlayer.Rewind()
			p.game.shieldsUpPlayer.Play()
		}

		p.isShielded = true
		p.shieldTimer = NewTimer(shieldDuration)
		p.game.shield = NewShield(p.game, Vector{}, p.rotation)
		p.shieldRemaining--
		p.shieldIndicators = p.shieldIndicators[:len(p.shieldIndicators)-1]
	}

	if p.shieldTimer != nil && p.isShielded {
		p.shieldTimer.Update()
	}

	if p.shieldTimer != nil && p.shieldTimer.IsReady() {
		p.shieldTimer = nil
		p.isShielded = false
		if p.game.shield != nil {
			p.game.space.Remove(p.game.shield.shiledObj)
			p.game.shield = nil
		}
	}
}

func (p *Player) isPlayerDead() {
	if p.isDead {
		p.game.playerIsDead = true
	}
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

				switch shotsFired {
				case 1:
					if !p.game.laserOnePlayer.IsPlaying() {
						p.game.laserOnePlayer.Rewind()
						p.game.laserOnePlayer.Play()
					}
				case 2:
					if !p.game.laserTwoPlayer.IsPlaying() {
						p.game.laserTwoPlayer.Rewind()
						p.game.laserTwoPlayer.Play()
					}
				case 3:
					if !p.game.laserThirdPlayer.IsPlaying() {
						p.game.laserThirdPlayer.Rewind()
						p.game.laserThirdPlayer.Play()
					}
				}
			} else {
				p.burstCoolDown.Reset()
				shotsFired = 0
			}
		}
	}
}

func (p *Player) isDoneAccelerating() {
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Pause()
		}

		if p.velocity < curAcceleration*10 {
			p.velocity = curAcceleration*10 - 5.0
		}

		if p.velocity < 0 {
			p.velocity = 0
		}

		curAcceleration = 0

		p.driftTimer = NewTimer(driftTime)

		p.driftAngle = p.rotation
	}
}

func (p *Player) updateExhaustSprite() {
	if !ebiten.IsKeyPressed(ebiten.KeyUp) && !ebiten.IsKeyPressed(ebiten.KeyDown) && p.game.exhaust != nil {
		p.game.exhaust = nil
	}
}

func (p *Player) reverse() {
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.driftTimer = nil

		p.keepOnScreen()

		dx := math.Sin(p.rotation) * -3
		dy := math.Cos(p.rotation) * 3

		halfW, halfH := HalfOfTheImage(p.sprite)

		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)*-exhaustSpawnOffset,
			p.position.Y + halfH + math.Cos(p.rotation)*exhaustSpawnOffset,
		}

		p.game.exhaust = NewExhaust(spawnPos, p.rotation+180.0*math.Pi/180.0)
		p.position.X += dx
		p.position.Y += dy

		p.playerObj.SetPosition(p.position.X, p.position.Y)

		if !p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Rewind()
			p.game.thrustPlayer.Play()
		}
	}
}

func (p *Player) isDoneReversing() {
	if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		if p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Pause()
		}
	}
}

func (p *Player) accelerate() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.driftTimer = nil

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

		halfW, halfH := HalfOfTheImage(p.sprite)

		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)*exhaustSpawnOffset,
			p.position.Y + halfH + math.Cos(p.rotation)*-exhaustSpawnOffset,
		}

		p.game.exhaust = NewExhaust(spawnPos, p.rotation+180.0*math.Pi/180.0)

		p.position.X += dx
		p.position.Y += dy

		if !p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Rewind()
			p.game.thrustPlayer.Play()
		}
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
