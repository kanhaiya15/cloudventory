package ec2inventory

import (
	"context"
	"errors"
	"fmt"
	"kanhaiya1501/cloudventory/aws/cfg"
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2InventoryClient is used to fetch EC2 instances
type EC2InventoryClient struct {
	awsCfg aws.Config
}

// New returns a new EC2InventoryClient
func New(cfg aws.Config) *EC2InventoryClient {
	return &EC2InventoryClient{awsCfg: cfg}
}

// Options controls how EC2 inventory is fetched
type Options struct {
	ChunkSize          int
	RegionFilter       []string
	MaxRetries         int
	EnableParallelScan bool
	Timeout            time.Duration
}

func (o Options) Validate() error {
	if o.ChunkSize <= 0 {
		return errors.New("chunk size must be greater than 0")
	}
	if o.MaxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}
	if o.Timeout <= 0 {
		return errors.New("timeout must be set")
	}
	return nil
}

// FetchInventoryAcrossRegions gathers EC2 instances across filtered regions
func (client *EC2InventoryClient) FetchInventoryAcrossRegions(ctx context.Context, opts Options) ([]EC2Inventory, error) {
	start := time.Now()
	defer func(t time.Time) {
		log.Printf("EC2 inventory fetch duration: %s", time.Since(t))
	}(start)

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	outCallerIdentity, err := cfg.GetCallerIdentity(ctx, client.awsCfg)
	if err != nil {
		return nil, err
	}

	outListRegions, err := cfg.ListRegions(ctx, client.awsCfg)
	if err != nil {
		return nil, err
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []EC2Inventory
		sem     = make(chan struct{}, opts.ChunkSize)
	)

	for _, regionEntry := range outListRegions.Regions {
		region := aws.ToString(regionEntry.RegionName)

		if len(opts.RegionFilter) > 0 && !slices.Contains(opts.RegionFilter, region) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(region string) {
			defer wg.Done()
			defer func() { <-sem }()

			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("panic in region %s: %v", region, r)
				}
			}()

			regionCfg := client.awsCfg.Copy()
			regionCfg.Region = region
			ec2Client := ec2.NewFromConfig(regionCfg)

			paginator := ec2.NewDescribeInstancesPaginator(ec2Client, &ec2.DescribeInstancesInput{})

			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					fmt.Printf("failed to fetch instances in %s: %v", region, err)
					break
				}
				for _, res := range page.Reservations {
					for _, inst := range res.Instances {
						inventory := EC2Inventory{
							Account:    *outCallerIdentity.Account,
							AccountArn: *outCallerIdentity.Arn,
							Region:     region,
						}
						inventory.InstanceID = aws.ToString(inst.InstanceId)
						inventory.Arn = fmt.Sprintf("arn:aws:ec2:%s:%s:instance/%s", region, inventory.Account, inventory.InstanceID)
						inventory.Name = getTagValue(inst.Tags, "Name")
						inventory.State = string(inst.State.Name)
						inventory.LaunchTime = aws.ToTime(inst.LaunchTime)
						inventory.Region = region
						inventory.AvailabilityZone = aws.ToString(inst.Placement.AvailabilityZone)
						inventory.InstanceType = string(inst.InstanceType)
						inventory.Monitoring = string(inst.Monitoring.State)
						inventory.IAMRole = extractIAMRoleName(inst)
						inventory.KeyName = aws.ToString(inst.KeyName)
						inventory.SubnetID = aws.ToString(inst.SubnetId)
						inventory.VpcID = aws.ToString(inst.VpcId)
						inventory.PublicIP = aws.ToString(inst.PublicIpAddress)
						inventory.PublicDNSName = aws.ToString(inst.PublicDnsName)
						inventory.PrivateIP = aws.ToString(inst.PrivateIpAddress)
						inventory.PrivateDNSName = aws.ToString(inst.PrivateDnsName)
						inventory.NetworkInterfaces = inst.NetworkInterfaces
						inventory.SecurityGroups = inst.SecurityGroups
						inventory.Tags = extractTags(inst.Tags)

						mu.Lock()
						results = append(results, inventory)
						mu.Unlock()
					}
				}

			}
		}(region)
	}

	wg.Wait()
	return results, nil
}

func extractIAMRoleName(instance types.Instance) string {
	if instance.IamInstanceProfile == nil || instance.IamInstanceProfile.Arn == nil {
		return ""
	}
	arn := aws.ToString(instance.IamInstanceProfile.Arn)
	parts := strings.Split(arn, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}

func extractTags(tags []types.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[aws.ToString(t.Key)] = aws.ToString(t.Value)
	}
	return result
}

func getTagValue(tags []types.Tag, key string) string {
	for _, t := range tags {
		if aws.ToString(t.Key) == key {
			return aws.ToString(t.Value)
		}
	}
	return ""
}
