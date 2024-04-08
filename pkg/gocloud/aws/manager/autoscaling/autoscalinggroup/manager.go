package autoscalingautoscalinggroupmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewFromConfig Creates a new instance of the resource manager
func NewFromConfig(config aws.Config) types.AutoScalingGroupResourceManager {
	return &manager{
		client: autoscaling.NewFromConfig(config),
	}
}

type manager struct {
	client *autoscaling.Client
}

func (rm *manager) Create(input *types.AutoScalingGroupInput) (*types.AutoScalingGroupID, *types.AutoScalingGroupOutput, error) {
	if *input.AutoScalingGroupName == "" {
		return nil, nil, fmt.Errorf("AutoScalingGroupName is required and is used as the external id")
	}
	_, err := rm.client.CreateAutoScalingGroup(context.TODO(), input)
	if err != nil {
		return nil, nil, err
	}
	asg, err := rm.Retrieve(input.AutoScalingGroupName)
	if err != nil {
		return nil, nil, err
	}
	return asg.AutoScalingGroupName, asg, nil
}

func (rm *manager) Update(id *types.AutoScalingGroupID, input *types.AutoScalingGroupInput) (*types.AutoScalingGroupID, *types.AutoScalingGroupOutput, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}
func (rm *manager) Retrieve(id *types.AutoScalingGroupID) (*types.AutoScalingGroupOutput, error) {
	output, err := rm.client.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{*id},
	})
	if err != nil {
		return nil, err
	}
	if len(output.AutoScalingGroups) == 0 {
		return nil, fmt.Errorf("AutoScalingGroup with id %s not found", *id)
	}
	return &output.AutoScalingGroups[0], nil
}

func (rm *manager) Delete(id *types.AutoScalingGroupID) (bool, error) {
	return false, fmt.Errorf("TODO: Need to implement")
}
