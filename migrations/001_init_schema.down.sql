-- Drop all tables in reverse order
DROP TABLE IF EXISTS schema_migrations;
DROP TABLE IF EXISTS dynamodb_tables;
DROP TABLE IF EXISTS lambda_functions;
DROP TABLE IF EXISTS rds_instances;
DROP TABLE IF EXISTS ec2_instances;
DROP TABLE IF EXISTS s3_buckets;