// Package mcstatus reports the status of a Minecraft Java server via the mcstatus.io API.
package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BASE_URL is the mcstatus.io endpoint queried for Java server status.
const BASE_URL = "https://api.mcstatus.io/v2/status/java"

// MCStatusResponse holds the subset of the mcstatus.io response that is parsed.
type MCStatusResponse struct {
	// Online reports whether the server is reachable.
	Online bool `json:"online"`

	// Players holds player count information for the server.
	Players struct {
		// Online is the number of players currently connected.
		Online int `json:"online"`
	} `json:"players"`
}

// GetMCStatus queries mcstatus.io for the given Minecraft server and returns whether it is online and the number of players currently connected.
func GetMCStatus(serverURL string) (bool, int, error) {
	response, err := http.Get(fmt.Sprintf("%s/%s", BASE_URL, serverURL))
	if err != nil {
		return false, 0, err
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, 0, err
	}

	var status MCStatusResponse
	err = json.Unmarshal(body, &status)
	if err != nil {
		return false, 0, err
	}

	return status.Online, status.Players.Online, nil
}
