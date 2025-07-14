package main

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		setEnv       bool
		envValue     string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			setEnv:       true,
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "Environment variable does not exist",
			key:          "NONEXISTENT_VAR",
			defaultValue: "default",
			setEnv:       false,
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.key, tt.envValue)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	config := &Config{
		DatabaseURL: "postgres://test",
		AWSRegion:   "us-west-2",
		Parallel:    true,
	}

	if config.DatabaseURL != "postgres://test" {
		t.Errorf("Expected DatabaseURL to be 'postgres://test', got %s", config.DatabaseURL)
	}

	if config.AWSRegion != "us-west-2" {
		t.Errorf("Expected AWSRegion to be 'us-west-2', got %s", config.AWSRegion)
	}

	if !config.Parallel {
		t.Errorf("Expected Parallel to be true, got %v", config.Parallel)
	}
}