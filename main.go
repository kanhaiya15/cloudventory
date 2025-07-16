package main

import (
	"context"
	"encoding/json"
	"fmt"
	awsCventSvc "kanhaiya1501/cloudventory/aws"
	"kanhaiya1501/cloudventory/aws/cfg"
	ddbinventory "kanhaiya1501/cloudventory/aws/ddb/inventory"
	ec2inventory "kanhaiya1501/cloudventory/aws/ec2/inventory"
	elbinventory "kanhaiya1501/cloudventory/aws/elb/inventory"
	iaminventory "kanhaiya1501/cloudventory/aws/iam/inventory"
	kmsinventory "kanhaiya1501/cloudventory/aws/kms/inventory"
	rdsinventory "kanhaiya1501/cloudventory/aws/rds/inventory"
	s3inventory "kanhaiya1501/cloudventory/aws/s3/inventory"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	awsCfg, cfgErr := cfg.LoadAWSConfig(ctx)
	if cfgErr != nil {
		log.Fatalf("cfg error: %v", cfgErr)
	}
	services := []awsCventSvc.InventoryService{
		&s3inventory.S3InventoryService{
			Client: s3inventory.New(awsCfg),
			Options: s3inventory.Options{
				ChunkSize:    100,
				RegionFilter: []string{"us-east-1", "us-west-2"},
				MaxRetries:   3,
				Timeout:      30 * time.Minute,
			},
		},
		&ec2inventory.EC2InventoryService{
			Client: ec2inventory.New(awsCfg),
			Options: ec2inventory.Options{
				ChunkSize:    100,
				RegionFilter: []string{"us-east-1", "us-west-2"},
				MaxRetries:   3,
				Timeout:      30 * time.Minute,
			},
		},
		&rdsinventory.RDSInventoryService{
			Client: rdsinventory.New(awsCfg),
			Options: rdsinventory.Options{
				ChunkSize:  100,
				MaxRetries: 3,
				Timeout:    30 * time.Minute,
			},
		},
		&iaminventory.IAMInventoryService{
			Client: iaminventory.New(awsCfg),
			Options: iaminventory.Options{
				ChunkSize:    100,
				RegionFilter: []string{"us-east-1", "us-west-2"},
				MaxRetries:   3,
				Timeout:      30 * time.Minute,
			},
		},
		&kmsinventory.KMSInventoryService{
			Client: kmsinventory.New(awsCfg),
			Options: kmsinventory.Options{
				ChunkSize:    100,
				RegionFilter: []string{"us-east-1", "us-west-2"},
				MaxRetries:   3,
				Timeout:      30 * time.Minute,
			},
		},
		&elbinventory.ELBInventoryService{
			Client: elbinventory.New(awsCfg),
			Options: elbinventory.Options{
				ChunkSize:    100,
				RegionFilter: []string{"us-east-1", "us-west-2"},
				MaxRetries:   3,
				Timeout:      30 * time.Minute,
			},
		},
		&ddbinventory.DynamoDBInventoryService{
			Client: ddbinventory.New(awsCfg),
			Options: ddbinventory.Options{
				ChunkSize:    100,
				RegionFilter: []string{"us-east-1", "us-west-2"},
				MaxRetries:   3,
				Timeout:      30 * time.Minute,
			},
		},
	}

	results, err := awsCventSvc.RunAllInventories(ctx, services)
	if err != nil {
		log.Printf("Some inventory fetches failed: %v", err)
	}

	for name, inv := range results {
		log.Printf("%s inventory fetched: %T", name, inv)
		if err := storeCloudventory(name, inv); err != nil {
			log.Printf("failed to store inventory for %s: %v", name, err)
		}
	}
}

func storeCloudventory(service string, inventory interface{}) error {
	outputDir := ".cloudventory"
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_inventory.json", service))

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing file: %w", err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(inventory); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	fmt.Printf("%s inventory written to %s\n", service, outputFile)
	return nil
}
