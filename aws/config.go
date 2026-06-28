// Package aws provides helpers for interacting with AWS services used by the chat bots.
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// loadConfig loads the default AWS configuration scoped to the given region.
func loadConfig(ctx context.Context, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}
