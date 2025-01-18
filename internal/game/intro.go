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
	introText     = "P A R A G O P H E R"
	introSkipText = "Press ENTER to skip intro..."
	scaleFactor   = 4
)

var textFaceSource *text.GoTextFaceSource = nil

var colourLayers = []color.Color{
	config.ColourDarkGrey,
	config.ColourPink,
	config.ColourTeal,
}

func (g *Game) initIntro() {
	audio.Play(g.soundProfile.IntroPlayer)
	var err error
	textFaceSource, err = text.NewGoTextFaceSource(
		bytes.NewReader(fonts.PressStart2P_ttf),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) drawIntro(screen *ebiten.Image) {
	fontFace := &text.GoTextFace{
		Source: textFaceSource,
		Size:   32,
	}
	textW, textH := text.Measure(introText, fontFace, 1.0)
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

	fontFace.Size = 12
	textSkipW, _ := text.Measure(introSkipText, fontFace, 1.0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(
		(config.ScreenWidth-textSkipW)/2.0,
		config.ScreenHeight-textH,
	)
	op.ColorScale.ScaleWithColor(config.ColourWhite)
	text.Draw(screen, introSkipText, fontFace, op)

	if g.introStep < len(introText) &&
		time.Since(g.lastIntroStep).Milliseconds() > 300 {
		g.introStep++
		g.lastIntroStep = time.Now()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) || g.isIntroFinished() {
		g.soundProfile.IntroPlayer.Close()
		g.showIntro = false
	}
}

func (g *Game) isIntroFinished() bool {
	return g.introStep == len(introText) &&
		time.Since(g.lastIntroStep).Seconds() > 2
}
