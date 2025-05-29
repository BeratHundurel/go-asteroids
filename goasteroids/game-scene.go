package goasteroids

import (
	"go-asteroids/assets"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity   = 0.25
	meteorSpawnTime      = 100 * time.Millisecond
	meteorSpeedUpAmount  = 0.1
	meteorSpeedUpTime    = 1000 * time.Millisecond
	cleanUpExplosionTime = 200 * time.Millisecond
	baseBeatWaitTime     = 1600
)

type GameScene struct {
	player               *Player
	baseVelocity         float64
	meteorCount          int
	meteorsSpawnTimer    *Timer
	meteors              map[int]*Meteor
	meteorsForLevel      int
	velocityTimer        *Timer
	space                *resolv.Space
	lasers               map[int]*Laser
	laserCount           int
	score                int
	explosionSmallSprite *ebiten.Image
	explosionSprite      *ebiten.Image
	explosionFrames      []*ebiten.Image
	cleanUpTimer         *Timer
	playerIsDead         bool
	audioContext         *audio.Context
	thrustPlayer         *audio.Player
	exhaust              *Exhaust
	laserOnePlayer       *audio.Player
	laserTwoPlayer       *audio.Player
	laserThirdPlayer     *audio.Player
	explosionPlayer      *audio.Player
	beatOnePlayer        *audio.Player
	beatTwoPlayer        *audio.Player
	beatTimer            *Timer
	beatWaitTime         int
	playBeatOne          bool
}

func NewGameScene() *GameScene {
	g := &GameScene{
		meteorsSpawnTimer:    NewTimer(meteorSpawnTime),
		baseVelocity:         baseMeteorVelocity,
		velocityTimer:        NewTimer(meteorSpeedUpTime),
		meteors:              make(map[int]*Meteor),
		meteorsForLevel:      2,
		meteorCount:          0,
		space:                resolv.NewSpace(ScreenWidth, ScreenHeight, 16, 16),
		lasers:               make(map[int]*Laser),
		laserCount:           0,
		explosionSprite:      assets.ExplosionSprite,
		explosionSmallSprite: assets.ExplosionSmallSprite,
		cleanUpTimer:         NewTimer(cleanUpExplosionTime),
		beatTimer:            NewTimer(2 * time.Second),
		beatWaitTime:         baseBeatWaitTime,
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)
	g.explosionFrames = assets.Explosion
	g.audioContext = audio.NewContext(48000)
	g.thrustPlayer, _ = g.audioContext.NewPlayer(assets.ThrustSound)
	g.laserOnePlayer, _ = g.audioContext.NewPlayer(assets.LaserOneSound)
	g.laserTwoPlayer, _ = g.audioContext.NewPlayer(assets.LaserTwoSound)
	g.laserThirdPlayer, _ = g.audioContext.NewPlayer(assets.LaserThirdSound)
	g.explosionPlayer, _ = g.audioContext.NewPlayer(assets.ExplosionSound)
	g.beatOnePlayer, _ = g.audioContext.NewPlayer(assets.BeatOneSound)
	g.beatTwoPlayer, _ = g.audioContext.NewPlayer(assets.BeatTwoSound)
	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.updateExhaust()

	g.isPlayerDying()

	g.isPlayerDead(state)

	g.spawnMeteors()

	for _, m := range g.meteors {
		m.Update()
	}

	for _, l := range g.lasers {
		l.Update()
	}

	g.speedUpMeteors()

	g.isPlayerCollidingWithMeteor()

	g.isMeteorHitByPlayerLaser()

	g.cleanUp()
	
	g.beatSound()

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}

	for _, m := range g.meteors {
		m.Draw(screen)
	}

	for _, l := range g.lasers {
		l.Draw(screen)
	}
}

func (g *GameScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameScene) beatSound() {
	g.beatTimer.Update()
	if g.beatTimer.IsReady() {
		if g.playBeatOne {
			g.beatOnePlayer.Rewind()
			g.beatOnePlayer.Play()
			g.beatTimer.Reset()
		} else {
			g.beatTwoPlayer.Rewind()
			g.beatTwoPlayer.Play()
			g.beatTimer.Reset()
		}
		g.playBeatOne = !g.playBeatOne

		if g.beatWaitTime > 400 {
			g.beatWaitTime = g.beatWaitTime - 25
			g.beatTimer = NewTimer(time.Millisecond * time.Duration(g.beatWaitTime))
		}
	}
}

