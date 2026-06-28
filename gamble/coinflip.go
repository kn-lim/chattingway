// Package gamble provides chance-based games such as coin flips and dice rolls.
package gamble

import (
	"math/rand"
	"time"
)

// CoinFlip returns "Heads" or "Tails" with equal probability.
func CoinFlip() string {
	// Create a new rand.Rand instance with a seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number (0 or 1)
	if r.Intn(2) == 0 {
		return "Heads"
	}

	return "Tails"
}
