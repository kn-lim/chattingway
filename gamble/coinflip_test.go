package gamble_test

import (
	"testing"

	"github.com/kn-lim/chattingway/v2/gamble"
	"github.com/stretchr/testify/assert"
)

func TestCoinflip(t *testing.T) {
	result := gamble.CoinFlip()
	assert.Contains(t, []string{"Heads", "Tails"}, result)
}
