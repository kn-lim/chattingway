package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// ROUTE53_TTL is the time-to-live, in seconds, applied to managed Route 53 records.
const ROUTE53_TTL = 300

// CreateRoute53Record upserts an A record named url in the given hosted zone so that it resolves to publicIP.
// An existing record with the same name is overwritten.
func CreateRoute53Record(ctx context.Context, cfg aws.Config, publicIP, zoneID, url string) error {
	client := route53.NewFromConfig(cfg)

	_, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &url,
						Type: types.RRTypeA,
						TTL:  aws.Int64(ROUTE53_TTL),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: &publicIP,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteRoute53Record deletes the A record named url that resolves to publicIP from the given hosted zone.
func DeleteRoute53Record(ctx context.Context, cfg aws.Config, publicIP, zoneID, url string) error {
	client := route53.NewFromConfig(cfg)

	_, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionDelete,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &url,
						Type: types.RRTypeA,
						TTL:  aws.Int64(ROUTE53_TTL),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: &publicIP,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
