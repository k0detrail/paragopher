package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ystepanoff/paragophers/internal/config"
	"github.com/ystepanoff/paragophers/internal/game"
)

func main() {
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("ParaGophers")

	g := &game.Game{}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
