package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/audio"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/utils"
)

type Bullet struct {
	x, y   float32
	vx, vy float32
}

func (g *Game) initTurretImage() {
	w := int(config.BaseWidth)
	h := int(config.BaseWidth/3.0 + config.BaseHeight + 1.0)
	g.turretBaseImage = ebiten.NewImage(w, h)
	g.turretBaseImage.Fill(config.TransparentBlack)
	vector.DrawFilledRect(
		g.turretBaseImage,
		0.0,
		float32(h)-config.BaseHeight,
		config.BaseWidth,
		config.BaseHeight,
		config.ColourWhite,
		true,
	)
	vector.DrawFilledRect(
		g.turretBaseImage,
		config.BaseWidth/3.0,
		0.0,
		config.BaseWidth/3.0,
		config.BaseWidth/3.0,
		config.ColourPink,
		false,
	)
}

func (g *Game) initBarrelImage() {
	w := config.BaseWidth
	g.barrelImage = ebiten.NewImage(int(w), int(w))
	g.barrelImage.Fill(config.TransparentBlack)

	rectX := w/2 - w/12
	rectY := w / 12
	rectW := w / 6
	rectH := w / 3
	vector.DrawFilledRect(
		g.barrelImage,
		rectX,
		rectY,
		rectW,
		rectH,
		config.ColourTeal,
		false,
	)

	circleX := w / 2
	circleY := w / 2
	pinkCircleRadius := w / 6
	tealCircleRaduis := w / 24
	vector.DrawFilledCircle(
		g.barrelImage,
		circleX,
		circleY,
		pinkCircleRadius,
		config.ColourPink,
		true,
	)
	vector.DrawFilledCircle(
		g.barrelImage,
		circleX,
		circleY,
		tealCircleRaduis,
		config.ColourTeal,
		true,
	)
	topCircleX, topCircleY := w/2, w/12
	topCircleRadius := w / 12
	vector.DrawFilledCircle(
		g.barrelImage,
		topCircleX,
		topCircleY,
		topCircleRadius,
		config.ColourTeal,
		true,
	)
}

func (g *Game) initBulletImage() {
	w := int(2 * config.BulletRadius)
	g.bulletImage = ebiten.NewImage(w, w)
	vector.DrawFilledCircle(
		g.bulletImage,
		config.BulletRadius,
		config.BulletRadius,
		config.BulletRadius,
		config.ColourWhite,
		true,
	)
}

func (g *Game) drawTurret(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(config.ScreenWidth-config.BaseWidth)/2.0,
		float64(config.ScreenHeight-g.turretBaseImage.Bounds().Dy()),
	)
	screen.DrawImage(g.turretBaseImage, op)

	op = &ebiten.DrawImageOptions{}
	centerX := float64(config.ScreenWidth) / 2.0
	centerY := float64(config.ScreenHeight)
	centerY -= float64(config.BaseHeight)
	centerY -= float64(config.BaseWidth) / 3.0
	centerY -= float64(config.BaseWidth) / 24.0
	barrelW := float64(g.barrelImage.Bounds().Dx())
	barrelH := float64(g.barrelImage.Bounds().Dy())
	op.GeoM.Translate(-barrelW/2.0, -barrelH/2.0)
	op.GeoM.Rotate(g.barrelAngle * math.Pi / 180)
	op.GeoM.Translate(centerX, centerY)
	screen.DrawImage(g.barrelImage, op)
}

func (g *Game) drawBullets(screen *ebiten.Image) {
	for _, b := range g.bullets {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			float64(b.x-config.BulletRadius/2.0),
			float64(b.y-config.BulletRadius/2.0),
		)
		screen.DrawImage(g.bulletImage, op)
	}
}

func (g *Game) shoot() {
	audio.Play(g.soundProfile.ShootPlayer)

	barrelCircleX := float32(config.ScreenWidth) / 2.0
	barrelCircleY := float32(config.ScreenHeight)
	barrelCircleY -= config.BaseHeight
	barrelCircleY -= config.BaseWidth / 3.0
	barrelCircleY -= config.BaseWidth / 24.0

	width := config.BaseWidth
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
	g.Score = max(g.Score-1, 0)
	g.lastShot = time.Now()
}

func (g *Game) updateBullets() {
	active := make([]*Bullet, 0, len(g.bullets))
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

func (g *Game) checkHits() {
	activeBullets := make([]*Bullet, 0, len(g.bullets))

bulletLoop:
	for _, b := range g.bullets {
		for i, h := range g.helicopters {
			if utils.Overlap2D(
				b.x-config.BulletRadius/2.0,
				b.y-config.BulletRadius/2.0,
				config.BulletRadius,
				config.BulletRadius,
				h.x-float32(g.helicopterImage.Bounds().Dx())/2.0,
				h.y-float32(g.helicopterImage.Bounds().Dy())/2.0,
				float32(g.helicopterImage.Bounds().Dx()),
				float32(g.helicopterImage.Bounds().Dy()),
			) {
				audio.Play(g.soundProfile.HitPlayer)
				g.helicopters = append(g.helicopters[:i], g.helicopters[i+1:]...)
				g.Score += 10
				continue bulletLoop
			}
		}
		for i, p := range g.paratroopers {
			if utils.Overlap2D(
				b.x-config.BulletRadius/2.0,
				b.y-config.BulletRadius/2.0,
				config.BulletRadius,
				config.BulletRadius,
				p.x-config.ParatrooperWidth/2.0,
				p.y,
				config.ParatrooperWidth,
				config.ParatrooperHeight,
			) || utils.Overlap2D(
				b.x-config.BulletRadius/2.0,
				b.y-config.BulletRadius/2.0,
				config.BulletRadius,
				config.BulletRadius,
				p.x-config.ParachuteRadius,
				p.y-config.ParachuteRadius*2.0,
				config.ParachuteRadius*2.0,
				config.ParachuteRadius*2.0,
			) {
				audio.Play(g.soundProfile.HitPlayer)
				g.paratroopers = append(g.paratroopers[:i], g.paratroopers[i+1:]...)
				g.Score += 5
				continue bulletLoop
			}
		}
		activeBullets = append(activeBullets, b)
	}
	g.bullets = activeBullets
	if g.Score > g.gameData.HiScore {
		g.gameData.HiScore = g.Score
	}
}
