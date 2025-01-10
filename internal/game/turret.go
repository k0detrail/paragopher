package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragophers/internal/config"
)

type Bullet struct {
	x, y   float32
	vx, vy float32
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
