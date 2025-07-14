package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Resource struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	CreationDate string `json:"creation_date"`
	Versioning   string `json:"versioning"`
	Encryption   string `json:"encryption"`
	PublicAccess string `json:"public_access"`
}

func FetchResources(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := s3.NewFromConfig(cfg)
	
	// List all buckets
	result, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 buckets: %w", err)
	}

	var resources []Resource
	for _, bucket := range result.Buckets {
		resource := Resource{
			Name:         aws.ToString(bucket.Name),
			CreationDate: bucket.CreationDate.String(),
		}

		// Get bucket location
		locationResult, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err == nil {
			if locationResult.LocationConstraint == "" {
				resource.Region = "us-east-1" // Default region
			} else {
				resource.Region = string(locationResult.LocationConstraint)
			}
		}

		// Get bucket versioning
		versioningResult, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
			Bucket: bucket.Name,
		})
		if err == nil {
			if versioningResult.Status == types.BucketVersioningStatusEnabled {
				resource.Versioning = "Enabled"
			} else {
				resource.Versioning = "Suspended"
			}
		}

		// Get bucket encryption
		encryptionResult, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
			Bucket: bucket.Name,
		})
		if err == nil && len(encryptionResult.ServerSideEncryptionConfiguration.Rules) > 0 {
			rule := encryptionResult.ServerSideEncryptionConfiguration.Rules[0]
			resource.Encryption = string(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		} else {
			resource.Encryption = "None"
		}

		// Get public access block
		publicAccessResult, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
			Bucket: bucket.Name,
		})
		if err == nil {
			config := publicAccessResult.PublicAccessBlockConfiguration
			if aws.ToBool(config.BlockPublicAcls) && aws.ToBool(config.BlockPublicPolicy) &&
				aws.ToBool(config.IgnorePublicAcls) && aws.ToBool(config.RestrictPublicBuckets) {
				resource.PublicAccess = "Blocked"
			} else {
				resource.PublicAccess = "Allowed"
			}
		} else {
			resource.PublicAccess = "Unknown"
		}

		resources = append(resources, resource)
	}

	return resources, nil
}