package ec2inventory

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Inventory struct {
	Account           string                           `json:"account"`
	AccountArn        string                           `json:"accountArn"`
	Arn               string                           `json:"arn"`
	InstanceID        string                           `json:"instanceId"`
	Name              string                           `json:"name"`
	State             string                           `json:"state"`
	LaunchTime        time.Time                        `json:"launchTime"`
	Region            string                           `json:"region"`
	AvailabilityZone  string                           `json:"availabilityZone"`
	InstanceType      string                           `json:"instanceType"`
	Monitoring        string                           `json:"monitoring"`
	IAMRole           string                           `json:"iamRole"`
	KeyName           string                           `json:"keyName"`
	SubnetID          string                           `json:"subnetId"`
	VpcID             string                           `json:"vpcId"`
	PublicIP          string                           `json:"publicIp"`
	PublicDNSName     string                           `json:"publicDnsName"`
	PrivateIP         string                           `json:"privateIp"`
	PrivateDNSName    string                           `json:"privateDnsName"`
	NetworkInterfaces []types.InstanceNetworkInterface `json:"networkInterfaces"`
	SecurityGroups    []types.GroupIdentifier          `json:"securityGroups"`
	Tags              map[string]string                `json:"tags"`
	Errors            []string                         `json:"errors"`
}
