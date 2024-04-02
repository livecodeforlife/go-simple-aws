package provider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"

	"github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra"
	autoscalingautoscalinggroupmanager "github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra/managers/autoscaling/autoscalinggroup"
	ec2launchtemplatemanager "github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra/managers/ec2/launchtemplate"
	ec2subnetmanager "github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra/managers/ec2/subnet"
	ec2vpcmanager "github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra/managers/ec2/vpc"
	elasticloadbalancingv2loadbalancermanager "github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra/managers/elasticloadbalacingv2/loadbalancer"
	route53resourcerecodsetmanager "github.com/livecodeforlife/go-simple-aws/internal/pkg/awsinfra/managers/route53/resourcerecordset"
)

// NewResourceProvider returns a new aws provider
func NewResourceProvider(config aws.Config) awsinfra.ResourceProvider {
	return &provider{config}
}
func (p *provider) VPC() awsinfra.ResourceManager[*ec2.CreateVpcInput, *ec2types.Vpc] {
	return ec2vpcmanager.New(ec2.NewFromConfig(p.config))
}
func (p *provider) DNSRecordSet() awsinfra.ResourceManager[*route53.ChangeResourceRecordSetsInput, *route53types.ChangeInfo] {
	return route53resourcerecodsetmanager.New(route53.NewFromConfig(p.config))
}
func (p *provider) Subnet() awsinfra.ResourceManager[*ec2.CreateSubnetInput, *ec2types.Subnet] {
	return ec2subnetmanager.New(ec2.NewFromConfig(p.config))
}
func (p *provider) LaunchTemplate() awsinfra.ResourceManager[*ec2.CreateLaunchTemplateInput, *ec2types.LaunchTemplate] {
	return ec2launchtemplatemanager.New(ec2.NewFromConfig(p.config))
}
func (p *provider) LoadBalancer() awsinfra.ResourceManager[*elbv2.CreateLoadBalancerInput, []elbv2types.LoadBalancer] {
	return elasticloadbalancingv2loadbalancermanager.New(elbv2.NewFromConfig(p.config))
}
func (p *provider) AutoScalingGroup() awsinfra.ResourceManager[*autoscaling.CreateAutoScalingGroupInput, *autoscalingtypes.AutoScalingGroup] {
	return autoscalingautoscalinggroupmanager.New(autoscaling.NewFromConfig(p.config))
}

type provider struct {
	config aws.Config
}
