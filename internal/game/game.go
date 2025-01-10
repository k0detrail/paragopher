package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
)

type Game struct {
	barrelAngle float64
	barrelImage *ebiten.Image

	bullets  []*Bullet
	lastShot time.Time

	score    int
	hiScore  int
	gameOver bool
}

func NewGame(hiScore int) *Game {
	game := &Game{
		bullets:  make([]*Bullet, 0),
		lastShot: time.Now(),
		hiScore:  hiScore,
	}
	width := config.BaseW
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
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return config.ErrEscPressed
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

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
