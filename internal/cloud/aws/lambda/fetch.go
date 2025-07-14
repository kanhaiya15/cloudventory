package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type Resource struct {
	FunctionName     string `json:"function_name"`
	FunctionArn      string `json:"function_arn"`
	Runtime          string `json:"runtime"`
	Role             string `json:"role"`
	Handler          string `json:"handler"`
	CodeSize         int64  `json:"code_size"`
	Description      string `json:"description"`
	Timeout          int32  `json:"timeout"`
	MemorySize       int32  `json:"memory_size"`
	LastModified     string `json:"last_modified"`
	CodeSha256       string `json:"code_sha256"`
	Version          string `json:"version"`
	Environment      string `json:"environment"`
	VpcConfig        string `json:"vpc_config"`
	DeadLetterConfig string `json:"dead_letter_config"`
	State            string `json:"state"`
	StateReason      string `json:"state_reason"`
}

func FetchResources(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := lambda.NewFromConfig(cfg)
	
	result, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Lambda functions: %w", err)
	}

	var resources []Resource
	for _, function := range result.Functions {
		resource := Resource{
			FunctionName: aws.ToString(function.FunctionName),
			FunctionArn:  aws.ToString(function.FunctionArn),
			Runtime:      string(function.Runtime),
			Role:         aws.ToString(function.Role),
			Handler:      aws.ToString(function.Handler),
			CodeSize:     function.CodeSize,
			Description:  aws.ToString(function.Description),
			Timeout:      aws.ToInt32(function.Timeout),
			MemorySize:   aws.ToInt32(function.MemorySize),
			LastModified: aws.ToString(function.LastModified),
			CodeSha256:   aws.ToString(function.CodeSha256),
			Version:      aws.ToString(function.Version),
			State:        string(function.State),
			StateReason:  aws.ToString(function.StateReason),
		}

		// Extract environment variables
		if function.Environment != nil && len(function.Environment.Variables) > 0 {
			var envVars []string
			for key, value := range function.Environment.Variables {
				envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
			}
			resource.Environment = fmt.Sprintf("{%s}", joinStrings(envVars, ","))
		}

		// Extract VPC configuration
		if function.VpcConfig != nil {
			var vpcInfo []string
			if function.VpcConfig.VpcId != nil {
				vpcInfo = append(vpcInfo, fmt.Sprintf("VpcId=%s", aws.ToString(function.VpcConfig.VpcId)))
			}
			if len(function.VpcConfig.SubnetIds) > 0 {
				vpcInfo = append(vpcInfo, fmt.Sprintf("Subnets=[%s]", joinStrings(function.VpcConfig.SubnetIds, ",")))
			}
			if len(function.VpcConfig.SecurityGroupIds) > 0 {
				vpcInfo = append(vpcInfo, fmt.Sprintf("SecurityGroups=[%s]", joinStrings(function.VpcConfig.SecurityGroupIds, ",")))
			}
			resource.VpcConfig = fmt.Sprintf("{%s}", joinStrings(vpcInfo, ","))
		}

		// Extract dead letter configuration
		if function.DeadLetterConfig != nil && function.DeadLetterConfig.TargetArn != nil {
			resource.DeadLetterConfig = aws.ToString(function.DeadLetterConfig.TargetArn)
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