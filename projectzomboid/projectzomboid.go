// Package projectzomboid orchestrates the lifecycle of a Project Zomboid game server, coordinating its AWS EC2 host,
// Cloudflare DNS record, and RCON connection.
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
	// waitTime is the delay, in seconds, used between sequential server operations such as saving and shutting down.
	waitTime = 10

	// statusCheckInterval is the interval, in seconds, between server status polls while waiting for the server to
	// come online or go offline.
	statusCheckInterval = 30
)

// Start boots the Project Zomboid server: it starts the EC2 host, waits for it to run, points the Cloudflare DNS record
// at the host's public IP, and blocks until the game server responds to RCON.
func Start(ctx context.Context, instanceID, region, host, port, password, cfToken, cfZoneID, cfRecordName string) error {
	if err := aws.StartInstance(ctx, instanceID, region); err != nil {
		return err
	}

	for {
		output, err := aws.InstanceState(ctx, instanceID, region)
		if err != nil {
			if !errors.Is(err, aws.ErrNoPublicIP) {
				return err
			}
		}

		if output == string(types.InstanceStateNameRunning) {
			break
		}

		time.Sleep(time.Duration(waitTime) * time.Second)
	}

	publicIP, err := aws.InstancePublicIP(ctx, instanceID, region)
	if err != nil {
		return err
	}

	if err := cloudflare.CreateDNSRecord(ctx, cfToken, cfZoneID, cfRecordName, publicIP); err != nil {
		return err
	}

	for {
		time.Sleep(time.Duration(statusCheckInterval) * time.Second)

		if status, _ := Status(host, port, password); status {
			break
		}
	}

	return nil
}

// Status reports whether the Project Zomboid server is online by issuing an RCON "players" command to the given host.
func Status(host, port, password string) (bool, error) {
	output, err := rcon.Run(host, port, password, "players")
	if err != nil {
		return false, err
	}

	if output == "" {
		return false, errors.New("projectzomboid: received empty output for status")
	}

	return true, nil
}

// Stop shuts down the Project Zomboid server: it saves and quits the game via RCON, waits for the server to go offline,
// stops the EC2 host, and removes the Cloudflare DNS record.
func Stop(ctx context.Context, instanceID, region, host, port, password, cfToken, cfZoneID, cfRecordName string) error {
	var err error
	_, err = rcon.Run(host, port, password, "save")
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(waitTime) * time.Second)

	_, err = rcon.Run(host, port, password, "quit")
	if err != nil {
		return err
	}

	for {
		time.Sleep(time.Duration(statusCheckInterval) * time.Second)

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
