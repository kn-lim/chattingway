package gamble

import (
	"math/rand"
	"time"
)

// CoinFlip returns either "Heads" or "Tails"
func CoinFlip() string {
	// Create a new rand.Rand instance with a seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number (0 or 1)
	if r.Intn(2) == 0 {
		return "Heads"
	} else {
		return "Tails"
	}
}
