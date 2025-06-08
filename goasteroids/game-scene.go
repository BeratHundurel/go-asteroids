package goasteroids

import (
	"fmt"
	"go-asteroids/assets"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity   = 0.25
	meteorSpawnTime      = 100 * time.Millisecond
	meteorSpeedUpAmount  = 0.1
	meteorSpeedUpTime    = 1000 * time.Millisecond
	cleanUpExplosionTime = 200 * time.Millisecond
	baseBeatWaitTime     = 1600
	numberOfStars        = 1000
	alienAttackTime      = 3 * time.Second
	alienSpawnTime       = 8 * time.Second
	baseAlienVelocity    = 0.5
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
	stars                []*Star
	currentLevel         int
	shield               *Shield
	shieldsUpPlayer      *audio.Player
	alienAttackTimer     *Timer
	alienCount           int
	alienLaserCount      int
	alienLaserPlayer     *audio.Player
	alienLasers          map[int]*AlienLaser
	alienSoundPLayer     *audio.Player
	aliens               map[int]*Alien
	alienSpawnTimer      *Timer
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
		currentLevel:         1,
		aliens:               make(map[int]*Alien),
		alienCount:           0,
		alienLasers:          make(map[int]*AlienLaser),
		alienLaserCount:      0,
		alienSpawnTimer:      NewTimer(alienSpawnTime),
		alienAttackTimer:     NewTimer(alienAttackTime),
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)
	g.stars = GenerateStars(numberOfStars)

	g.explosionFrames = assets.Explosion
	g.audioContext = audio.NewContext(48000)
	g.thrustPlayer, _ = g.audioContext.NewPlayer(assets.ThrustSound)
	g.laserOnePlayer, _ = g.audioContext.NewPlayer(assets.LaserOneSound)
	g.laserTwoPlayer, _ = g.audioContext.NewPlayer(assets.LaserTwoSound)
	g.laserThirdPlayer, _ = g.audioContext.NewPlayer(assets.LaserThirdSound)
	g.explosionPlayer, _ = g.audioContext.NewPlayer(assets.ExplosionSound)
	g.beatOnePlayer, _ = g.audioContext.NewPlayer(assets.BeatOneSound)
	g.beatTwoPlayer, _ = g.audioContext.NewPlayer(assets.BeatTwoSound)
	g.shieldsUpPlayer, _ = g.audioContext.NewPlayer(assets.ShieldSound)
	g.alienSoundPLayer, _ = g.audioContext.NewPlayer(assets.AlienSound)
	g.alienSoundPLayer.SetVolume(0.5)
	g.alienLaserPlayer, _ = g.audioContext.NewPlayer(assets.AlienLaserSound)
	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.updateExhaust()

	g.updateShield()

	g.isPlayerDying()

	g.isPlayerDead(state)

	g.spawnMeteors()

	g.spawnAliens()

	for _, a := range g.aliens {
		a.Update()
	}

	g.letAliensAttack()

	for _, al := range g.alienLasers {
		al.Update()
	}

	for _, m := range g.meteors {
		m.Update()
	}

	for _, l := range g.lasers {
		l.Update()
	}

	g.speedUpMeteors()

	g.isPlayerCollidingWithMeteor()

	g.isMeteorHitByPlayerLaser()

	g.isPlayerCollidingWithAlien()

	g.isPlayerHitByAlienLaser()

	g.isAlienHitByPlayerLaser()

	g.cleanUp()

	g.beatSound()

	g.isLevelComplete(state)

	g.removeOffScreenAliens()

	g.removeOffScreenLasers()

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	for _, s := range g.stars {
		s.Draw(screen)
	}

	g.player.Draw(screen)

	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}

	if g.shield != nil {
		g.shield.Draw(screen)
	}

	for _, m := range g.meteors {
		m.Draw(screen)
	}

	for _, l := range g.lasers {
		l.Draw(screen)
	}

	if len(g.player.lifeIndicators) > 0 {
		for _, li := range g.player.lifeIndicators {
			li.Draw(screen)
		}
	}

	if len(g.player.shieldIndicators) > 0 {
		for _, si := range g.player.shieldIndicators {
			si.Draw(screen)
		}
	}

	if g.player.hyperSpaceTimer == nil || g.player.hyperSpaceTimer.IsReady() {
		g.player.hyperSpaceIndicator.Draw(screen)
	}

	for _, a := range g.aliens {
		a.Draw(screen)
	}

	for _, al := range g.alienLasers {
		al.Draw(screen)
	}

	textToDraw := fmt.Sprintf("%06d", g.score)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   24,
	}, op)

	if g.score > highScore {
		highScore = g.score
	}

	textToDraw = fmt.Sprintf("HIGH SCORE %06d", highScore)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 75)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   16,
	}, op)

	textToDraw = fmt.Sprintf("LEVEL %d", g.currentLevel)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight-40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.LevelFont,
		Size:   16,
	}, op)
}

