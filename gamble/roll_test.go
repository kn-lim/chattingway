package gamble_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/chattingway/gamble"
)

func TestRoll(t *testing.T) {
	testInput := "2d6+3"

	resultString, resultInt, err := gamble.Roll(testInput)

	assert.Nil(t, err)
	assert.NotEmpty(t, resultString)
	assert.NotZero(t, resultInt)
}
