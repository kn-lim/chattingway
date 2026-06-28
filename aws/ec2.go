package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ErrNoPublicIP is returned when an instance has no public IP address assigned.
var ErrNoPublicIP = errors.New("aws: instance does not have a public IP")

// ErrInstanceNotFound is returned when no instance matches the given ID.
var ErrInstanceNotFound = errors.New("aws: instance not found")

// StartInstance starts the EC2 instance identified by instanceID in the given region.
// It is a no-op if the instance is not currently in the stopped state.
func StartInstance(ctx context.Context, instanceID, region string) error {
	// Setup AWS session
	cfg, err := loadConfig(ctx, region)
	if err != nil {
		return err
	}

	// Get EC2 instance
	client, instance, err := describeInstance(ctx, cfg, instanceID)
	if err != nil {
		return err
	}

	// Start EC2 instance
	if instance.State.Name == types.InstanceStateNameStopped {
		input := &ec2.StartInstancesInput{
			InstanceIds: []string{*instance.InstanceId},
		}

		_, err := client.StartInstances(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}

// StopInstance stops the EC2 instance identified by instanceID in the given region.
// It is a no-op if the instance is not currently in the running state.
func StopInstance(ctx context.Context, instanceID, region string) error {
	// Setup AWS session
	cfg, err := loadConfig(ctx, region)
	if err != nil {
		return err
	}

	// Get EC2 instance
	client, instance, err := describeInstance(ctx, cfg, instanceID)
	if err != nil {
		return err
	}

	// Stop EC2 instance
	if instance.State.Name == types.InstanceStateNameRunning {
		input := &ec2.StopInstancesInput{
			InstanceIds: []string{*instance.InstanceId},
		}

		_, err := client.StopInstances(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}

// InstancePublicIP returns the public IPv4 address of the EC2 instance identified by instanceID in the given region.
// It returns ErrNoPublicIP if the instance has no public IP assigned.
func InstancePublicIP(ctx context.Context, instanceID, region string) (string, error) {
	// Setup AWS session
	cfg, err := loadConfig(ctx, region)
	if err != nil {
		return "", err
	}

	// Get EC2 instance
	_, instance, err := describeInstance(ctx, cfg, instanceID)
	if err != nil {
		return "", err
	}

	// Return public IP address
	if instance.PublicIpAddress == nil {
		return "", ErrNoPublicIP
	}

	return *instance.PublicIpAddress, nil
}

// InstanceState returns the current lifecycle state of the EC2 instance identified by instanceID in the given region
// (for example "running" or "stopped").
func InstanceState(ctx context.Context, instanceID, region string) (string, error) {
	// Setup AWS session
	cfg, err := loadConfig(ctx, region)
	if err != nil {
		return "", err
	}

	// Get EC2 instance
	_, instance, err := describeInstance(ctx, cfg, instanceID)
	if err != nil {
		return "", err
	}

	return string(instance.State.Name), nil
}

// describeInstance describes the instance identified by instanceID and returns the EC2 client along with the resolved instance.
// It returns ErrInstanceNotFound if no matching instance exists.
func describeInstance(ctx context.Context, cfg aws.Config, instanceID string) (*ec2.Client, types.Instance, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	}

	result, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, types.Instance{}, err
	}
	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, types.Instance{}, ErrInstanceNotFound
	}

	return client, result.Reservations[0].Instances[0], nil
}
