package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragophers/internal/config"
)

type Game struct {
	barrelAngle float64
	barrelImage *ebiten.Image

	bullets []*Bullet

	score    int
	hiScore  int
	gameOver bool
}

type Bullet struct {
	x, y   float32
	vx, vy float32
}

type Helicopter struct {
	x, y     float32
	vx       float32
	lastDrop time.Time
}

type Paratrooper struct {
	x, y      float32
	vy        float32
	parachute bool
	landed    bool
	onBase    bool
	climbing  bool
	onTopOf   *Paratrooper
}

func NewGame(hiScore int) *Game {
	game := &Game{
		bullets: make([]*Bullet, 0),
		hiScore: hiScore,
	}
	width := config.BaseW
	game.barrelImage = ebiten.NewImage(int(width), int(width))
	game.barrelImage.Fill(config.TransparentBlack)

	xRect, yRect := width/2-width/12, width/12
	wRect, hRect := width/6, width/3
	vector.DrawFilledRect(
		game.barrelImage,
		xRect,
		yRect,
		wRect,
		hRect,
		config.ColourTeal,
		true,
	)

	circleX, circleY := width/2, width/2
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
	screen.Fill(config.ColourBlack)
	baseX := (config.ScreenWidth - config.BaseW) / 2
	baseY := config.ScreenHeight - config.BaseH
	vector.DrawFilledRect(
		screen,
		baseX,
		baseY,
		config.BaseW,
		config.BaseH,
		config.ColourWhite,
		false,
	)

	pinkBaseX := (float32(config.ScreenWidth) - float32(config.BaseW)/3) / 2
	pinkBaseY := float32(config.ScreenHeight)
	pinkBaseY -= float32(config.BaseH)
	pinkBaseY -= float32(config.BaseW) / 3
	pinkBaseW := float32(config.BaseW) / 3
	pinkBaseH := float32(config.BaseW) / 3

	vector.DrawFilledRect(
		screen,
		pinkBaseX,
		pinkBaseY,
		pinkBaseW,
		pinkBaseH,
		config.ColourPink,
		false,
	)

	op := &ebiten.DrawImageOptions{}
	centerX := float64(config.ScreenWidth / 2)
	centerY := float64(config.ScreenHeight)
	centerY -= float64(config.BaseH)
	centerY -= float64(config.BaseW) / 3
	centerY -= float64(config.BaseW) / 24
	barrelW := float64(g.barrelImage.Bounds().Dx())
	barrelH := float64(g.barrelImage.Bounds().Dy())
	op.GeoM.Translate(-float64(barrelW)/2, -float64(barrelH)/2)
	op.GeoM.Rotate(g.barrelAngle * math.Pi / 180)
	op.GeoM.Translate(centerX, centerY)
	screen.DrawImage(g.barrelImage, op)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
