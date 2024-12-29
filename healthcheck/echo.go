package healthcheck

import "fmt"

func Echo(user, msg string) string {
	return fmt.Sprintf("Received Echo from %s: `%v`", user, msg)
}
