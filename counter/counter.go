// Package counter stores named counters per Discord guild in a DynamoDB table.
package counter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Counter is one named counter in one guild.
type Counter struct {
	Name    string `dynamodbav:"name"`
	Count   int64  `dynamodbav:"count"`
	Message string `dynamodbav:"message"`
}

// CooldownError reports that the counter changed too recently.
type CooldownError struct {
	// Remaining is the time left before the next increment can succeed.
	Remaining time.Duration
}

func (e *CooldownError) Error() string {
	return fmt.Sprintf("counter is on cooldown for %s", e.Remaining)
}

// newClient returns a DynamoDB client for the given region.
func newClient(ctx context.Context, region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

// key returns the primary key for one counter.
func key(guildID, name string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"guild_id": &types.AttributeValueMemberS{Value: guildID},
		"name":     &types.AttributeValueMemberS{Value: name},
	}
}

// Increment adds one to a counter and returns the new value.
// If cooldown is more than zero, the counter must be older than the cooldown.
// If the counter changed too recently, Increment returns a CooldownError.
// If the counter does not exist, Increment creates it with a count of one.
func Increment(ctx context.Context, table, region, guildID, name string, cooldown time.Duration) (Counter, error) {
	client, err := newClient(ctx, region)
	if err != nil {
		return Counter{}, err
	}

	now := time.Now().Unix()

	values := map[string]types.AttributeValue{
		":one": &types.AttributeValueMemberN{Value: "1"},
		":now": &types.AttributeValueMemberN{Value: strconv.FormatInt(now, 10)},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(table),
		Key:              key(guildID, name),
		UpdateExpression: aws.String("ADD #c :one SET last_incremented_at = :now"),
		// count is a DynamoDB reserved word.
		ExpressionAttributeNames:  map[string]string{"#c": "count"},
		ExpressionAttributeValues: values,
		ReturnValues:              types.ReturnValueAllNew,
	}

	// Add the condition and its value together. DynamoDB rejects a value that no expression uses.
	if cooldown > 0 {
		cutoff := now - int64(cooldown.Seconds())
		values[":cutoff"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(cutoff, 10)}
		input.ConditionExpression = aws.String(
			"attribute_not_exists(last_incremented_at) OR last_incremented_at < :cutoff")
		input.ReturnValuesOnConditionCheckFailure = types.ReturnValuesOnConditionCheckFailureAllOld
	}

	out, err := client.UpdateItem(ctx, input)
	if err != nil {
		var failed *types.ConditionalCheckFailedException
		if errors.As(err, &failed) {
			return Counter{}, &CooldownError{Remaining: remaining(failed.Item, now, cooldown)}
		}

		return Counter{}, err
	}

	return unmarshal(out.Attributes, name)
}

// remaining returns the time left in the cooldown.
// If the item cannot be read, remaining returns the full cooldown.
func remaining(item map[string]types.AttributeValue, now int64, cooldown time.Duration) time.Duration {
	if item == nil {
		return cooldown
	}

	var old struct {
		LastIncrementedAt int64 `dynamodbav:"last_incremented_at"`
	}
	if err := attributevalue.UnmarshalMap(item, &old); err != nil {
		return cooldown
	}

	left := cooldown - time.Duration(now-old.LastIncrementedAt)*time.Second
	if left < 0 {
		return 0
	}

	return left
}

// Set gives a counter an exact value and returns the new value.
func Set(ctx context.Context, table, region, guildID, name string, value int64) (Counter, error) {
	client, err := newClient(ctx, region)
	if err != nil {
		return Counter{}, err
	}

	out, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                aws.String(table),
		Key:                      key(guildID, name),
		UpdateExpression:         aws.String("SET #c = :v"),
		ExpressionAttributeNames: map[string]string{"#c": "count"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberN{Value: strconv.FormatInt(value, 10)},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return Counter{}, err
	}

	return unmarshal(out.Attributes, name)
}

// SetMessage gives a counter a new reply template and returns the counter.
func SetMessage(ctx context.Context, table, region, guildID, name, message string) (Counter, error) {
	client, err := newClient(ctx, region)
	if err != nil {
		return Counter{}, err
	}

	out, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                aws.String(table),
		Key:                      key(guildID, name),
		UpdateExpression:         aws.String("SET #m = :m"),
		ExpressionAttributeNames: map[string]string{"#m": "message"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":m": &types.AttributeValueMemberS{Value: message},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return Counter{}, err
	}

	return unmarshal(out.Attributes, name)
}

// List returns every counter in one guild.
func List(ctx context.Context, table, region, guildID string) ([]Counter, error) {
	client, err := newClient(ctx, region)
	if err != nil {
		return nil, err
	}

	var counters []Counter

	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("guild_id = :g"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":g": &types.AttributeValueMemberS{Value: guildID},
		},
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		var pageCounters []Counter
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &pageCounters); err != nil {
			return nil, err
		}

		counters = append(counters, pageCounters...)
	}

	return counters, nil
}

// unmarshal converts DynamoDB attributes into a Counter struct.
func unmarshal(item map[string]types.AttributeValue, name string) (Counter, error) {
	var c Counter
	if err := attributevalue.UnmarshalMap(item, &c); err != nil {
		return Counter{}, err
	}

	if c.Name == "" {
		c.Name = name
	}

	return c, nil
}
