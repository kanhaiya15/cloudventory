package rds

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
		INSERT INTO rds_instances (
			db_instance_identifier, db_instance_class, engine, engine_version,
			db_instance_status, master_username, db_name, allocated_storage,
			storage_type, encrypted, availability_zone, multi_az, vpc_id,
			subnet_group, security_groups, backup_retention_period,
			instance_create_time, last_updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW())
		ON CONFLICT (db_instance_identifier) DO UPDATE SET
			db_instance_class = EXCLUDED.db_instance_class,
			engine = EXCLUDED.engine,
			engine_version = EXCLUDED.engine_version,
			db_instance_status = EXCLUDED.db_instance_status,
			master_username = EXCLUDED.master_username,
			db_name = EXCLUDED.db_name,
			allocated_storage = EXCLUDED.allocated_storage,
			storage_type = EXCLUDED.storage_type,
			encrypted = EXCLUDED.encrypted,
			availability_zone = EXCLUDED.availability_zone,
			multi_az = EXCLUDED.multi_az,
			vpc_id = EXCLUDED.vpc_id,
			subnet_group = EXCLUDED.subnet_group,
			security_groups = EXCLUDED.security_groups,
			backup_retention_period = EXCLUDED.backup_retention_period,
			instance_create_time = EXCLUDED.instance_create_time,
			last_updated = NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, resource := range resources {
		_, err := stmt.ExecContext(ctx,
			resource.DBInstanceIdentifier,
			resource.DBInstanceClass,
			resource.Engine,
			resource.EngineVersion,
			resource.DBInstanceStatus,
			resource.MasterUsername,
			resource.DBName,
			resource.AllocatedStorage,
			resource.StorageType,
			resource.Encrypted,
			resource.AvailabilityZone,
			resource.MultiAZ,
			resource.VpcID,
			resource.SubnetGroup,
			resource.SecurityGroups,
			resource.BackupRetentionPeriod,
			resource.InstanceCreateTime,
		)
		if err != nil {
			return fmt.Errorf("failed to insert RDS instance %s: %w", resource.DBInstanceIdentifier, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}