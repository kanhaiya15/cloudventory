-- S3 Buckets Table
CREATE TABLE IF NOT EXISTS s3_buckets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    region VARCHAR(50),
    creation_date TEXT,
    versioning VARCHAR(20),
    encryption VARCHAR(50),
    public_access VARCHAR(20),
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_s3_buckets_name ON s3_buckets(name);
CREATE INDEX IF NOT EXISTS idx_s3_buckets_region ON s3_buckets(region);
CREATE INDEX IF NOT EXISTS idx_s3_buckets_last_updated ON s3_buckets(last_updated);

-- EC2 Instances Table
CREATE TABLE IF NOT EXISTS ec2_instances (
    id SERIAL PRIMARY KEY,
    instance_id VARCHAR(50) UNIQUE NOT NULL,
    instance_type VARCHAR(50),
    state VARCHAR(20),
    availability_zone VARCHAR(50),
    public_ip INET,
    private_ip INET,
    launch_time TEXT,
    image_id VARCHAR(50),
    vpc_id VARCHAR(50),
    subnet_id VARCHAR(50),
    security_groups TEXT,
    tags TEXT,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ec2_instances_instance_id ON ec2_instances(instance_id);
CREATE INDEX IF NOT EXISTS idx_ec2_instances_state ON ec2_instances(state);
CREATE INDEX IF NOT EXISTS idx_ec2_instances_instance_type ON ec2_instances(instance_type);
CREATE INDEX IF NOT EXISTS idx_ec2_instances_availability_zone ON ec2_instances(availability_zone);
CREATE INDEX IF NOT EXISTS idx_ec2_instances_last_updated ON ec2_instances(last_updated);

-- RDS Instances Table
CREATE TABLE IF NOT EXISTS rds_instances (
    id SERIAL PRIMARY KEY,
    db_instance_identifier VARCHAR(255) UNIQUE NOT NULL,
    db_instance_class VARCHAR(50),
    engine VARCHAR(50),
    engine_version VARCHAR(50),
    db_instance_status VARCHAR(50),
    master_username VARCHAR(100),
    db_name VARCHAR(100),
    allocated_storage INTEGER,
    storage_type VARCHAR(50),
    encrypted BOOLEAN,
    availability_zone VARCHAR(50),
    multi_az BOOLEAN,
    vpc_id VARCHAR(50),
    subnet_group VARCHAR(100),
    security_groups TEXT,
    backup_retention_period INTEGER,
    instance_create_time TEXT,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rds_instances_db_instance_identifier ON rds_instances(db_instance_identifier);
CREATE INDEX IF NOT EXISTS idx_rds_instances_engine ON rds_instances(engine);
CREATE INDEX IF NOT EXISTS idx_rds_instances_db_instance_status ON rds_instances(db_instance_status);
CREATE INDEX IF NOT EXISTS idx_rds_instances_availability_zone ON rds_instances(availability_zone);
CREATE INDEX IF NOT EXISTS idx_rds_instances_last_updated ON rds_instances(last_updated);

-- Lambda Functions Table
CREATE TABLE IF NOT EXISTS lambda_functions (
    id SERIAL PRIMARY KEY,
    function_name VARCHAR(255) UNIQUE NOT NULL,
    function_arn TEXT,
    runtime VARCHAR(50),
    role TEXT,
    handler VARCHAR(100),
    code_size BIGINT,
    description TEXT,
    timeout INTEGER,
    memory_size INTEGER,
    last_modified TEXT,
    code_sha256 VARCHAR(100),
    version VARCHAR(50),
    environment TEXT,
    vpc_config TEXT,
    dead_letter_config TEXT,
    state VARCHAR(50),
    state_reason TEXT,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lambda_functions_function_name ON lambda_functions(function_name);
CREATE INDEX IF NOT EXISTS idx_lambda_functions_runtime ON lambda_functions(runtime);
CREATE INDEX IF NOT EXISTS idx_lambda_functions_state ON lambda_functions(state);
CREATE INDEX IF NOT EXISTS idx_lambda_functions_last_updated ON lambda_functions(last_updated);

-- DynamoDB Tables Table
CREATE TABLE IF NOT EXISTS dynamodb_tables (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) UNIQUE NOT NULL,
    table_arn TEXT,
    table_status VARCHAR(50),
    creation_date_time TEXT,
    provisioned_throughput TEXT,
    billing_mode VARCHAR(50),
    item_count BIGINT,
    table_size_bytes BIGINT,
    global_secondary_indexes TEXT,
    local_secondary_indexes TEXT,
    stream_specification TEXT,
    sse_description TEXT,
    point_in_time_recovery TEXT,
    tags TEXT,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dynamodb_tables_table_name ON dynamodb_tables(table_name);
CREATE INDEX IF NOT EXISTS idx_dynamodb_tables_table_status ON dynamodb_tables(table_status);
CREATE INDEX IF NOT EXISTS idx_dynamodb_tables_billing_mode ON dynamodb_tables(billing_mode);
CREATE INDEX IF NOT EXISTS idx_dynamodb_tables_last_updated ON dynamodb_tables(last_updated);

-- Migration tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL DEFAULT FALSE,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert migration version
INSERT INTO schema_migrations (version) VALUES (1) ON CONFLICT (version) DO NOTHING;