package main

import (
	"go-asteroids/goasteroids"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Go Asteroids")
	ebiten.SetWindowSize(goasteroids.ScreenWidth, goasteroids.ScreenHeight)

	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(&goasteroids.Game{}); err != nil {
		log.Fatal(err)
	}
}
