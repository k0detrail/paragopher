package config

import (
	"image/color"
	"time"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600

	BarrelAngleMin = -90.0
	BarrelAngleMax = 90.0
)

var (
	BaseW = ScreenWidth / 10
	BaseH = ScreenHeight / 10

	ColourTeal  = color.RGBA{R: 101, G: 247, B: 246, A: 255}
	ColourPink  = color.RGBA{R: 255, G: 82, B: 242, A: 255}
	ColourWhite = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColourBlack = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	HelicopterSpawnChance = 0.01
	HelicopterSpeed       = 1.5
	HelicopterDropRate    = 2 * time.Second

	ParatrooperFallSpeed = 1.2
	ParatrooperWalkSpeed = 0.5
	ParatrooperWidth     = 6.0
	ParatrooperHeight    = 10.0
	ParaschuteRadius     = 20.0

	TurretRadius = 20.0
	BulletSpeed  = 3.0
	BulletRadius = 2.0

	GroundY = 580
)
