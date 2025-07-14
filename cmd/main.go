package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"aws-inventory-system/internal/cloud/aws/dynamodb"
	"aws-inventory-system/internal/cloud/aws/ec2"
	"aws-inventory-system/internal/cloud/aws/lambda"
	"aws-inventory-system/internal/cloud/aws/rds"
	"aws-inventory-system/internal/cloud/aws/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseURL string
	AWSRegion   string
	Parallel    bool
}

type ServiceRunner struct {
	db       *sql.DB
	awsConfig aws.Config
	parallel  bool
}

func main() {
	var (
		dbURL    = flag.String("db-url", getEnv("DATABASE_URL", "postgres://user:password@localhost/aws_inventory?sslmode=disable"), "Database URL")
		region   = flag.String("region", getEnv("AWS_REGION", "us-east-1"), "AWS Region")
		parallel = flag.Bool("parallel", true, "Run services in parallel")
	)
	flag.Parse()

	cfg := &Config{
		DatabaseURL: *dbURL,
		AWSRegion:   *region,
		Parallel:    *parallel,
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *Config) error {
	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize AWS config
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	runner := &ServiceRunner{
		db:        db,
		awsConfig: awsConfig,
		parallel:  cfg.Parallel,
	}

	log.Println("Starting AWS inventory collection...")
	start := time.Now()

	if err := runner.collectInventory(context.TODO()); err != nil {
		return fmt.Errorf("failed to collect inventory: %w", err)
	}

	log.Printf("Inventory collection completed in %v", time.Since(start))
	return nil
}

func (r *ServiceRunner) collectInventory(ctx context.Context) error {
	if r.parallel {
		return r.runParallel(ctx)
	}
	return r.runSequential(ctx)
}

func (r *ServiceRunner) runParallel(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 5)

	services := []func(context.Context) error{
		r.collectS3,
		r.collectEC2,
		r.collectRDS,
		r.collectLambda,
		r.collectDynamoDB,
	}

	for _, service := range services {
		wg.Add(1)
		go func(svc func(context.Context) error) {
			defer wg.Done()
			if err := svc(ctx); err != nil {
				errChan <- err
			}
		}(service)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ServiceRunner) runSequential(ctx context.Context) error {
	services := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"S3", r.collectS3},
		{"EC2", r.collectEC2},
		{"RDS", r.collectRDS},
		{"Lambda", r.collectLambda},
		{"DynamoDB", r.collectDynamoDB},
	}

	for _, service := range services {
		log.Printf("Collecting %s inventory...", service.name)
		if err := service.fn(ctx); err != nil {
			return fmt.Errorf("failed to collect %s inventory: %w", service.name, err)
		}
	}

	return nil
}

func (r *ServiceRunner) collectS3(ctx context.Context) error {
	log.Println("Fetching S3 resources...")
	resources, err := s3.FetchResources(ctx, r.awsConfig)
	if err != nil {
		return err
	}

	log.Printf("Inserting %d S3 resources...", len(resources))
	return s3.InsertResources(ctx, r.db, resources)
}

func (r *ServiceRunner) collectEC2(ctx context.Context) error {
	log.Println("Fetching EC2 resources...")
	resources, err := ec2.FetchResources(ctx, r.awsConfig)
	if err != nil {
		return err
	}

	log.Printf("Inserting %d EC2 resources...", len(resources))
	return ec2.InsertResources(ctx, r.db, resources)
}

func (r *ServiceRunner) collectRDS(ctx context.Context) error {
	log.Println("Fetching RDS resources...")
	resources, err := rds.FetchResources(ctx, r.awsConfig)
	if err != nil {
		return err
	}

	log.Printf("Inserting %d RDS resources...", len(resources))
	return rds.InsertResources(ctx, r.db, resources)
}

func (r *ServiceRunner) collectLambda(ctx context.Context) error {
	log.Println("Fetching Lambda resources...")
	resources, err := lambda.FetchResources(ctx, r.awsConfig)
	if err != nil {
		return err
	}

	log.Printf("Inserting %d Lambda resources...", len(resources))
	return lambda.InsertResources(ctx, r.db, resources)
}

func (r *ServiceRunner) collectDynamoDB(ctx context.Context) error {
	log.Println("Fetching DynamoDB resources...")
	resources, err := dynamodb.FetchResources(ctx, r.awsConfig)
	if err != nil {
		return err
	}

	log.Printf("Inserting %d DynamoDB resources...", len(resources))
	return dynamodb.InsertResources(ctx, r.db, resources)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}