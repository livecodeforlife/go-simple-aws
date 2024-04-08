package ec2launchtemplatemanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewFromConfig Creates a new instsance of the resource manager
func NewFromConfig(config aws.Config) types.LaunchTemplateResourceManager {
	return &manager{
		client: ec2.NewFromConfig(config),
	}
}

type manager struct {
	client *ec2.Client
}

func (rm *manager) Create(input *types.LaunchTemplateInput) (*types.LaunchTemplateID, *types.LaunchTemplateOutput, error) {
	output, err := rm.client.CreateLaunchTemplate(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.LaunchTemplate.LaunchTemplateId, output.LaunchTemplate, nil
}

func (rm *manager) Update(id *types.LaunchTemplateID, input *types.LaunchTemplateInput) (*types.LaunchTemplateID, *types.LaunchTemplateOutput, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *manager) Retrieve(id *types.LaunchTemplateID) (*types.LaunchTemplateOutput, error) {
	output, err := rm.client.DescribeLaunchTemplates(context.TODO(), &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{*id},
		MaxResults:        aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(output.LaunchTemplates) == 0 {
		return nil, fmt.Errorf("Launch Template %s not found", *id)
	}
	return &output.LaunchTemplates[0], nil
}

func (rm *manager) Delete(id *types.LaunchTemplateID) (bool, error) {
	return false, fmt.Errorf("TODO: Need to implement")
}
