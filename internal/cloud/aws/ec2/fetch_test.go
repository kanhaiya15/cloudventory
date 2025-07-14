package ec2

import (
	"testing"
)

func TestJoinStrings(t *testing.T) {
	tests := []struct {
		name      string
		slice     []string
		separator string
		expected  string
	}{
		{
			name:      "Empty slice",
			slice:     []string{},
			separator: ",",
			expected:  "",
		},
		{
			name:      "Single element",
			slice:     []string{"a"},
			separator: ",",
			expected:  "a",
		},
		{
			name:      "Multiple elements",
			slice:     []string{"a", "b", "c"},
			separator: ",",
			expected:  "a,b,c",
		},
		{
			name:      "Different separator",
			slice:     []string{"x", "y"},
			separator: "|",
			expected:  "x|y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinStrings(tt.slice, tt.separator)
			if result != tt.expected {
				t.Errorf("joinStrings() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResourceStruct(t *testing.T) {
	resource := Resource{
		InstanceID:       "i-1234567890abcdef0",
		InstanceType:     "t3.micro",
		State:            "running",
		AvailabilityZone: "us-east-1a",
		PublicIP:         "1.2.3.4",
		PrivateIP:        "10.0.1.100",
		ImageID:          "ami-0abcdef1234567890",
		VpcID:            "vpc-12345678",
		SubnetID:         "subnet-12345678",
		SecurityGroups:   "[default,web]",
		Tags:             "{Name=test,Environment=dev}",
	}

	if resource.InstanceID != "i-1234567890abcdef0" {
		t.Errorf("Expected InstanceID to be 'i-1234567890abcdef0', got %s", resource.InstanceID)
	}

	if resource.InstanceType != "t3.micro" {
		t.Errorf("Expected InstanceType to be 't3.micro', got %s", resource.InstanceType)
	}

	if resource.State != "running" {
		t.Errorf("Expected State to be 'running', got %s", resource.State)
	}
}