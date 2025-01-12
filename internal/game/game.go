package game

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
)

type Game struct {
	Score          int
	hiScore        int
	gameOver       bool
	showExitDialog bool

	barrelAngle float64
	barrelImage *ebiten.Image

	bullets  []*Bullet
	lastShot time.Time

	helicopters  []*Helicopter
	paratroopers []*Paratrooper
}

func NewGame(hiScore int) *Game {
	game := &Game{
		bullets:  make([]*Bullet, 0),
		lastShot: time.Now(),
		hiScore:  hiScore,
	}
	width := config.BaseWidth
	game.barrelImage = ebiten.NewImage(int(width), int(width))
	game.barrelImage.Fill(config.TransparentBlack)

	rectX := width/2 - width/12
	rectY := width / 12
	rectW := width / 6
	rectH := width / 3
	vector.DrawFilledRect(
		game.barrelImage,
		rectX,
		rectY,
		rectW,
		rectH,
		config.ColourTeal,
		true,
	)

	circleX := width / 2
	circleY := width / 2
	pinkCircleRadius := width / 6
	tealCircleRaduis := width / 24
	vector.DrawFilledCircle(
		game.barrelImage,
		circleX,
		circleY,
		pinkCircleRadius,
		config.ColourPink,
		true,
	)
	vector.DrawFilledCircle(
		game.barrelImage,
		circleX,
		circleY,
		tealCircleRaduis,
		config.ColourTeal,
		true,
	)

	topCircleX, topCircleY := width/2, width/12
	topCircleRadius := width / 12
	vector.DrawFilledCircle(
		game.barrelImage,
		topCircleX,
		topCircleY,
		topCircleRadius,
		config.ColourTeal,
		true,
	)

	return game
}

// Ebiten Game Interface
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawTurret(screen)
	g.drawBullets(screen)
	g.drawHelicopters(screen)
	g.drawParatroopers(screen)

	// Display Score
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("SCORE: %d    HI-SCORE: %d", g.Score, g.hiScore),
	)

	if g.showExitDialog {
		overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		overlay.Fill(config.SemiTransparentBlack)
		screen.DrawImage(overlay, nil)

		dialogWidth, dialogHeight := 300, 150
		dialogX := (screen.Bounds().Dx() - dialogWidth) / 2
		dialogY := (screen.Bounds().Dy() - dialogHeight) / 2
		dialog := ebiten.NewImage(dialogWidth, dialogHeight)
		dialog.Fill(config.ColourDarkGrey)

		vector.DrawFilledRect(
			dialog,
			0,
			0,
			float32(dialogWidth),
			5,
			config.ColourBlack,
			false,
		)
		vector.DrawFilledRect(
			dialog,
			0,
			float32(dialogHeight-5),
			float32(dialogWidth),
			5,
			config.ColourBlack,
			false,
		)
		vector.DrawFilledRect(
			dialog,
			0,
			0,
			5,
			float32(dialogHeight),
			config.ColourBlack,
			false,
		)
		vector.DrawFilledRect(
			dialog,
			float32(dialogWidth-5),
			0,
			5,
			float32(dialogHeight),
			config.ColourBlack,
			false,
		)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(dialogX), float64(dialogY))
		screen.DrawImage(dialog, op)

		msg := "Do you want to exit the game?"
		textX := dialogX + 50
		textY := dialogY + 40
		ebitenutil.DebugPrintAt(screen, msg, textX, textY)

		yesText := "Y: Yes"
		noText := "N: No"
		ebitenutil.DebugPrintAt(screen, yesText, dialogX+50, dialogY+90)
		ebitenutil.DebugPrintAt(screen, noText, dialogX+200, dialogY+90)
	}
}

func (g *Game) Update() error {
	if !g.showExitDialog {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.showExitDialog = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			if g.barrelAngle > config.BarrelAngleMin {
				g.barrelAngle--
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			if g.barrelAngle < config.BarrelAngleMax {
				g.barrelAngle++
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			if time.Since(g.lastShot).Milliseconds() > config.ShotCooldown {
				g.shoot()
			}
		}

		g.updateBullets()
		g.spawnHelicopter()
		g.updateHelicopters()
		g.updateParatroopers()
		g.checkHits()
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			return config.ErrQuit
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			g.showExitDialog = false
		}
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
