package ec2subnetmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewFromConfig Creates a new instsance of the resource manager
func NewFromConfig(config aws.Config) types.SubnetResourceManager {
	return &manager{
		client: ec2.NewFromConfig(config),
	}
}

type manager struct {
	client *ec2.Client
}

func (rm *manager) Create(input *types.SubnetInput) (*types.SubnetID, *types.SubnetOutput, error) {
	output, err := rm.client.CreateSubnet(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.Subnet.SubnetId, output.Subnet, nil
}
func (rm *manager) Update(id *types.SubnetID, input *types.SubnetInput) (*types.SubnetID, *types.SubnetOutput, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}
func (rm *manager) Retrieve(id *types.SubnetID) (*types.SubnetOutput, error) {
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

func (rm *manager) Delete(id *types.SubnetID) (bool, error) {
	return false, fmt.Errorf("TODO: Need to implement")
}
