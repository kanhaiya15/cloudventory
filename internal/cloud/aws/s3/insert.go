package s3

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
		INSERT INTO s3_buckets (name, region, creation_date, versioning, encryption, public_access, last_updated)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (name) DO UPDATE SET
			region = EXCLUDED.region,
			creation_date = EXCLUDED.creation_date,
			versioning = EXCLUDED.versioning,
			encryption = EXCLUDED.encryption,
			public_access = EXCLUDED.public_access,
			last_updated = NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, resource := range resources {
		_, err := stmt.ExecContext(ctx, 
			resource.Name,
			resource.Region,
			resource.CreationDate,
			resource.Versioning,
			resource.Encryption,
			resource.PublicAccess,
		)
		if err != nil {
			return fmt.Errorf("failed to insert S3 bucket %s: %w", resource.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}