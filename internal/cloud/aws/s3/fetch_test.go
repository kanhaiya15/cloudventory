package s3

import (
	"testing"
)

func TestResourceStruct(t *testing.T) {
	resource := Resource{
		Name:         "test-bucket",
		Region:       "us-east-1",
		CreationDate: "2023-01-01",
		Versioning:   "Enabled",
		Encryption:   "AES256",
		PublicAccess: "Blocked",
	}

	if resource.Name != "test-bucket" {
		t.Errorf("Expected Name to be 'test-bucket', got %s", resource.Name)
	}

	if resource.Region != "us-east-1" {
		t.Errorf("Expected Region to be 'us-east-1', got %s", resource.Region)
	}

	if resource.Versioning != "Enabled" {
		t.Errorf("Expected Versioning to be 'Enabled', got %s", resource.Versioning)
	}

	if resource.Encryption != "AES256" {
		t.Errorf("Expected Encryption to be 'AES256', got %s", resource.Encryption)
	}

	if resource.PublicAccess != "Blocked" {
		t.Errorf("Expected PublicAccess to be 'Blocked', got %s", resource.PublicAccess)
	}
}