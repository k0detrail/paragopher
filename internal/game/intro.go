package game

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/ystepanoff/paragopher/internal/config"
	"golang.org/x/image/font/basicfont"
)

const (
	introText   = "P A R A G O P H E R"
	scaleFactor = 4
)

func (g *Game) drawIntro(screen *ebiten.Image) {
	message := ""
	textWidth := 200
	textHeight := 100
	textImg := ebiten.NewImage(textWidth, textHeight)
	textImg.Fill(color.Transparent)

	message = introText[:g.introStep]

	face := basicfont.Face7x13

	text.Draw(textImg, message, face, 33, 51, config.ColourDarkGrey)
	text.Draw(textImg, message, face, 34, 50, config.ColourPink)
	text.Draw(textImg, message, face, 35, 50, config.ColourTeal)

	screen.Fill(config.ColourBlack) // Clear screen with black

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleFactor, scaleFactor) // Scale the text
	op.GeoM.Translate(
		(float64(config.ScreenWidth)-(float64(textImg.Bounds().Dx())*scaleFactor))/2,
		(float64(config.ScreenHeight)-(float64(textImg.Bounds().Dy())*scaleFactor))/4,
	)

	screen.DrawImage(textImg, op)

	if time.Since(g.lastIntroStep).Milliseconds() > 300 {
		g.introStep++
		g.lastIntroStep = time.Now()
	}

	if g.introStep == len(introText)+1 {
		g.showIntro = false
	}
}
