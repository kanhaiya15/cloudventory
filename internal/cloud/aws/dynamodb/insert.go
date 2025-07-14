package dynamodb

import (
	"context"
	"database/sql"
	"fmt"
)

func InsertResources(ctx context.Context, db *sql.DB, resources []Resource) error {
	if len(resources) == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO dynamodb_tables (
			table_name, table_arn, table_status, creation_date_time,
			provisioned_throughput, billing_mode, item_count, table_size_bytes,
			global_secondary_indexes, local_secondary_indexes, stream_specification,
			sse_description, point_in_time_recovery, tags, last_updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW())
		ON CONFLICT (table_name) DO UPDATE SET
			table_arn = EXCLUDED.table_arn,
			table_status = EXCLUDED.table_status,
			creation_date_time = EXCLUDED.creation_date_time,
			provisioned_throughput = EXCLUDED.provisioned_throughput,
			billing_mode = EXCLUDED.billing_mode,
			item_count = EXCLUDED.item_count,
			table_size_bytes = EXCLUDED.table_size_bytes,
			global_secondary_indexes = EXCLUDED.global_secondary_indexes,
			local_secondary_indexes = EXCLUDED.local_secondary_indexes,
			stream_specification = EXCLUDED.stream_specification,
			sse_description = EXCLUDED.sse_description,
			point_in_time_recovery = EXCLUDED.point_in_time_recovery,
			tags = EXCLUDED.tags,
			last_updated = NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, resource := range resources {
		_, err := stmt.ExecContext(ctx,
			resource.TableName,
			resource.TableArn,
			resource.TableStatus,
			resource.CreationDateTime,
			resource.ProvisionedThroughput,
			resource.BillingMode,
			resource.ItemCount,
			resource.TableSizeBytes,
			resource.GlobalSecondaryIndexes,
			resource.LocalSecondaryIndexes,
			resource.StreamSpecification,
			resource.SSEDescription,
			resource.PointInTimeRecovery,
			resource.Tags,
		)
		if err != nil {
			return fmt.Errorf("failed to insert DynamoDB table %s: %w", resource.TableName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}