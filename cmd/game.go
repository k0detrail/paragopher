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
		if err == config.ErrQuit {
			if g.Score > gameData.HiScore {
				fmt.Println("Updating HiScore...")
				gameData.HiScore = g.Score
				utils.SaveData(gameData)
			}
			os.Exit(0)
		}
		if err == config.ErrGameOver {
		}
		log.Fatal(err)
	}
}
