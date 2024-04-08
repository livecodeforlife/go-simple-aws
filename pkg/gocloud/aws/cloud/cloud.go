package cloud

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/aws/types"
	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/core"
)

// String  Return an aws string
func String(s string) *string {
	return aws.String(s)
}

// Cloud is the main struct for the aws package. It contains the resourceManagerProvider, resourceStorer, and resourcePlanner.
type Cloud struct {
	provider ResourceManagerProvider
	store    core.ResourceStorer
	planner  core.Planner
}

// Load resources
func (i *Cloud) Load() error {
	return i.store.Load()
}

// Save resources
func (i *Cloud) Save() error {
	return i.store.Save()
}

// ResourceStorer retrieve the resource storer being used
func (i *Cloud) ResourceStorer() core.ResourceStorer {
	return i.store
}

// New creates a new AwsInfra struct with default values for resourceManagerProvider, resourceStorer, and resourcePlanner.
func New(provider ResourceManagerProvider, store core.ResourceStorer, planner core.Planner) (*Cloud, error) {
	if err := store.Load(); err != nil {
		return nil, err
	}
	return &Cloud{provider, store, planner}, nil
}

// Apply creates all resources added to the plan
func (i *Cloud) Apply() error {
	if err := core.ApplyPlan(i.planner); err != nil {
		return err
	}
	if err := i.store.Save(); err != nil {
		return err
	}
	return nil
}

// Destroy deletes all resources added to the plan
func (i *Cloud) Destroy() error {
	if err := core.DestroyPlan(i.planner); err != nil {
		return err
	}
	if err := i.store.Save(); err != nil {
		return err
	}
	return nil
}

// ResourceManagerProvider aggregates interfaces for creating cloud resources. Implementations of ResourceManagerProvider
// enable the creation of VPCs, DNS records, and subnets, along with managing their resource handlers.
type ResourceManagerProvider interface {
	VPC() types.VpcResourceManager
	DNSRecordSet() types.DnsRecordSetResourceManager
	Subnet() types.SubnetResourceManager
	LoadBalancer() types.LoadBalancerResourceManager
	LaunchTemplate() types.LaunchTemplateResourceManager
	AutoScalingGroup() types.AutoScalingGroupResourceManager
}

// CreateVPC requests the creation of a VPC resource in the cloud, using the provided definition.
func (i *Cloud) CreateVPC(id core.ID, input *types.VpcInput) (*types.VpcLazyResource, error) {
	return core.CreateResource(i.planner, i.store, i.provider.VPC(), id, input)
}

// CreateDNSRecordSet requests the creation of a DNS record in the cloud, using the provided definition.
func (i *Cloud) CreateDNSRecordSet(id core.ID, input *types.DnsRecordSetInput) (*types.DnsRecordSetLazyResource, error) {
	return core.CreateResource(i.planner, i.store, i.provider.DNSRecordSet(), id, input)
}

// CreateSubnet requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *Cloud) CreateSubnet(id core.ID, input *types.SubnetInput) (*types.SubnetLazyResource, error) {
	return core.CreateResource(i.planner, i.store, i.provider.Subnet(), id, input)
}

// SetSubnetVpc sets the VPC ID on the Subnet resource.
func (i *Cloud) SetSubnetVpc(subnet *types.SubnetLazyResource, vpc *types.VpcLazyResource) {
	core.AddDependency(
		i.ResourceStorer(),
		subnet,
		vpc,
		func(subnet *types.SubnetInput, vpc *types.VpcResource) error {
			subnet.VpcId = vpc.Output.VpcId
			return nil
		})
}

// CreateLoadBalancer requests the creation of a LoadBalancer resource in the cloud, using the provided definition.
func (i *Cloud) CreateLoadBalancer(id core.ID, input *types.LoadBalancerInput) (*types.LoadBalancerLazyResource, error) {
	return core.CreateResource(i.planner, i.store, i.provider.LoadBalancer(), id, input)
}

// CreateLaunchTemplate requests the creation of a LaunchTemplate resource in the cloud, using the provided definition.
func (i *Cloud) CreateLaunchTemplate(id core.ID, input *types.LaunchTemplateInput) (*types.LaunchTemplateLazyResource, error) {
	return core.CreateResource(i.planner, i.store, i.provider.LaunchTemplate(), id, input)
}

// CreateAutoScalingGroup requests the creation of an AutoScalingGroup resource in the cloud, using the provided definition.
func (i *Cloud) CreateAutoScalingGroup(id core.ID, input *types.AutoScalingGroupInput) (*types.AutoScalingGroupLazyResource, error) {
	return core.CreateResource(i.planner, i.store, i.provider.AutoScalingGroup(), id, input)
}
