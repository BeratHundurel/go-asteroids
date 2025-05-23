package goasteroids

import "github.com/hajimehoshi/ebiten/v2"

func HalfOfTheImage(image *ebiten.Image) (float64, float64) {
	bounds := image.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	return halfW, halfH
}