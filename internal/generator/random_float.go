package generator

import "math/rand"

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
