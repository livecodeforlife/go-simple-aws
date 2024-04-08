package elasticloadbalancingv2loadbalancermanager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewFromConfig Creates a new instsance of the resource manager
func NewFromConfig(config aws.Config) types.LoadBalancerResourceManager {
	return &manager{
		client: elasticloadbalancingv2.NewFromConfig(config),
	}
}

type manager struct {
	client *elasticloadbalancingv2.Client
}

func (rm *manager) Create(input *types.LoadBalancerInput) (*types.LoadBalancerID, []types.LoadBalancerOutput, error) {
	output, err := rm.client.CreateLoadBalancer(context.TODO(), input)
	if err != nil {
		return nil, nil, err
	}
	//Encode arn ids into the ExternalID string
	arns := make([]string, len(output.LoadBalancers))
	for i, v := range output.LoadBalancers {
		arns[i] = *v.LoadBalancerArn
	}
	loadBalancerArns, err := json.Marshal(arns)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	}

	return aws.String(string(loadBalancerArns)), output.LoadBalancers, nil
}

func (rm *manager) Update(id *types.LoadBalancerID, input *types.LoadBalancerInput) (*types.LoadBalancerID, []types.LoadBalancerOutput, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *manager) Retrieve(id *types.LoadBalancerID) ([]types.LoadBalancerOutput, error) {
	var loadBalancerArns []string
	if err := json.Unmarshal([]byte(*id), &loadBalancerArns); err != nil {
		return nil, err
	}
	output, err := rm.client.DescribeLoadBalancers(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancersInput{
		LoadBalancerArns: loadBalancerArns,
	})
	if err != nil {
		return nil, err
	}
	return output.LoadBalancers, nil
}

func (rm *manager) Delete(id *types.LoadBalancerID) (bool, error) {
	return false, fmt.Errorf("TODO: Need to implement")
}
