// Package healthcheck provides simple liveness commands for verifying that a bot is responsive.
package healthcheck

// Ping returns "Pong!".
func Ping() string {
	return "Pong!"
}
