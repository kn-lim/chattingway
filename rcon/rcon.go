package rcon

import (
	"fmt"

	"github.com/gorcon/rcon"
)

// Run executes the provided RCON command
func Run(host, port, password, command string) (string, error) {
	conn, err := rcon.Dial(fmt.Sprintf("%s:%s", host, port), password)
	if err != nil {
		return "", err
	}
	defer conn.Close() //nolint:errcheck

	resp, err := conn.Execute(command)
	if err != nil {
		return "", err
	}

	return resp, nil
}
