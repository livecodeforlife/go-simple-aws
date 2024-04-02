package ec2vpcmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
)

// New Creates a new instsance of the resource manager
func New(client *ec2.Client) awsinfra.ResourceManager[*ec2.CreateVpcInput, *types.Vpc] {
	return &manager{
		client,
	}
}

type manager struct {
	client *ec2.Client
}

func (rm *manager) Create(input *ec2.CreateVpcInput) (awsinfra.ExternalID, *types.Vpc, error) {
	output, err := rm.client.CreateVpc(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.Vpc.VpcId, output.Vpc, nil
}
func (rm *manager) Update(input *ec2.CreateVpcInput, last *types.Vpc) (awsinfra.ExternalID, *types.Vpc, error) {
	return aws.String(""), &types.Vpc{}, fmt.Errorf("TODO: Need to implement")
}
func (rm *manager) Load(id awsinfra.ExternalID) (*types.Vpc, error) {
	output, err := rm.client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{
		VpcIds: []string{string(*id)},
	})
	if err != nil {
		return nil, err
	}
	if len(output.Vpcs) == 0 {
		return nil, fmt.Errorf("VPC %s not found", *id)
	}
	return &output.Vpcs[0], nil
}

func (rm *manager) Destroy(id awsinfra.ExternalID) error {
	return fmt.Errorf("TODO: Need to implement")
}
