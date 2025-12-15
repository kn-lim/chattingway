package cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

// CreateDNSRecord creates a DNS A record pointing to the given IP address
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

// DeleteDNSRecord finds and deletes a DNS record by name
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
