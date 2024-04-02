package ec2subnetmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
)

// New Creates a new instsance of the resource manager
func New(client *ec2.Client) awsinfra.ResourceManager[*ec2.CreateSubnetInput, *types.Subnet] {
	return &manager{
		client,
	}
}

type manager struct {
	client *ec2.Client
}

func (rm *manager) Create(input *ec2.CreateSubnetInput) (awsinfra.ExternalID, *types.Subnet, error) {
	output, err := rm.client.CreateSubnet(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.Subnet.SubnetId, output.Subnet, nil
}
func (rm *manager) Update(input *ec2.CreateSubnetInput, last *types.Subnet) (awsinfra.ExternalID, *types.Subnet, error) {
	return aws.String(""), &types.Subnet{}, fmt.Errorf("TODO: Need to implement")
}
func (rm *manager) Load(id awsinfra.ExternalID) (*types.Subnet, error) {
	output, err := rm.client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
		SubnetIds: []string{*id},
	})
	if err != nil {
		return nil, err
	}
	if len(output.Subnets) == 0 {
		return nil, fmt.Errorf("Subnet %s not found", *id)
	}
	return &output.Subnets[0], nil
}

func (rm *manager) Destroy(id awsinfra.ExternalID) error {
	return fmt.Errorf("TODO: Need to implement")
}
