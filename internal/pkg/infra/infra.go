package infra

import "fmt"

// New initializes a new infrastructure manager with the specified provider.
func New(provider Provider) *Infra {
	return &Infra{provider, make(map[string]interface{})}
}

// Infra manages the lifecycle of cloud infrastructure components, such as creation,
// update, and deletion.
type Infra struct {
	provider  Provider               // Provider implements resource creation interfaces.
	resources map[string]interface{} // Tracks created resources.
}

// Resource encapsulates a cloud resource that has been deployed, holding its ID and a
// provider-specific handler.
type Resource[T any] struct {
	id      string
	handler interface{} // Provider-specific reference to the created resource.
}

// ResourceManager create or update resources
type ResourceManager[T any] interface {
	Create(id string, resourceDef T) (*Resource[T], error)
	Delete(id string) error
}

// Provider aggregates interfaces for creating cloud resources. Implementations of Provider
// enable the creation of VPCs, DNS records, and subnets, along with managing their resource handlers.
type Provider interface {
	VPC() ResourceManager[VPC]
	DNS() ResourceManager[DNS]
	Subnet() ResourceManager[Subnet]
	LoadBalancer() ResourceManager[LoadBalancer]
	LaunchTemplate() ResourceManager[LaunchTemplate]
	AutoScale() ResourceManager[AutoScale]
}

// VPC represents the configuration for a Virtual Private Cloud to be deployed.
type VPC struct {
}

// LoadBalancer is a Load Balancer
type LoadBalancer struct {
	VPC    *Resource[VPC]
	DNS    *Resource[DNS]
	Subnet *Resource[Subnet]
}

// LaunchTemplate is a template to lauch an image
type LaunchTemplate struct {
}

// AutoScale defines an AutoScale resource
type AutoScale struct {
}

// Subnet represents the configuration for a subnet within a VPC.
type Subnet struct {
	VPC *Resource[VPC] // Reference to the VPC resource this subnet belongs to.
}

// DNS represents the configuration for a DNS record.
type DNS struct {
}

// CreateVPC requests the creation of a VPC resource in the cloud, using the provided definition.
func (i *Infra) CreateVPC(id string, resourceDef VPC) (*Resource[VPC], error) {
	return innerCreate[VPC](i, id, resourceDef, i.provider.VPC())
}

// CreateDNS requests the creation of a DNS record in the cloud, using the provided definition.
func (i *Infra) CreateDNS(id string, resourceDef DNS) (*Resource[DNS], error) {
	return innerCreate[DNS](i, id, resourceDef, i.provider.DNS())
}

// CreateSubnet requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *Infra) CreateSubnet(id string, resourceDef Subnet) (*Resource[Subnet], error) {
	return innerCreate[Subnet](i, id, resourceDef, i.provider.Subnet())
}

// CreateLoadBalancer requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *Infra) CreateLoadBalancer(id string, resourceDef LoadBalancer) (*Resource[LoadBalancer], error) {
	return innerCreate[LoadBalancer](i, id, resourceDef, i.provider.LoadBalancer())
}

// CreateLaunchTemplate requests the creation of a LaunchTemplate resource in the cloud, using the provided definition.
func (i *Infra) CreateLaunchTemplate(id string, resourceDef LaunchTemplate) (*Resource[LaunchTemplate], error) {
	return innerCreate[LaunchTemplate](i, id, resourceDef, i.provider.LaunchTemplate())
}

// CreateAutoScale requests the creation of a LaunchTemplate resource in the cloud, using the provided definition.
func (i *Infra) CreateAutoScale(id string, resourceDef AutoScale) (*Resource[AutoScale], error) {
	return innerCreate[AutoScale](i, id, resourceDef, i.provider.AutoScale())
}

// innerCreate is a generic function that encapsulates common logic for resource creation.
// It checks for the presence of a provider and resources, ensuring id uniqueness and provider
// ability to innerCreate and store the resource.
func innerCreate[T any](infra *Infra, id string, resourceDef T, manager ResourceManager[T]) (*Resource[T], error) {
	//Check initialization
	if infra == nil {
		return nil, fmt.Errorf("infra is nil")
	}
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	if infra.provider == nil || infra.resources == nil {
		return nil, fmt.Errorf("missing provider or resources; ensure Infra is initialized with New")
	}
	if manager == nil {
		return nil, fmt.Errorf("resource manager is nil")
	}

	//Check unique call
	if _, exists := infra.resources[id]; exists {
		return nil, fmt.Errorf("resource with id %s already exists", id)
	}

	//Creates the resource. the resource manager should take care of idempotency
	resource, err := manager.Create(id, resourceDef)
	if err != nil {
		return nil, err
	}

	//Check if resource manager returned a valid resource
	if resource == nil {
		return nil, fmt.Errorf("provider returned a nil or invalid resource for id %s", id)
	}
	if resource.handler == nil {
		return nil, fmt.Errorf("provider returned a nil or invalid resource handler for id %s", id)
	}
	if resource.id != id {
		return nil, fmt.Errorf("provider returned a resource with a different id[%s] for id %s", resource.id, id)
	}

	infra.resources[id] = resource

	return resource, nil
}
