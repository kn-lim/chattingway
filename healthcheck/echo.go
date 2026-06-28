package healthcheck

import "fmt"

// Echo returns a formatted string echoing back the message received from the given user.
func Echo(user, msg string) string {
	return fmt.Sprintf("Received Echo from %s: `%s`", user, msg)
}
