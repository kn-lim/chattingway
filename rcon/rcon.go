// Package rcon provides a helper for executing commands against a server over the RCON protocol.
package rcon

import (
	"fmt"

	"github.com/gorcon/rcon"
)

// Run opens an RCON connection to host:port, authenticates with password, executes command, and returns the server's response.
// The connection is closed before returning.
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
