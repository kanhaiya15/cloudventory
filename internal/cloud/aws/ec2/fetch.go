package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Resource struct {
	InstanceID       string `json:"instance_id"`
	InstanceType     string `json:"instance_type"`
	State            string `json:"state"`
	AvailabilityZone string `json:"availability_zone"`
	PublicIP         string `json:"public_ip"`
	PrivateIP        string `json:"private_ip"`
	LaunchTime       string `json:"launch_time"`
	ImageID          string `json:"image_id"`
	VpcID            string `json:"vpc_id"`
	SubnetID         string `json:"subnet_id"`
	SecurityGroups   string `json:"security_groups"`
	Tags             string `json:"tags"`
}

func FetchResources(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := ec2.NewFromConfig(cfg)
	
	result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe EC2 instances: %w", err)
	}

	var resources []Resource
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			resource := Resource{
				InstanceID:       aws.ToString(instance.InstanceId),
				InstanceType:     string(instance.InstanceType),
				State:            string(instance.State.Name),
				AvailabilityZone: aws.ToString(instance.Placement.AvailabilityZone),
				PublicIP:         aws.ToString(instance.PublicIpAddress),
				PrivateIP:        aws.ToString(instance.PrivateIpAddress),
				ImageID:          aws.ToString(instance.ImageId),
				VpcID:            aws.ToString(instance.VpcId),
				SubnetID:         aws.ToString(instance.SubnetId),
			}

			if instance.LaunchTime != nil {
				resource.LaunchTime = instance.LaunchTime.String()
			}

			// Extract security groups
			var sgNames []string
			for _, sg := range instance.SecurityGroups {
				sgNames = append(sgNames, aws.ToString(sg.GroupName))
			}
			resource.SecurityGroups = fmt.Sprintf("[%s]", joinStrings(sgNames, ","))

			// Extract tags
			var tagPairs []string
			for _, tag := range instance.Tags {
				tagPairs = append(tagPairs, fmt.Sprintf("%s=%s", 
					aws.ToString(tag.Key), aws.ToString(tag.Value)))
			}
			resource.Tags = fmt.Sprintf("{%s}", joinStrings(tagPairs, ","))

			resources = append(resources, resource)
		}
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