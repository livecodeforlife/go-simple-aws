package provider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/cloud"
	autoscalingautoscalinggroupmanager "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/manager/autoscaling/autoscalinggroup"
	ec2launchtemplatemanager "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/manager/ec2/launchtemplate"
	ec2subnetmanager "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/manager/ec2/subnet"
	ec2vpcmanager "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/manager/ec2/vpc"
	elasticloadbalancingv2loadbalancermanager "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/manager/elasticloadbalacingv2/loadbalancer"
	route53resourcerecodsetmanager "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/manager/route53/resourcerecordset"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
)

// NewResourceProvider returns a new aws provider
func NewResourceProvider(config aws.Config) cloud.ResourceManagerProvider {
	return &provider{config}
}
func (p *provider) VPC() types.VpcResourceManager {
	return ec2vpcmanager.NewFromConfig(p.config)
}
func (p *provider) DNSRecordSet() types.DnsRecordSetResourceManager {
	return route53resourcerecodsetmanager.NewFromConfig(p.config)
}
func (p *provider) Subnet() types.SubnetResourceManager {
	return ec2subnetmanager.NewFromConfig(p.config)
}
func (p *provider) LaunchTemplate() types.LaunchTemplateResourceManager {
	return ec2launchtemplatemanager.NewFromConfig(p.config)
}
func (p *provider) LoadBalancer() types.LoadBalancerResourceManager {
	return elasticloadbalancingv2loadbalancermanager.NewFromConfig(p.config)
}
func (p *provider) AutoScalingGroup() types.AutoScalingGroupResourceManager {
	return autoscalingautoscalinggroupmanager.NewFromConfig(p.config)
}

type provider struct {
	config aws.Config
}
