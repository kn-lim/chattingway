package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// StartInstance starts the EC2 instance identified by instanceID in the given region.
// It is a no-op if the instance is not currently in the stopped state.
func StartInstance(ctx context.Context, instanceID, region string) error {
	// Setup AWS session
	cfg, err := getConfig(ctx, region)
	if err != nil {
		return err
	}

	// Get EC2 instance
	client, instance, err := getInstance(ctx, cfg, instanceID)
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
	cfg, err := getConfig(ctx, region)
	if err != nil {
		return err
	}

	// Get EC2 instance
	client, instance, err := getInstance(ctx, cfg, instanceID)
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

// GetInstancePublicIP returns the public IPv4 address of the EC2 instance identified by instanceID in the given region.
// It returns an error if the instance has no public IP assigned.
func GetInstancePublicIP(ctx context.Context, instanceID, region string) (string, error) {
	// Setup AWS session
	cfg, err := getConfig(ctx, region)
	if err != nil {
		return "", err
	}

	// Get EC2 instance
	_, instance, err := getInstance(ctx, cfg, instanceID)
	if err != nil {
		return "", err
	}

	// Return public IP address
	if instance.PublicIpAddress == nil {
		return "", errors.New(ERR_INSTANCE_NO_PUBLIC_IP)
	}

	return *instance.PublicIpAddress, nil
}

// GetInstanceState returns the current lifecycle state of the EC2 instance identified by instanceID in the given region
// (for example "running" or "stopped").
func GetInstanceState(ctx context.Context, instanceID, region string) (string, error) {
	// Setup AWS session
	cfg, err := getConfig(ctx, region)
	if err != nil {
		return "", err
	}

	// Get EC2 instance
	_, instance, err := getInstance(ctx, cfg, instanceID)
	if err != nil {
		return "", err
	}

	return string(instance.State.Name), nil
}

// getInstance describes the instance identified by instanceID and returns the EC2 client along with the resolved instance.
func getInstance(ctx context.Context, cfg aws.Config, instanceID string) (*ec2.Client, types.Instance, error) {
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
		return nil, types.Instance{}, err
	}

	return client, result.Reservations[0].Instances[0], nil
}
