package rds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type Resource struct {
	DBInstanceIdentifier   string `json:"db_instance_identifier"`
	DBInstanceClass       string `json:"db_instance_class"`
	Engine                string `json:"engine"`
	EngineVersion         string `json:"engine_version"`
	DBInstanceStatus      string `json:"db_instance_status"`
	MasterUsername        string `json:"master_username"`
	DBName                string `json:"db_name"`
	AllocatedStorage      int32  `json:"allocated_storage"`
	StorageType           string `json:"storage_type"`
	Encrypted             bool   `json:"encrypted"`
	AvailabilityZone      string `json:"availability_zone"`
	MultiAZ               bool   `json:"multi_az"`
	VpcID                 string `json:"vpc_id"`
	SubnetGroup           string `json:"subnet_group"`
	SecurityGroups        string `json:"security_groups"`
	BackupRetentionPeriod int32  `json:"backup_retention_period"`
	InstanceCreateTime    string `json:"instance_create_time"`
}

func FetchResources(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := rds.NewFromConfig(cfg)
	
	result, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe RDS instances: %w", err)
	}

	var resources []Resource
	for _, dbInstance := range result.DBInstances {
		resource := Resource{
			DBInstanceIdentifier:   aws.ToString(dbInstance.DBInstanceIdentifier),
			DBInstanceClass:       aws.ToString(dbInstance.DBInstanceClass),
			Engine:                aws.ToString(dbInstance.Engine),
			EngineVersion:         aws.ToString(dbInstance.EngineVersion),
			DBInstanceStatus:      aws.ToString(dbInstance.DBInstanceStatus),
			MasterUsername:        aws.ToString(dbInstance.MasterUsername),
			DBName:                aws.ToString(dbInstance.DBName),
			AllocatedStorage:      aws.ToInt32(dbInstance.AllocatedStorage),
			StorageType:           aws.ToString(dbInstance.StorageType),
			Encrypted:             aws.ToBool(dbInstance.StorageEncrypted),
			AvailabilityZone:      aws.ToString(dbInstance.AvailabilityZone),
			MultiAZ:               aws.ToBool(dbInstance.MultiAZ),
			BackupRetentionPeriod: aws.ToInt32(dbInstance.BackupRetentionPeriod),
		}

		if dbInstance.InstanceCreateTime != nil {
			resource.InstanceCreateTime = dbInstance.InstanceCreateTime.String()
		}

		if dbInstance.DBSubnetGroup != nil {
			resource.SubnetGroup = aws.ToString(dbInstance.DBSubnetGroup.DBSubnetGroupName)
			if dbInstance.DBSubnetGroup.VpcId != nil {
				resource.VpcID = aws.ToString(dbInstance.DBSubnetGroup.VpcId)
			}
		}

		// Extract security groups
		var sgNames []string
		for _, sg := range dbInstance.VpcSecurityGroups {
			sgNames = append(sgNames, aws.ToString(sg.VpcSecurityGroupId))
		}
		resource.SecurityGroups = fmt.Sprintf("[%s]", joinStrings(sgNames, ","))

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