func (g *GameScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameScene) isPlayerCollidingWithAlien() {
	for _, a := range g.aliens {
		if a.alienObj.IsIntersecting(g.player.playerObj) {
			if !a.game.player.isShielded {
				if !a.game.explosionPlayer.IsPlaying() {
					a.game.explosionPlayer.Rewind()
					a.game.explosionPlayer.Play()
				}
				a.game.player.isDying = true
			}
		}
	}
}

func (g *GameScene) isPlayerHitByAlienLaser() {
	for _, l := range g.alienLasers {
		if l.laserObj.IsIntersecting(g.player.playerObj) {
			if !g.player.isShielded {
				if !g.explosionPlayer.IsPlaying() {
					g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
				g.player.isDying = true
			}
		}
	}
}

func (g *GameScene) isAlienHitByPlayerLaser() {
	for _, a := range g.aliens {
		for _, l := range g.lasers {
			if a.alienObj.IsIntersecting(l.laserObj) {
				laserData := l.laserObj.Data().(*ObjectData)
				delete(g.alienLasers, laserData.index)
				g.space.Remove(l.laserObj)
				a.sprite = g.explosionSprite
				g.score += g.score + 50
				if !g.explosionPlayer.IsPlaying() {
					g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
			}
		}
	}
}

func (g *GameScene) letAliensAttack() {
	if len(g.aliens) > 0 {
		if !g.alienSoundPLayer.IsPlaying() {
			g.alienSoundPLayer.Rewind()
			g.alienSoundPLayer.Play()
		}

		g.alienAttackTimer.Update()

		if g.alienAttackTimer.IsReady() {
			g.alienAttackTimer.Reset()
			for _, a := range g.aliens {
				halfW, halfH := HalfOfTheImage(a.sprite)

				var degreesRadian float64
				if !a.isIntelligent {
					degreesRadian = rand.Float64() * (math.Pi * 2)
				} else {
					degreesRadian = math.Atan2(g.player.position.Y-a.position.Y, g.player.position.X-a.position.X)
					degreesRadian = degreesRadian - math.Pi*-0.5
				}

				r := degreesRadian

				offsetX := float64(a.sprite.Bounds().Dx() - int(halfW))
				offsetY := float64(a.sprite.Bounds().Dy() - int(halfH))

				spawnPos := Vector{
					X: a.position.X + halfW + math.Sin(r) - offsetX,
					Y: a.position.Y + halfH + math.Cos(r) - offsetY,
				}

				laser := NewAlienLaser(spawnPos, r)
				g.alienLaserCount++
				g.alienLasers[g.alienLaserCount] = laser
				if !g.alienLaserPlayer.IsPlaying() {
					_ = g.alienLaserPlayer.Rewind()
					g.alienLaserPlayer.Play()
				}
			}
		}

	}
}

func (g *GameScene) removeOffScreenLasers() {
	for i, l := range g.lasers {
		if l.position.X > ScreenWidth+200 || l.position.Y > ScreenHeight+200 || l.position.X < -200 || l.position.Y < -200 {
			g.space.Remove(l.laserObj)
			delete(g.lasers, i)
		}
	}

	for i, l := range g.alienLasers {
		if l.position.X > ScreenWidth+200 || l.position.Y > ScreenHeight+200 || l.position.X < -200 || l.position.Y < -200 {
			g.space.Remove(l.laserObj)
			delete(g.alienLasers, i)
		}
	}
}

func (g *GameScene) spawnAliens() {
	g.alienSpawnTimer.Update()
	if len(g.aliens) <= 3 {
		if g.alienSpawnTimer.IsReady() {
			g.alienSpawnTimer.Reset()
			rnd := rand.Intn(100-1) + 1
			if rnd > 25 {
				a := NewAlien(baseAlienVelocity, g)
				g.space.Add(a.alienObj)
				g.alienCount++
				g.aliens[g.alienCount] = a
			}
		}
	}
}

func (g *GameScene) removeOffScreenAliens() {
	for i, a := range g.aliens {
		if a.position.X > ScreenWidth+200 || a.position.Y > ScreenHeight+200 || a.position.X < -200 || a.position.Y < -200 {
			g.space.Remove(a.alienObj)
			delete(g.aliens, i)
		}
	}
}

func (g *GameScene) updateShield() {
	if g.shield != nil {
		g.shield.Update()
	}
}

func (g *GameScene) isLevelComplete(state *State) {
	if g.meteorCount >= g.meteorsForLevel && len(g.meteors) == 0 {
		g.baseVelocity = baseMeteorVelocity
		g.currentLevel++

		if g.currentLevel%5 == 0 {
			if g.player.livesRemaining < 6 {
				g.player.livesRemaining++
				x := float64(20 + (len(g.player.lifeIndicators) * 50))
				y := 20.0
				g.player.lifeIndicators = append(g.player.lifeIndicators, NewLifeIndicator(Vector{X: x, Y: y}, 0))
			}
		}

		g.beatWaitTime = baseBeatWaitTime
		state.SceneManager.GoToScene(&LevelStartScene{
			game:           g,
			nextLevelTimer: NewTimer(time.Second * 2),
			stars:          GenerateStars(numberOfStars),
		})
	}
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

			if g.score >= highScore {
				highScore = g.score
				if err := updateHighScore(highScore); err != nil {
					log.Println("Error updating high score:", err)
				}
			}

			state.SceneManager.GoToScene(&GameOverScene{
				game:        g,
				meteors:     make(map[int]*Meteor),
				meteorCount: 5,
				stars:       GenerateStars(numberOfStars),
			})
		} else {
			score := g.score
			livesRemaining := g.player.livesRemaining
			lifeSlice := g.player.lifeIndicators[:len(g.player.lifeIndicators)-1]
			stars := g.stars
			shieldsRemaining := g.player.shieldRemaining
			shieldIndicatorSlice := g.player.shieldIndicators
			g.Reset()
			g.score = score
			g.player.livesRemaining = livesRemaining
			g.player.lifeIndicators = lifeSlice
			g.stars = stars
			g.player.shieldRemaining = shieldsRemaining
			g.player.shieldIndicators = shieldIndicatorSlice
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
				g.bounceMeteor(m)
			}
		}
	}
}

func (g *GameScene) bounceMeteor(m *Meteor) {
	direction := Vector{
		X: (ScreenWidth/2 - m.position.X) * -1,
		Y: (ScreenHeight/2 - m.position.Y) * -1,
	}
	normalizeDirection := direction.Normalize()
	velocity := g.baseVelocity

	movement := Vector{
		X: normalizeDirection.X * velocity,
		Y: normalizeDirection.Y * velocity,
	}

	m.movement = movement
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

		for i, a := range g.aliens {
			if a.sprite == g.explosionSprite {
				delete(g.aliens, i)
				g.space.Remove(a.alienObj)
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
	g.stars = GenerateStars(numberOfStars)
	g.player.shieldRemaining = numberOfShields
	g.player.isShielded = false
	g.aliens = make(map[int]*Alien)
	g.alienLasers = make(map[int]*AlienLaser)
	g.alienCount = 0
	g.alienLaserCount = 0
}
