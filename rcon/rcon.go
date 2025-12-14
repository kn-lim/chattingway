package rcon

import "github.com/gorcon/rcon"

// Run executes the provided RCON command
func Run(host, password, command string) (string, error) {
	conn, err := rcon.Dial(host, password)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	resp, err := conn.Execute(command)
	if err != nil {
		return "", err
	}

	return resp, nil
}
