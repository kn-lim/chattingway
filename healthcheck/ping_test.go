package healthcheck_test

import (
	"testing"

	"github.com/kn-lim/chattingway/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	// Run test
	expected := "Pong!"
	actual := healthcheck.Ping()

	assert.Equal(t, expected, actual)
}
