package goasteroids

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity  = 0.25
	meteorSpawnTime     = 100 * time.Millisecond
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 1000 * time.Millisecond
)

type GameScene struct {
	player            *Player
	baseVelocity      float64
	meteorCount       int
	meteorsSpawnTimer *Timer
	meteors           map[int]*Meteor
	meteorsForLevel   int
	velocityTimer     *Timer
	space             *resolv.Space
	lasers            map[int]*Laser
	laserCount        int
}

func NewGameScene() *GameScene {
	g := &GameScene{
		meteorsSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:      baseMeteorVelocity,
		velocityTimer:     NewTimer(meteorSpeedUpTime),
		meteors:           make(map[int]*Meteor),
		meteorsForLevel:   2,
		meteorCount:       0,
		space:             resolv.NewSpace(ScreenWidth, ScreenHeight, 16, 16),
		lasers:            make(map[int]*Laser),
		laserCount:        0,
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)
	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()
	g.spawnMeteors()

	for _, m := range g.meteors {
		m.Update()
	}

	for _, l := range g.lasers {
		l.Update()
	}

	g.speedUpMeteors()

	g.isPlayerCollidingWithMeteor()

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

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
			data := m.meteorObj.Data().(*ObjectData)
			fmt.Println("Player collided with meteor", data.index)
		}
	}
}
