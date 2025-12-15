package projectzomboid

import (
	"context"
	"errors"
	"time"

	"github.com/kn-lim/chattingway/aws"
	"github.com/kn-lim/chattingway/rcon"
)

const (
	SAVE_WAIT_TIME        = 10 // Save wait time in seconds
	STATUS_CHECK_INTERVAL = 30 // Status check interval in seconds
)

// Start starts the Project Zomboid server
func Start(ctx context.Context, instanceID, region, host, port, password string) error {
	if err := aws.StartInstance(ctx, instanceID, region); err != nil {
		return err
	}

	for {
		time.Sleep(time.Duration(STATUS_CHECK_INTERVAL) * time.Second)

		if status, _ := Status(host, port, password); status {
			break
		}
	}

	return nil
}

// Status returns whether the Project Zomboid server is online or offline
func Status(host, port, password string) (bool, error) {
	output, err := rcon.Run(host, port, password, "players")
	if err != nil {
		return false, err
	}

	if output == "" {
		return false, errors.New("received empty output for status")
	}

	return true, nil
}

// Stop stops the Project Zomboid server
func Stop(ctx context.Context, instanceID, region, host, port, password string) error {
	var err error
	_, err = rcon.Run(host, port, password, "save")
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(SAVE_WAIT_TIME) * time.Second)

	_, err = rcon.Run(host, port, password, "quit")
	if err != nil {
		return err
	}

	for {
		time.Sleep(time.Duration(STATUS_CHECK_INTERVAL) * time.Second)

		if status, _ := Status(host, port, password); !status {
			break
		}
	}

	if err := aws.StopInstance(ctx, instanceID, region); err != nil {
		return err
	}

	return nil
}
