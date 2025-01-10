package game

type Paratrooper struct {
	x, y      float32
	vy        float32
	parachute bool
	landed    bool
	onBase    bool
	climbing  bool
	onTopOf   *Paratrooper
}
