package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/game"
)

func main() {
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("ParaGophers")

	g := game.NewGame(1000)

	if err := ebiten.RunGame(g); err != nil {
		if err == config.ErrEscPressed {
			os.Exit(0)
		}
		log.Fatal(err)
	}
}
