package autoscalingautoscalinggroupmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
)

// New Creates a new instsance of the resource manager
func New(client *autoscaling.Client) awsinfra.ResourceManager[*autoscaling.CreateAutoScalingGroupInput, *types.AutoScalingGroup] {
	return &manager{
		client,
	}
}

type manager struct {
	client *autoscaling.Client
}

func (rm *manager) Create(input *autoscaling.CreateAutoScalingGroupInput) (awsinfra.ExternalID, *types.AutoScalingGroup, error) {
	if *input.AutoScalingGroupName == "" {
		return nil, nil, fmt.Errorf("AutoScalingGroupName is required and is used as the external id")
	}
	_, err := rm.client.CreateAutoScalingGroup(context.TODO(), input)
	if err != nil {
		return nil, nil, err
	}
	asg, err := rm.Load(input.AutoScalingGroupName)
	if err != nil {
		return nil, nil, err
	}
	return asg.AutoScalingGroupName, asg, nil
}

func (rm *manager) Update(input *autoscaling.CreateAutoScalingGroupInput, last *types.AutoScalingGroup) (awsinfra.ExternalID, *types.AutoScalingGroup, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}
func (rm *manager) Load(id awsinfra.ExternalID) (*types.AutoScalingGroup, error) {
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

func (rm *manager) Destroy(id awsinfra.ExternalID) error {
	return fmt.Errorf("TODO: Need to implement")
}
