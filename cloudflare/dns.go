// Package cloudflare provides helpers for managing Cloudflare DNS records used by the chat bots.
package cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

// CreateDNSRecord creates an unproxied DNS A record named recordName in the given zone, pointing to ipAddress with an automatic TTL.
// The apiToken is used to authenticate with the Cloudflare API.
func CreateDNSRecord(ctx context.Context, apiToken string, zoneID string, recordName string, ipAddress string) error {
	client := cloudflare.NewClient(option.WithAPIToken(apiToken))

	_, err := client.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Body: dns.ARecordParam{
			Type:    cloudflare.F(dns.ARecordTypeA),
			Name:    cloudflare.F(recordName),
			Content: cloudflare.F(ipAddress),
			Proxied: cloudflare.F(false),
			TTL:     cloudflare.F(dns.TTL(1)),
		},
	})

	return err
}

// DeleteDNSRecord finds and deletes the A record named recordName in the given zone.
// It returns an error if no matching record exists or if more than one is found.
// The apiToken is used to authenticate with the Cloudflare API.
func DeleteDNSRecord(ctx context.Context, apiToken string, zoneID string, recordName string) error {
	client := cloudflare.NewClient(option.WithAPIToken(apiToken))

	page, err := client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(recordName),
		}),
		Type: cloudflare.F(dns.RecordListParamsTypeA),
	})
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}

	if len(page.Result) == 0 {
		return fmt.Errorf("DNS record not found: %s", recordName)
	} else if len(page.Result) > 1 {
		return fmt.Errorf("found multiple A records for %s", recordName)
	}

	_, err = client.DNS.Records.Delete(ctx, page.Result[0].ID, dns.RecordDeleteParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}
