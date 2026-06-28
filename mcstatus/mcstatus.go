// Package mcstatus reports the status of a Minecraft Java server via the mcstatus.io API.
package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// baseURL is the mcstatus.io endpoint queried for Java server status.
const baseURL = "https://api.mcstatus.io/v2/status/java"

// response holds the subset of the mcstatus.io response that is parsed.
type response struct {
	// Online reports whether the server is reachable.
	Online bool `json:"online"`

	// Players holds player count information for the server.
	Players struct {
		// Online is the number of players currently connected.
		Online int `json:"online"`
	} `json:"players"`
}

// Query queries mcstatus.io for the given Minecraft server and returns whether it is online and the number of players currently connected.
func Query(serverURL string) (bool, int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, serverURL))
	if err != nil {
		return false, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}

	var status response
	if err := json.Unmarshal(body, &status); err != nil {
		return false, 0, err
	}

	return status.Online, status.Players.Online, nil
}
