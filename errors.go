package ec2inventory

import "fmt"

func ErrListBuckets(err error) error {
	return fmt.Errorf("failed to list buckets: %w", err)
}

func ErrGetRegion(err error) error {
	return fmt.Errorf("failed to get region: %w", err)
}
