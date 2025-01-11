package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/game"
	"github.com/ystepanoff/paragopher/internal/utils"
)

func main() {
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("ParaGophers")

	gameData, err := utils.LoadData()
	if err != nil {
		fmt.Println("Error loading game data!")
		gameData = &utils.GameData{}
	}

	g := game.NewGame(gameData.HiScore)

	if err := ebiten.RunGame(g); err != nil {
		if err == config.ErrEscPressed {
			os.Exit(0)
		}
		log.Fatal(err)
	}
}
