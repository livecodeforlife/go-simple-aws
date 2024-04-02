package ec2launchtemplatemanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
)

// New Creates a new instsance of the resource manager
func New(client *ec2.Client) awsinfra.ResourceManager[*ec2.CreateLaunchTemplateInput, *types.LaunchTemplate] {
	return &manager{
		client,
	}
}

type manager struct {
	client *ec2.Client
}

func (rm *manager) Create(input *ec2.CreateLaunchTemplateInput) (awsinfra.ExternalID, *types.LaunchTemplate, error) {
	output, err := rm.client.CreateLaunchTemplate(context.TODO(), input)
	if err != nil {
		return aws.String(""), nil, err
	}
	return output.LaunchTemplate.LaunchTemplateId, output.LaunchTemplate, nil
}
func (rm *manager) Update(input *ec2.CreateLaunchTemplateInput, last *types.LaunchTemplate) (awsinfra.ExternalID, *types.LaunchTemplate, error) {
	return aws.String(""), &types.LaunchTemplate{}, fmt.Errorf("TODO: Need to implement")
}
func (rm *manager) Load(id awsinfra.ExternalID) (*types.LaunchTemplate, error) {
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
func (rm *manager) Destroy(id awsinfra.ExternalID) error {
	return fmt.Errorf("TODO: Need to implement")
}
