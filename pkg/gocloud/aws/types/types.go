package types

import (
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elasticloadbalancingv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/core"
)

type VpcInput = ec2.CreateVpcInput
type VpcOutput = ec2types.Vpc
type VpcID = string
type VpcResourceManager = core.ResourceManager[*VpcInput, *VpcOutput, *VpcID]
type VpcResource = core.Resource[*VpcInput, *VpcOutput, *VpcID]
type VpcLazyResource = core.LazyResource[*VpcInput]

type SubnetInput = ec2.CreateSubnetInput
type SubnetOutput = ec2types.Subnet
type SubnetID = string
type SubnetResourceManager = core.ResourceManager[*SubnetInput, *SubnetOutput, *SubnetID]
type SubnetResource = core.Resource[*SubnetInput, *SubnetOutput, *SubnetID]
type SubnetLazyResource = core.LazyResource[*SubnetInput]

type DnsRecordSetInput = route53.ChangeResourceRecordSetsInput
type DnsRecordSetOutput = route53types.ChangeInfo
type DnsRecordSetID = string
type DnsRecordSetResourceManager = core.ResourceManager[*DnsRecordSetInput, *DnsRecordSetOutput, *DnsRecordSetID]
type DnsRecordSetResource = core.Resource[*DnsRecordSetInput, *DnsRecordSetOutput, *DnsRecordSetID]
type DnsRecordSetLazyResource = core.LazyResource[*DnsRecordSetInput]

type AutoScalingGroupInput = autoscaling.CreateAutoScalingGroupInput
type AutoScalingGroupOutput = autoscalingtypes.AutoScalingGroup
type AutoScalingGroupID = string
type AutoScalingGroupResourceManager = core.ResourceManager[*AutoScalingGroupInput, *AutoScalingGroupOutput, *AutoScalingGroupID]
type AutoScalingGroupResource = core.Resource[*AutoScalingGroupInput, *AutoScalingGroupOutput, *AutoScalingGroupID]
type AutoScalingGroupLazyResource = core.LazyResource[*AutoScalingGroupInput]

type LaunchTemplateInput = ec2.CreateLaunchTemplateInput
type LaunchTemplateOutput = ec2types.LaunchTemplate
type LaunchTemplateID = string
type LaunchTemplateResourceManager = core.ResourceManager[*LaunchTemplateInput, *LaunchTemplateOutput, *LaunchTemplateID]
type LaunchTemplateResource = core.Resource[*LaunchTemplateInput, *LaunchTemplateOutput, *LaunchTemplateID]
type LaunchTemplateLazyResource = core.LazyResource[*LaunchTemplateInput]

type LoadBalancerInput = elasticloadbalancingv2.CreateLoadBalancerInput
type LoadBalancerOutput = elasticloadbalancingv2types.LoadBalancer
type LoadBalancerID = string
type LoadBalancerResourceManager = core.ResourceManager[*LoadBalancerInput, []LoadBalancerOutput, *LoadBalancerID]
type LoadBalancerResource = core.Resource[*LoadBalancerInput, []LoadBalancerOutput, *LoadBalancerID]
type LoadBalancerLazyResource = core.LazyResource[*LoadBalancerInput]
