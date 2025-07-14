package ec2

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
		INSERT INTO ec2_instances (
			instance_id, instance_type, state, availability_zone, 
			public_ip, private_ip, launch_time, image_id, 
			vpc_id, subnet_id, security_groups, tags, last_updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		ON CONFLICT (instance_id) DO UPDATE SET
			instance_type = EXCLUDED.instance_type,
			state = EXCLUDED.state,
			availability_zone = EXCLUDED.availability_zone,
			public_ip = EXCLUDED.public_ip,
			private_ip = EXCLUDED.private_ip,
			launch_time = EXCLUDED.launch_time,
			image_id = EXCLUDED.image_id,
			vpc_id = EXCLUDED.vpc_id,
			subnet_id = EXCLUDED.subnet_id,
			security_groups = EXCLUDED.security_groups,
			tags = EXCLUDED.tags,
			last_updated = NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, resource := range resources {
		_, err := stmt.ExecContext(ctx,
			resource.InstanceID,
			resource.InstanceType,
			resource.State,
			resource.AvailabilityZone,
			resource.PublicIP,
			resource.PrivateIP,
			resource.LaunchTime,
			resource.ImageID,
			resource.VpcID,
			resource.SubnetID,
			resource.SecurityGroups,
			resource.Tags,
		)
		if err != nil {
			return fmt.Errorf("failed to insert EC2 instance %s: %w", resource.InstanceID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}