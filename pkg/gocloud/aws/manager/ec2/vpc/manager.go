package ec2vpcmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewFromConfig  Creates a new instance of the resource manager
func NewFromConfig(config aws.Config) types.VpcResourceManager {
	return &manager{
		client: ec2.NewFromConfig(config),
	}
}

type manager struct {
	client *ec2.Client
}

func (rm *manager) Create(input *types.VpcInput) (*types.VpcID, *types.VpcOutput, error) {
	output, err := rm.client.CreateVpc(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	for {
		vpc, err := rm.Retrieve(output.Vpc.VpcId)
		if err != nil {
			return nil, nil, err
		}
		if vpc.State == "available" {
			output.Vpc = vpc
			break
		}
		time.Sleep(time.Duration(time.Second * 5))
	}
	return output.Vpc.VpcId, output.Vpc, nil
}

func (rm *manager) Update(id *types.VpcID, input *types.VpcInput) (*types.VpcID, *types.VpcOutput, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *manager) Retrieve(id *types.VpcID) (*types.VpcOutput, error) {
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

func (rm *manager) Delete(id *types.VpcID) (bool, error) {
	return false, fmt.Errorf("TODO: Need to implement")
}
