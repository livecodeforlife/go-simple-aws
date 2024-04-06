package awsinfra

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
	"github.com/livecodeforlife/go-simple-aws/pkg/coreinfra"
)

// String  Return an aws string
func String(s string) *string {
	return aws.String(s)
}

// VpcResource is an alias for the Vpc type
type VpcResource = coreinfra.Resource[*ec2.CreateVpcInput, *ec2types.Vpc, *string]

// AwsInfra is the main struct for the awsinfra package. It contains the resourceManagerProvider, resourceStorer, and resourcePlanner.
type AwsInfra struct {
	resourceManagerProvider ResourceManagerProvider
	resourceStorer          coreinfra.ResourceStorer
	resourcePlanner         coreinfra.Planner
}

// ResourceStorer retrieve the resource storer being used
func (i *AwsInfra) ResourceStorer() coreinfra.ResourceStorer {
	return i.resourceStorer
}

// NewWithDefaults creates a new AwsInfra struct with default values for resourceManagerProvider, resourceStorer, and resourcePlanner.
func NewWithDefaults() *AwsInfra {
	return &AwsInfra{
		resourceManagerProvider: nil, //TODO: Create a resource manager provider
		resourceStorer:          nil, //TODO: Create a resource store
		resourcePlanner:         coreinfra.NewSimplePlanner(),
	}
}

// Apply creates all resources added to the plan
func (i *AwsInfra) Apply() error {
	return coreinfra.ApplyPlan(i.resourcePlanner)
}

// Destroy deletes all resources added to the plan
func (i *AwsInfra) Destroy() error {
	return coreinfra.DestroyPlan(i.resourcePlanner)
}

type AwsID = *string
type CreateVpcInput = ec2.CreateVpcInput
type Vpc = ec2types.Vpc
type ChangeResourceRecordSetsInput = route53.ChangeResourceRecordSetsInput
type ChangeInfo = route53types.ChangeInfo
type CreateSubnetInput = ec2.CreateSubnetInput
type Subnet = ec2types.Subnet
type CreateLoadBalancerInput = elbv2.CreateLoadBalancerInput
type LoadBalancer = elbv2types.LoadBalancer
type CreateLaunchTemplateInput = ec2.CreateLaunchTemplateInput
type LaunchTemplate = ec2types.LaunchTemplate
type CreateAutoScalingGroupInput = autoscaling.CreateAutoScalingGroupInput
type AutoScalingGroup = autoscalingtypes.AutoScalingGroup

// New creates a new AwsInfra struct.
func New(resourceManagerProvider ResourceManagerProvider, resourceStorer coreinfra.ResourceStorer, resourcePlanner coreinfra.Planner) *AwsInfra {
	return &AwsInfra{
		resourceManagerProvider: resourceManagerProvider,
		resourceStorer:          resourceStorer,
		resourcePlanner:         resourcePlanner,
	}
}

// ResourceManagerProvider aggregates interfaces for creating cloud resources. Implementations of ResourceManagerProvider
// enable the creation of VPCs, DNS records, and subnets, along with managing their resource handlers.
type ResourceManagerProvider interface {
	VPC() coreinfra.ResourceManager[*CreateVpcInput, *Vpc, AwsID]
	DNSRecordSet() coreinfra.ResourceManager[*ChangeResourceRecordSetsInput, *ChangeInfo, AwsID]
	Subnet() coreinfra.ResourceManager[*CreateSubnetInput, *Subnet, AwsID]
	LoadBalancer() coreinfra.ResourceManager[*CreateLoadBalancerInput, []LoadBalancer, AwsID]
	LaunchTemplate() coreinfra.ResourceManager[*CreateLaunchTemplateInput, *LaunchTemplate, AwsID]
	AutoScalingGroup() coreinfra.ResourceManager[*CreateAutoScalingGroupInput, *AutoScalingGroup, AwsID]
}

// CreateVPC requests the creation of a VPC resource in the cloud, using the provided definition.
func (i *AwsInfra) CreateVPC(id coreinfra.ID, input *CreateVpcInput) (*coreinfra.LazyResource[*CreateVpcInput], error) {
	return coreinfra.CreateResource(i.resourcePlanner, i.resourceStorer, i.resourceManagerProvider.VPC(), id, input)
}

// CreateDNS requests the creation of a DNS record in the cloud, using the provided definition.
func (i *AwsInfra) CreateDNS(id coreinfra.ID, input *ChangeResourceRecordSetsInput) (*coreinfra.LazyResource[*ChangeResourceRecordSetsInput], error) {
	return coreinfra.CreateResource(i.resourcePlanner, i.resourceStorer, i.resourceManagerProvider.DNSRecordSet(), id, input)
}

// CreateSubnet requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *AwsInfra) CreateSubnet(id coreinfra.ID, input *CreateSubnetInput) (*coreinfra.LazyResource[*CreateSubnetInput], error) {
	return coreinfra.CreateResource(i.resourcePlanner, i.resourceStorer, i.resourceManagerProvider.Subnet(), id, input)
}

// CreateLoadBalancer requests the creation of a LoadBalancer resource in the cloud, using the provided definition.
func (i *AwsInfra) CreateLoadBalancer(id coreinfra.ID, input *CreateLoadBalancerInput) (*coreinfra.LazyResource[*CreateLoadBalancerInput], error) {
	return coreinfra.CreateResource(i.resourcePlanner, i.resourceStorer, i.resourceManagerProvider.LoadBalancer(), id, input)
}

// CreateLaunchTemplate requests the creation of a LaunchTemplate resource in the cloud, using the provided definition.
func (i *AwsInfra) CreateLaunchTemplate(id coreinfra.ID, input *CreateLaunchTemplateInput) (*coreinfra.LazyResource[*CreateLaunchTemplateInput], error) {
	return coreinfra.CreateResource(i.resourcePlanner, i.resourceStorer, i.resourceManagerProvider.LaunchTemplate(), id, input)
}

// CreateAutoScalingGroup requests the creation of an AutoScalingGroup resource in the cloud, using the provided definition.
func (i *AwsInfra) CreateAutoScalingGroup(id coreinfra.ID, input *CreateAutoScalingGroupInput) (*coreinfra.LazyResource[*CreateAutoScalingGroupInput], error) {
	return coreinfra.CreateResource(i.resourcePlanner, i.resourceStorer, i.resourceManagerProvider.AutoScalingGroup(), id, input)
}
