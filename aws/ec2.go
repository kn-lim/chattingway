package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

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
