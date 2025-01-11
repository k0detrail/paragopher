package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
)

type Paratrooper struct {
	x, y      float32
	vy        float32
	parachute bool
	landed    bool
	onBase    bool
	climbing  bool
	onTopOf   *Paratrooper
}

func (g *Game) spawnParatrooper(x, y float32) {
	g.paratroopers = append(g.paratroopers, &Paratrooper{
		x:         x,
		y:         y,
		vy:        config.ParatrooperFallSpeed,
		parachute: true,
		landed:    false,
	})
}

func (g *Game) drawParatrooper(screen *ebiten.Image, p *Paratrooper) {
	if !p.landed && p.parachute {
		vector.DrawFilledCircle(
			screen,
			p.x,
			p.y-config.ParachuteRadius,
			config.ParachuteRadius,
			config.ColourTeal,
			false,
		)
		vector.DrawFilledRect(
			screen,
			p.x-config.ParachuteRadius,
			p.y-config.ParachuteRadius,
			2*config.ParachuteRadius,
			config.ParachuteRadius,
			config.ColourBlack,
			false,
		)
		vector.StrokeLine(
			screen,
			p.x-config.ParachuteRadius,
			p.y,
			p.x-config.ParachuteRadius/2.0,
			config.ParachuteRadius,
			1,
			config.ColourTeal,
			false,
		)
		vector.StrokeLine(
			screen,
			p.x+config.ParachuteRadius,
			p.y,
			p.x-config.ParachuteRadius/2.0,
			config.ParachuteRadius,
			1,
			config.ColourTeal,
			false,
		)
		vector.DrawFilledRect(
			screen,
			p.x-config.ParatrooperWidth/2.0,
			p.y,
			config.ParatrooperWidth,
			config.ParatrooperHeight,
			config.ColourTeal,
			false,
		)
	} else {
		vector.DrawFilledRect(
			screen,
			p.x-config.ParatrooperWidth/2.0,
			p.y,
			config.ParatrooperWidth,
			config.ParatrooperHeight,
			config.ColourTeal,
			false,
		)
	}
}
