package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragophers/internal/config"
)

type Game struct {
	barrelAngle float64
	barrelImage ebiten.Image

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

// Ebiten Game Interface
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(config.ColourBlack)
	baseX := (config.ScreenWidth - config.BaseW) / 2
	baseY := config.ScreenHeight - config.BaseH
	vector.DrawFilledRect(
		screen,
		baseX,
		baseY,
		config.BaseH,
		config.BaseW,
		config.ColourWhite,
		true,
	)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
