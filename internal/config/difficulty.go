package config

type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyHard
	DifficultyVeteran
)

var CurrentDifficulty = DifficultyEasy

func HelicopterSpeedByDifficulty() float64 {
	switch CurrentDifficulty {
	case DifficultyEasy:
		return 0.4
	case DifficultyHard:
		return 0.6
	case DifficultyVeteran:
		return 0.8
	}
	return HelicopterSpeed
}

func ParatrooperFallSpeedByDifficulty() float64 {
	switch CurrentDifficulty {
	case DifficultyEasy:
		return 0.2
	case DifficultyHard:
		return 0.35
	case DifficultyVeteran:
		return 0.6
	}
	return ParatrooperFallSpeed
}

func ParatrooperSpawnChanceByDifficulty() float64 {
	switch CurrentDifficulty {
	case DifficultyEasy:
		return 0.005
	case DifficultyHard:
		return 0.008
	case DifficultyVeteran:
		return 0.012
	}
	return ParatrooperSpawnChance
}

func HelicopterSpawnChanceByDifficulty() float64 {
	switch CurrentDifficulty {
	case DifficultyEasy:
		return 0.002
	case DifficultyHard:
		return 0.004
	case DifficultyVeteran:
		return 0.008
	}
	return HelicopterSpawnChance
}
