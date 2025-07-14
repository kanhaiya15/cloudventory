package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Resource struct {
	TableName           string `json:"table_name"`
	TableArn            string `json:"table_arn"`
	TableStatus         string `json:"table_status"`
	CreationDateTime    string `json:"creation_date_time"`
	ProvisionedThroughput string `json:"provisioned_throughput"`
	BillingMode         string `json:"billing_mode"`
	ItemCount           int64  `json:"item_count"`
	TableSizeBytes      int64  `json:"table_size_bytes"`
	GlobalSecondaryIndexes string `json:"global_secondary_indexes"`
	LocalSecondaryIndexes  string `json:"local_secondary_indexes"`
	StreamSpecification    string `json:"stream_specification"`
	SSEDescription         string `json:"sse_description"`
	PointInTimeRecovery    string `json:"point_in_time_recovery"`
	Tags                   string `json:"tags"`
}

func FetchResources(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodb.NewFromConfig(cfg)
	
	result, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list DynamoDB tables: %w", err)
	}

	var resources []Resource
	for _, tableName := range result.TableNames {
		// Get detailed table information
		describeResult, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			continue // Skip this table if we can't describe it
		}

		table := describeResult.Table
		resource := Resource{
			TableName:    aws.ToString(table.TableName),
			TableArn:     aws.ToString(table.TableArn),
			TableStatus:  string(table.TableStatus),
			ItemCount:    aws.ToInt64(table.ItemCount),
			TableSizeBytes: aws.ToInt64(table.TableSizeBytes),
			BillingMode:  string(table.BillingModeSummary.BillingMode),
		}

		if table.CreationDateTime != nil {
			resource.CreationDateTime = table.CreationDateTime.String()
		}

		// Extract provisioned throughput
		if table.ProvisionedThroughput != nil {
			resource.ProvisionedThroughput = fmt.Sprintf("Read:%d,Write:%d",
				aws.ToInt64(table.ProvisionedThroughput.ReadCapacityUnits),
				aws.ToInt64(table.ProvisionedThroughput.WriteCapacityUnits))
		}

		// Extract Global Secondary Indexes
		if len(table.GlobalSecondaryIndexes) > 0 {
			var gsiInfo []string
			for _, gsi := range table.GlobalSecondaryIndexes {
				gsiInfo = append(gsiInfo, aws.ToString(gsi.IndexName))
			}
			resource.GlobalSecondaryIndexes = fmt.Sprintf("[%s]", joinStrings(gsiInfo, ","))
		}

		// Extract Local Secondary Indexes
		if len(table.LocalSecondaryIndexes) > 0 {
			var lsiInfo []string
			for _, lsi := range table.LocalSecondaryIndexes {
				lsiInfo = append(lsiInfo, aws.ToString(lsi.IndexName))
			}
			resource.LocalSecondaryIndexes = fmt.Sprintf("[%s]", joinStrings(lsiInfo, ","))
		}

		// Extract Stream Specification
		if table.StreamSpecification != nil {
			resource.StreamSpecification = fmt.Sprintf("Enabled:%t,StreamViewType:%s",
				aws.ToBool(table.StreamSpecification.StreamEnabled),
				string(table.StreamSpecification.StreamViewType))
		}

		// Extract SSE Description
		if table.SSEDescription != nil {
			resource.SSEDescription = string(table.SSEDescription.Status)
		}

		// Get Point-in-Time Recovery status
		pitrResult, err := client.DescribeContinuousBackups(ctx, &dynamodb.DescribeContinuousBackupsInput{
			TableName: aws.String(tableName),
		})
		if err == nil && pitrResult.ContinuousBackupsDescription != nil {
			resource.PointInTimeRecovery = string(pitrResult.ContinuousBackupsDescription.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus)
		}

		// Get table tags
		tagsResult, err := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
			ResourceArn: table.TableArn,
		})
		if err == nil && len(tagsResult.Tags) > 0 {
			var tagPairs []string
			for _, tag := range tagsResult.Tags {
				tagPairs = append(tagPairs, fmt.Sprintf("%s=%s",
					aws.ToString(tag.Key), aws.ToString(tag.Value)))
			}
			resource.Tags = fmt.Sprintf("{%s}", joinStrings(tagPairs, ","))
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func joinStrings(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}