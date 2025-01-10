package game

import "time"

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
