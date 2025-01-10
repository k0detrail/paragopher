package config

import (
	"errors"
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
	BaseW = float32(ScreenWidth) / 10.0
	BaseH = float32(ScreenHeight) / 10.0

	ColourTeal       = color.RGBA{R: 101, G: 247, B: 246, A: 255}
	ColourPink       = color.RGBA{R: 255, G: 82, B: 242, A: 255}
	ColourWhite      = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColourBlack      = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	TransparentBlack = color.RGBA{R: 0, G: 0, B: 0, A: 0}

	HelicopterSpawnChance = 0.01
	HelicopterSpeed       = 1.5
	HelicopterDropRate    = 2 * time.Second

	ParatrooperFallSpeed = 1.2
	ParatrooperWalkSpeed = 0.5
	ParatrooperWidth     = 6.0
	ParatrooperHeight    = 10.0
	ParachuteRadius      = 20.0

	BulletSpeed  float64 = 3.0
	BulletRadius float32 = 2.0
	ShotCooldown int64   = 500

	GroundY = 580

	ErrEscPressed = errors.New("ESC pressed")
)