func (g *GameScene) updateExhaust() {
	if g.exhaust != nil {
		g.exhaust.Update()
	}
}

func (g *GameScene) isPlayerDying() {
	if g.player.isDying {
		g.player.dyingTimer.Update()
		if g.player.dyingTimer.IsReady() {
			g.player.dyingTimer.Reset()
			g.player.dyingCounter++
			if g.player.dyingCounter == 12 {
				g.player.isDying = false
				g.player.isDead = true
			} else if g.player.dyingCounter < 12 {
				g.player.sprite = g.explosionFrames[g.player.dyingCounter]
			} else {
				//Do nothing
			}
		}
	}
}

func (g *GameScene) isPlayerDead(state *State) {
	if g.playerIsDead {
		g.player.livesRemaining--
		if g.player.livesRemaining == 0 {
			g.Reset()
			state.SceneManager.GoToScene(g)
		}
	}
}

func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, m := range g.meteors {
		for _, l := range g.lasers {
			if m.meteorObj.IsIntersecting(l.laserObj) {
				if m.meteorObj.Tags().Has(TagSmall) {
					m.sprite = g.explosionSmallSprite
					g.score++

					if !g.explosionPlayer.IsPlaying() {
						g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}
				} else {
					oldPos := m.position

					m.sprite = g.explosionSprite

					g.score++

					if !g.explosionPlayer.IsPlaying() {
						g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}

					numToSpawn := rand.Intn(numberOfSmallMeteorsFromLargeMeteor)
					for range numToSpawn {
						meteor := NewSmallMeteor(baseMeteorVelocity, g, len(m.game.meteors)-1)
						meteor.position = Vector{oldPos.X + float64(rand.Intn(100-50)+50), oldPos.Y + float64(rand.Intn(100-50)+50)}
						meteor.meteorObj.SetPosition(meteor.position.X, meteor.position.Y)
						g.space.Add(meteor.meteorObj)

						g.meteorCount++
						g.meteors[m.game.meteorCount] = meteor
					}
				}
			}
		}
	}
}

func (g *GameScene) spawnMeteors() {
	g.meteorsSpawnTimer.Update()
	if g.meteorsSpawnTimer.IsReady() {
		g.meteorsSpawnTimer.Reset()
		if len(g.meteors) < g.meteorsForLevel && g.meteorCount < g.meteorsForLevel {
			m := NewMeteor(g.baseVelocity, g, len(g.meteors)-1)
			g.space.Add(m.meteorObj)
			g.meteorCount++
			g.meteors[g.meteorCount] = m
		}
	}
}

func (g *GameScene) speedUpMeteors() {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}
}

func (g *GameScene) isPlayerCollidingWithMeteor() {
	for _, m := range g.meteors {
		if m.meteorObj.IsIntersecting(g.player.playerObj) {
			if !g.player.isShielded {
				m.game.player.isDying = true

				if !g.explosionPlayer.IsPlaying() {
					g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
				break
			} else {

			}
		}
	}
}

func (g *GameScene) cleanUp() {
	g.cleanUpTimer.Update()
	if g.cleanUpTimer.IsReady() {
		for i, m := range g.meteors {
			if m.sprite == g.explosionSprite || m.sprite == g.explosionSmallSprite {
				delete(g.meteors, i)
				g.space.Remove(m.meteorObj)
			}
		}
		g.cleanUpTimer.Reset()
	}
}

func (g *GameScene) Reset() {
	g.player = NewPlayer(g)
	g.meteors = make(map[int]*Meteor)
	g.meteorCount = 0
	g.meteorsSpawnTimer.Reset()
	g.lasers = make(map[int]*Laser)
	g.laserCount = 0
	g.score = 0
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
	g.playerIsDead = false
	g.exhaust = nil
	g.space.RemoveAll()
	g.space.Add(g.player.playerObj)
}
