package game

import (
	"bytes"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/ystepanoff/paragopher/internal/audio"
	"github.com/ystepanoff/paragopher/internal/config"
)

const (
	introText   = "P A R A G O P H E R"
	scaleFactor = 4
)

var (
	fontFace     *text.GoTextFace = nil
	textW, textH float64
)

var colourLayers = []color.Color{
	config.ColourDarkGrey,
	config.ColourPink,
	config.ColourTeal,
}

func (g *Game) initIntro() {
	audio.Play(g.soundProfile.IntroPlayer)

	faceSource, err := text.NewGoTextFaceSource(
		bytes.NewReader(fonts.PressStart2P_ttf),
	)
	if err != nil {
		log.Fatal(err)
	}
	fontFace = &text.GoTextFace{
		Source: faceSource,
		Size:   32,
	}
	textW, textH = text.Measure(introText, fontFace, 1.0)
}

func (g *Game) drawIntro(screen *ebiten.Image) {
	message := introText[:g.introStep]

	for i, colour := range colourLayers {
		op := &text.DrawOptions{}
		op.GeoM.Translate(
			(config.ScreenWidth-textW)/2.0+float64((i-1)*5),
			(config.ScreenHeight-textH)/2.0,
		)
		op.ColorScale.ScaleWithColor(colour)
		text.Draw(screen, message, fontFace, op)
	}

	if g.introStep < len(introText) &&
		time.Since(g.lastIntroStep).Milliseconds() > 300 {
		g.introStep++
		g.lastIntroStep = time.Now()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) || g.introIsDone() {
		g.showIntro = false
	}
}

func (g *Game) introIsDone() bool {
	return g.introStep == len(introText) &&
		time.Since(g.lastIntroStep).Seconds() > 2
}
