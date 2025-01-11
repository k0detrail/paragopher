package config

import (
	"errors"
	"image/color"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600

	BarrelAngleMin = -90.0
	BarrelAngleMax = 90.0

	HelicopterSpawnChance = 0.009
	HelicopterSpeed       = 1.5
	HelicopterDropRate    = 2
	HelicopterBodyW       = 30.0
	HelicopterBodyH       = 10.0
	HelicopterTailW       = 16.0
	HelicopterTailH       = 3.0
	HelicopterRotorLen    = 20.0

	ParatrooperFallSpeed = 1.2
	ParatrooperWalkSpeed = 0.5
	ParatrooperWidth     = 6.0
	ParatrooperHeight    = 10.0
	ParachuteRadius      = 20.0

	BulletSpeed  = 3.0
	BulletRadius = 2.0
	ShotCooldown = 200

	GroundY = 580
)

var (
	BaseW = float32(ScreenWidth) / 10.0
	BaseH = float32(ScreenHeight) / 10.0

	ColourTeal       = color.RGBA{R: 101, G: 247, B: 246, A: 255}
	ColourPink       = color.RGBA{R: 255, G: 82, B: 242, A: 255}
	ColourMagenta    = color.RGBA{255, 0, 255, 255}
	ColourWhite      = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColourBlack      = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	TransparentBlack = color.RGBA{R: 0, G: 0, B: 0, A: 0}

	ErrEscPressed = errors.New("ESC pressed")
)
