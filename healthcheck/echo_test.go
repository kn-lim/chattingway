package healthcheck_test

import (
	"fmt"
	"testing"

	"github.com/kn-lim/chattingway/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	user := "testUser"
	msg := "testMessage"

	expected := fmt.Sprintf("Received Echo from %s: `%v`", user, msg)
	actual := healthcheck.Echo(user, msg)

	assert.Equal(t, expected, actual)
}
