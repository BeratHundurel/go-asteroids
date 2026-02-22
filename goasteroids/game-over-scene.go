package goasteroids

import (
	"go-asteroids/assets"
	"image/color"
	"os"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameOverScene struct {
	game        *GameScene
	meteors     map[int]*Meteor
	meteorCount int
	stars       []*Star
}

func (o *GameOverScene) Draw(screen *ebiten.Image) {
	for _, s := range o.stars {
		s.Draw(screen)
	}

	meteorKeys := make([]int, 0, len(o.meteors))
	for k := range o.meteors {
		meteorKeys = append(meteorKeys, k)
	}
	sort.Ints(meteorKeys)
	for _, k := range meteorKeys {
		o.meteors[k].Draw(screen)
	}

	textToDraw := "Game Over Press Space to Restart"
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight/2+100)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   48,
	}, op)

	if o.game.score > originalHighScore {
		textToDraw = "New High Score!"
		op = &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignCenter,
			},
		}
		op.ColorScale.ScaleWithColor(color.White)
		op.GeoM.Translate(ScreenWidth/2, ScreenHeight/2-200)
		text.Draw(screen, textToDraw, &text.GoTextFace{
			Source: assets.TitleFont,
			Size:   48,
		}, op)
	}
}

func (o *GameOverScene) Update(state *State) error {
	if len(o.meteors) < 10 {
		m := NewMeteor(0.25, &GameScene{}, len(o.meteors)-1)
		o.meteorCount++
		o.meteors[o.meteorCount] = m
	}

	for _, m := range o.meteors {
		m.Update()
	}

	o.separateOverlappingMeteors()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		o.game.Reset()
		state.SceneManager.GoToScene(o.game)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	return nil
}

func (o *GameOverScene) separateOverlappingMeteors() {
	list := make([]*Meteor, 0, len(o.meteors))
	for _, m := range o.meteors {
		list = append(list, m)
	}
	separateMeteors(list)
}
