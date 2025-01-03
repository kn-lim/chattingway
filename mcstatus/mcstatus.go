package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BASE_URL = "https://api.mcstatus.io/v2/status/java"
)

type MCStatusResponse struct {
	Online  bool `json:"online"`
	Players struct {
		Online int `json:"online"`
	} `json:"players"`
}

// GetMCStatus checks with mcstatus.io to get information about the Minecraft server
func GetMCStatus(serverURL string) (bool, int, error) {
	response, err := http.Get(fmt.Sprintf("%s/%s", BASE_URL, serverURL))
	if err != nil {
		return false, 0, err
	}
	defer response.Body.Close()

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
