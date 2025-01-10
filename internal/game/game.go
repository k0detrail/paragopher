package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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

func (g *Game) drawTurret(screen *ebiten.Image) {
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

	pinkBaseX := (float32(config.ScreenWidth) - config.BaseW/3.0) / 2.0
	pinkBaseY := float32(config.ScreenHeight)
	pinkBaseY -= config.BaseH
	pinkBaseY -= config.BaseW / 3
	pinkBaseW := config.BaseW / 3
	pinkBaseH := config.BaseW / 3

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
	centerX := float64(config.ScreenWidth) / 2.0
	centerY := float64(config.ScreenHeight)
	centerY -= float64(config.BaseH)
	centerY -= float64(config.BaseW) / 3.0
	centerY -= float64(config.BaseW) / 24.0
	barrelW := float64(g.barrelImage.Bounds().Dx())
	barrelH := float64(g.barrelImage.Bounds().Dy())
	op.GeoM.Translate(-barrelW/2.0, -barrelH/2.0)
	op.GeoM.Rotate(g.barrelAngle * math.Pi / 180)
	op.GeoM.Translate(centerX, centerY)
	screen.DrawImage(g.barrelImage, op)
}

func (g *Game) drawBullets(screen *ebiten.Image) {
	for _, b := range g.bullets {
		vector.DrawFilledCircle(
			screen,
			b.x,
			b.y,
			config.BulletRadius,
			config.ColourWhite,
			false,
		)
	}
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
		g.shoot()
	}

	g.updateBullets()

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) shoot() {
	barrelCircleX := float32(config.ScreenWidth) / 2.0
	barrelCircleY := float32(config.ScreenHeight)
	barrelCircleY -= config.BaseH
	barrelCircleY -= config.BaseW / 3.0
	barrelCircleY -= config.BaseW / 24.0

	width := config.BaseW
	localTipX := width / 2
	localTipY := width / 12
	angleRadians := float64(g.barrelAngle * math.Pi / 180.0)
	dx := float64(localTipX - width/2)
	dy := float64(localTipY - width/2)
	rx := float32(dx*math.Cos(angleRadians) - dy*math.Sin(angleRadians))
	ry := float32(dx*math.Sin(angleRadians) + dy*math.Cos(angleRadians))
	tipX := barrelCircleX + rx
	tipY := barrelCircleY + ry
	realAngleRadians := (90.0 - g.barrelAngle) * math.Pi / 180.0
	vx := float32(config.BulletSpeed * math.Cos(realAngleRadians))
	vy := -float32(config.BulletSpeed * math.Sin(realAngleRadians))
	g.bullets = append(g.bullets, &Bullet{
		x:  tipX,
		y:  tipY,
		vx: vx,
		vy: vy,
	})
}

func (g *Game) updateBullets() {
	active := make([]*Bullet, 0)
	for _, b := range g.bullets {
		b.x += b.vx
		b.y += b.vy
		if b.x < 0 || b.x > config.ScreenWidth || b.y < 0 ||
			b.y > config.ScreenHeight {
			continue
		}
		active = append(active, b)
	}
	g.bullets = active
}
