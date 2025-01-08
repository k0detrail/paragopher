package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ystepanoff/paragophers/internal/config"
	"github.com/ystepanoff/paragophers/internal/config/game"
)

func main() {
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("ParaGophers")

	g := &game.Game{}
}
