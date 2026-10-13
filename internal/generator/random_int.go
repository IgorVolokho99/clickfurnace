package generator

import "math/rand"

func randomUint64(min, max uint64) uint64 {
	if min == 0 && max == ^uint64(0) {
		return rand.Uint64()
	}

	n := max - min + 1

	threshold := -n % n

	for {
		value := rand.Uint64()

		if value >= threshold {
			return min + value%n
		}
	}
}
