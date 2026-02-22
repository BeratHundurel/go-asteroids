package goasteroids

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Star struct {
	x          float32
	y          float32
	r          float32
	brightness float32
}

func NewStar() *Star {
	return &Star{
		x:          rand.Float32() * ScreenWidth,
		y:          rand.Float32() * ScreenHeight,
		r:          rand.Float32() * (3 - 1),
		brightness: rand.Float32() * 0xff,
	}
}

func (s *Star) Draw(screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff,
	}
	vector.FillCircle(screen, s.x, s.y, s.r, c, false)
}

func (s *Star) Update() {}

func GenerateStars(n int) []*Star {
	stars := make([]*Star, n)
	for i := range n {
		stars[i] = NewStar()
	}
	return stars
}

// RenderStarField pre-renders n randomly generated stars onto a single
// ScreenWidth x ScreenHeight image.
func RenderStarField(n int) *ebiten.Image {
	img := ebiten.NewImage(ScreenWidth, ScreenHeight)
	for range n {
		NewStar().Draw(img)
	}
	return img
}
