package gamble_test

import (
	"testing"

	"github.com/kn-lim/chattingway/gamble"
	"github.com/stretchr/testify/assert"
)

func TestRoll(t *testing.T) {
	testInput := "2d6+3"

	resultString, resultInt, err := gamble.Roll(testInput)

	assert.Nil(t, err)
	assert.NotEmpty(t, resultString)
	assert.NotZero(t, resultInt)
}
