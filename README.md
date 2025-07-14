# AWS Inventory System

A comprehensive Go-based system for collecting and storing AWS cloud resource information in PostgreSQL. Supports parallel/sequential orchestration and includes migration tooling.

## Features

- **Multi-service AWS inventory**: S3, EC2, RDS, Lambda, DynamoDB (extensible to Azure, GCP, etc.)
- **Normalized database storage**: Inserts resource data into PostgreSQL with proper indexing
- **Flexible orchestration**: Parallel or sequential execution via runner
- **Database migrations**: Automatic schema management with `make migrate`
- **Dockerized deployment**: Ready for local, development, and production environments
- **Comprehensive testing**: Unit and integration tests included
- **Production-ready**: Multi-stage Docker builds, health checks, and security best practices

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (local or Docker)
- AWS credentials (profile or access keys)

### Installation

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd aws-inventory-system
   make setup
   ```

2. **Configure AWS credentials**:
   ```bash
   # Option 1: AWS CLI
   aws configure
   
   # Option 2: Environment variables
   export AWS_ACCESS_KEY_ID="your-access-key"
   export AWS_SECRET_ACCESS_KEY="your-secret-key"
   export AWS_REGION="us-east-1"
   ```

3. **Run the application**:
   ```bash
   make run
   ```

## Directory Structure

```
aws-inventory-system/
├── cmd/
│   ├── main.go                    # Application entry point
│   ├── main_test.go              # Main tests
│   └── main_integration_test.go  # Integration tests
├── internal/cloud/aws/           # AWS service implementations
│   ├── s3/
│   │   ├── fetch.go              # S3 resource fetching
│   │   ├── insert.go             # S3 database insertion
│   │   ├── fetch_test.go         # S3 fetch tests
│   │   ├── insert_test.go        # S3 insert tests
│   │   └── integration_test.go   # S3 integration tests
│   ├── ec2/                      # EC2 implementation
│   ├── rds/                      # RDS implementation
│   ├── lambda/                   # Lambda implementation
│   └── dynamodb/                 # DynamoDB implementation
├── migrations/
│   ├── 001_init_schema.up.sql    # Database schema
│   └── 001_init_schema.down.sql  # Schema rollback
├── docs/
│   └── api.yaml                  # OpenAPI specification
├── Dockerfile                    # Multi-stage Docker build
├── docker-compose.yml           # Local development environment
├── Makefile                      # Build and deployment automation
├── .env.example                  # Environment configuration template
└── README.md                     # This file
```

## Usage

### Command Line Options

```bash
./aws-inventory-system [flags]

Flags:
  -db-url string
        Database URL (default "postgres://user:password@localhost/aws_inventory?sslmode=disable")
  -region string
        AWS Region (default "us-east-1")
  -parallel
        Run services in parallel (default true)
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/aws_inventory?sslmode=disable

# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key_here
AWS_SECRET_ACCESS_KEY=your_secret_key_here

# Application
PARALLEL_EXECUTION=true
```

### Examples

**Run with specific region and sequential execution**:
```bash
go run ./cmd -region=us-west-2 -parallel=false
```

**Run with custom database**:
```bash
go run ./cmd -db-url="postgres://user:pass@localhost:5432/inventory?sslmode=disable"
```

**Docker deployment**:
```bash
# Development environment
make docker-dev

# Production build
make docker-build
make docker-run
```

## Development

### Makefile Commands

```bash
make help              # Show all available commands
make setup             # Complete development environment setup
make build             # Build the binary
make run               # Run locally
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make lint              # Run linter
make migrate           # Run database migrations
make docker-build      # Build Docker image
make docker-dev        # Start development environment
make clean             # Clean build artifacts
```

### Database Schema

The system creates the following tables:

- `s3_buckets` - S3 bucket information
- `ec2_instances` - EC2 instance details
- `rds_instances` - RDS database instances
- `lambda_functions` - Lambda function metadata
- `dynamodb_tables` - DynamoDB table information
- `schema_migrations` - Migration tracking

Each table includes:
- Unique constraints on resource identifiers
- Proper indexing for performance
- `last_updated` and `created_at` timestamps
- Upsert functionality for incremental updates

### Adding New Services

1. **Create service directory**:
   ```bash
   mkdir -p internal/cloud/aws/newservice
   ```

2. **Implement fetch.go**:
   ```go
   package newservice
   
   type Resource struct {
       // Define resource structure
   }
   
   func FetchResources(ctx context.Context, cfg aws.Config) ([]Resource, error) {
       // Implement AWS API calls
   }
   ```

3. **Implement insert.go**:
   ```go
   func InsertResources(ctx context.Context, db *sql.DB, resources []Resource) error {
       // Implement database insertion with upserts
   }
   ```

4. **Add to main.go**:
   ```go
   func (r *ServiceRunner) collectNewService(ctx context.Context) error {
       resources, err := newservice.FetchResources(ctx, r.awsConfig)
       if err != nil {
           return err
       }
       return newservice.InsertResources(ctx, r.db, resources)
   }
   ```

5. **Add migration**:
   ```bash
   make migrate-create name=add_newservice_table
   ```

### Testing

**Unit tests**:
```bash
make test
```

**Integration tests** (requires AWS credentials):
```bash
make test-integration
```

**Coverage report**:
```bash
make test-coverage
open coverage.html
```

## Deployment

### Local Development

```bash
# Setup environment
make dev-setup

# Run application
make run
```

### Docker Production

```bash
# Build production image
make docker-build

# Deploy with docker-compose
docker-compose up -d
```

### Cloud Deployment

The application is designed to run in containerized environments:

- **ECS/Fargate**: Use the provided Dockerfile
- **Kubernetes**: Create deployments with the Docker image
- **Lambda**: Modify main.go for Lambda handler

## Security Considerations

- **Non-root container**: Runs as user `appuser` (UID 1001)
- **Minimal image**: Based on Alpine Linux
- **Credential management**: Supports IAM roles, AWS profiles, and environment variables
- **Network security**: No exposed ports by default
- **Database security**: Uses connection pooling and prepared statements

## Performance

### Optimizations

- **Parallel execution**: Concurrent AWS API calls
- **Database indexing**: Optimized queries on resource identifiers
- **Connection pooling**: Efficient database connections
- **Prepared statements**: SQL injection prevention and performance

### Scaling Considerations

- **Horizontal scaling**: Multiple instances can run concurrently
- **Database partitioning**: Consider partitioning large tables by region/date
- **Rate limiting**: AWS API rate limits are handled gracefully
- **Memory management**: Configurable connection pools

## Monitoring

### Logs

The application provides structured logging:

```bash
# View logs in Docker
docker-compose logs aws-inventory

# Follow logs
docker-compose logs -f aws-inventory
```

### Metrics

Key metrics to monitor:

- **Execution time**: Total and per-service collection time
- **Resource counts**: Number of resources collected per service
- **Error rates**: Failed API calls or database operations
- **Database performance**: Query execution times

## Troubleshooting

### Common Issues

**AWS Credentials**:
```bash
# Check credentials
make check-aws

# Verify AWS access
aws sts get-caller-identity
```

**Database Connection**:
```bash
# Test database connection
psql "postgres://postgres:password@localhost:5432/aws_inventory"

# Run migrations manually
make migrate
```

**Docker Issues**:
```bash
# Restart services
make docker-stop
make docker-dev

# Clean Docker environment
make docker-clean
```

### Debug Mode

Enable debug logging:
```bash
export LOG_LEVEL=debug
make run
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Run `make test lint`
5. Submit a pull request

## License

[Add your license here]

## Support

For issues and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review the API documentation in `docs/api.yaml`
