package projectzomboid

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/kn-lim/chattingway/aws"
	"github.com/kn-lim/chattingway/rcon"
)

const (
	STATUS_CHECK_INTERVAL = 30 // Status check interval in seconds
)

// Start starts the Project Zomboid server
func Start() error {
	if err := aws.StartInstance(context.TODO(), os.Getenv("PZ_HOST_INSTANCE_ID"), os.Getenv("PZ_HOST_REGION")); err != nil {
		return err
	}

	for {
		time.Sleep(time.Duration(STATUS_CHECK_INTERVAL) * time.Second)

		if status, _ := Status(); status {
			break
		}
	}

	return nil
}

// Status returns whether the Project Zomboid server is online or offline
func Status() (bool, error) {
	output, err := rcon.Run(os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PASSWORD"), "players")
	if err != nil {
		return false, err
	}

	if output == "" {
		return false, errors.New("received empty output for status")
	}

	return true, nil
}

// Stop stops the Project Zomboid server
func Stop() error {
	_, err := rcon.Run(os.Getenv("PZ_HOST"), os.Getenv("PZ_RCON_PASSWORD"), "quit")
	if err != nil {
		return err
	}

	if err := aws.StopInstance(context.TODO(), os.Getenv("PZ_HOST_INSTANCE_ID"), os.Getenv("PZ_HOST_REGION")); err != nil {
		return err
	}

	for {
		time.Sleep(time.Duration(STATUS_CHECK_INTERVAL) * time.Second)

		if status, _ := Status(); !status {
			break
		}
	}

	return nil
}
