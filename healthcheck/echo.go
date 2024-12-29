package healthcheck

import "fmt"

// Echo returns a string with the user and message received
func Echo(user, msg string) string {
	return fmt.Sprintf("Received Echo from %s: `%v`", user, msg)
}
