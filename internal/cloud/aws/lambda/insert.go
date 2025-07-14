package lambda

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
		INSERT INTO lambda_functions (
			function_name, function_arn, runtime, role, handler, code_size,
			description, timeout, memory_size, last_modified, code_sha256,
			version, environment, vpc_config, dead_letter_config, state,
			state_reason, last_updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW())
		ON CONFLICT (function_name) DO UPDATE SET
			function_arn = EXCLUDED.function_arn,
			runtime = EXCLUDED.runtime,
			role = EXCLUDED.role,
			handler = EXCLUDED.handler,
			code_size = EXCLUDED.code_size,
			description = EXCLUDED.description,
			timeout = EXCLUDED.timeout,
			memory_size = EXCLUDED.memory_size,
			last_modified = EXCLUDED.last_modified,
			code_sha256 = EXCLUDED.code_sha256,
			version = EXCLUDED.version,
			environment = EXCLUDED.environment,
			vpc_config = EXCLUDED.vpc_config,
			dead_letter_config = EXCLUDED.dead_letter_config,
			state = EXCLUDED.state,
			state_reason = EXCLUDED.state_reason,
			last_updated = NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, resource := range resources {
		_, err := stmt.ExecContext(ctx,
			resource.FunctionName,
			resource.FunctionArn,
			resource.Runtime,
			resource.Role,
			resource.Handler,
			resource.CodeSize,
			resource.Description,
			resource.Timeout,
			resource.MemorySize,
			resource.LastModified,
			resource.CodeSha256,
			resource.Version,
			resource.Environment,
			resource.VpcConfig,
			resource.DeadLetterConfig,
			resource.State,
			resource.StateReason,
		)
		if err != nil {
			return fmt.Errorf("failed to insert Lambda function %s: %w", resource.FunctionName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}