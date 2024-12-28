package gamble_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/chattingway/gamble"
)

func TestCoinflip(t *testing.T) {
	result := gamble.CoinFlip()
	assert.Contains(t, []string{"Heads", "Tails"}, result)
}
