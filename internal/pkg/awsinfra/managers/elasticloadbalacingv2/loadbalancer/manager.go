package elasticloadbalancingv2loadbalancermanager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
)

// New Creates a new instsance of the resource manager
func New(client *elasticloadbalancingv2.Client) awsinfra.ResourceManager[*elasticloadbalancingv2.CreateLoadBalancerInput, []types.LoadBalancer] {
	return &manager{
		client,
	}
}

type manager struct {
	client *elasticloadbalancingv2.Client
}

func (rm *manager) Create(input *elasticloadbalancingv2.CreateLoadBalancerInput) (awsinfra.ExternalID, []types.LoadBalancer, error) {
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

func (rm *manager) Update(input *elasticloadbalancingv2.CreateLoadBalancerInput, last []types.LoadBalancer) (awsinfra.ExternalID, []types.LoadBalancer, error) {
	return nil, nil, fmt.Errorf("TODO: Need to implement")
}

func (rm *manager) Load(id awsinfra.ExternalID) ([]types.LoadBalancer, error) {
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

func (rm *manager) Destroy(id awsinfra.ExternalID) error {
	return fmt.Errorf("TODO: Need to implement")
}
