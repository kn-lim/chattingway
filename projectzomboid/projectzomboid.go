package projectzomboid

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/kn-lim/chattingway/aws"
	"github.com/kn-lim/chattingway/cloudflare"
	"github.com/kn-lim/chattingway/rcon"
)

const (
	WAIT_TIME             = 10 // Wait time in seconds
	STATUS_CHECK_INTERVAL = 30 // Status check interval in seconds
)

// Start starts the Project Zomboid server
func Start(ctx context.Context, instanceID, region, host, port, password, cfToken, cfZoneID, cfRecordName string) error {
	if err := aws.StartInstance(ctx, instanceID, region); err != nil {
		return err
	}

	for {
		output, err := aws.GetInstanceState(ctx, instanceID, region)
		if err != nil {
			return err
		}

		if output == string(types.InstanceStateNameRunning) {
			break
		}

		time.Sleep(time.Duration(WAIT_TIME) * time.Second)
	}

	publicIP, err := aws.GetInstancePublicIP(ctx, instanceID, region)
	if err != nil {
		return err
	}

	if err := cloudflare.CreateDNSRecord(ctx, cfToken, cfZoneID, cfRecordName, publicIP); err != nil {
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
func Stop(ctx context.Context, instanceID, region, host, port, password, cfToken, cfZoneID, cfRecordName string) error {
	var err error
	_, err = rcon.Run(host, port, password, "save")
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(WAIT_TIME) * time.Second)

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

	if err := cloudflare.DeleteDNSRecord(ctx, cfToken, cfZoneID, cfRecordName); err != nil {
		return err
	}

	return nil
}
