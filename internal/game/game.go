package game

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	barrelAngle float64
	barrelImage ebiten.Image

	bullets []*Bullet

	score    int
	hiScore  int
	gameOver bool
}

type Bullet struct {
	x, y   float64
	vx, vy float64
}
