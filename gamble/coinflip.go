// Package gamble provides chance-based games such as coin flips and dice rolls.
package gamble

import (
	"math/rand/v2"
)

// CoinFlip returns "Heads" or "Tails" with equal probability.
func CoinFlip() string {
	// Generate a random number (0 or 1)
	if rand.IntN(2) == 0 {
		return "Heads"
	}

	return "Tails"
}
