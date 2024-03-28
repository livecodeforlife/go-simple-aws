package infra

import "fmt"

// DefinitionStore abstracts the storage of resource definitions.
type DefinitionStore interface {
	mapInterface
}

// Provider aggregates interfaces for creating cloud resources. Implementations of Provider
// enable the creation of VPCs, DNS records, and subnets, along with managing their handlers.
type Provider interface {
	VPC() ResourceManager[VPC]
	VPCCreator
	DNSCreator
	SubnetCreator
	HandlerStore
}

// ResourceManager create or update
type ResourceManager[T any] interface {
	Create(resource T) (*Resource[T], error)
	Update(resource T, old *Resource[T]) (*Resource[T], error)
	Delete(old *Resource[T]) error
}

// HandlerStore abstracts the storage of handlers for cloud resources. Handlers represent
// references or identifiers to resources created in the cloud.
type HandlerStore interface {
	mapInterface
}

// SubnetCreator defines the interface for creating subnets within a VPC.
type SubnetCreator interface {
	CreateSubnet(Subnet) (*Resource[Subnet], error)
}

// VPCCreator defines the interface for creating VPCs.
type VPCCreator interface {
	CreateVPC(VPC) (*Resource[VPC], error)
}

// DNSCreator defines the interface for creating DNS records.
type DNSCreator interface {
	CreateDNS(DNS) (*Resource[DNS], error)
}

// Resource encapsulates a cloud resource that has been deployed, holding its ID and a
// provider-specific handler.
type Resource[T any] struct {
	id      string
	handler interface{} // Provider-specific reference to the created resource.
}

// Infra manages the lifecycle of cloud infrastructure components, such as creation,
// update, and deletion.
type Infra struct {
	provider  Provider               // Provider implements resource creation interfaces.
	resources map[string]interface{} // Tracks created resources.
}

// New initializes a new infrastructure manager with the specified provider.
func New(provider Provider) *Infra {
	return &Infra{provider, make(map[string]interface{})}
}

// VPC represents the configuration for a Virtual Private Cloud to be deployed.
type VPC struct {
}

// Subnet represents the configuration for a subnet within a VPC.
type Subnet struct {
	vpc *Resource[VPC] // Reference to the VPC resource this subnet belongs to.
}

// DNS represents the configuration for a DNS record.
type DNS struct {
}

// CreateVPC requests the creation of a VPC resource in the cloud, using the provided definition.
func (i *Infra) CreateVPC(id string, resourceDef VPC) (*Resource[VPC], error) {
	return apply[VPC](i, id, resourceDef, i.provider.CreateVPC)
}

// CreateSubnet requests the creation of a Subnet resource in the cloud, using the provided definition.
func (i *Infra) CreateSubnet(id string, resourceDef Subnet) (*Resource[Subnet], error) {
	return apply[Subnet](i, id, resourceDef, i.provider.CreateSubnet)
}

// CreateDNS requests the creation of a DNS record in the cloud, using the provided definition.
func (i *Infra) CreateDNS(id string, resourceDef DNS) (*Resource[DNS], error) {
	return apply[DNS](i, id, resourceDef, i.provider.CreateDNS)
}

// apply is a generic function that encapsulates common logic for resource creation.
// It checks for the presence of a provider and resources, ensuring id uniqueness and provider
// ability to create and store the resource.
func apply[T any](infra *Infra, id string, resourceDef T, create func(T) (*Resource[T], error)) (*Resource[T], error) {
	if infra.provider == nil || infra.resources == nil {
		return nil, fmt.Errorf("missing provider or resources; ensure Infra is initialized with New")
	}
	if _, exists := infra.resources[id]; exists {
		return nil, fmt.Errorf("resource with id %s already exists", id)
	}
	if !infra.provider.Exists(id) {
		resource, err := create(resourceDef)
		if err != nil {
			return nil, err
		}
		infra.provider.Set(id, resource)
	}
	handler, ok := infra.provider.Get(id).(*Resource[T])
	if !ok || handler == nil {
		return nil, fmt.Errorf("provider returned a nil or invalid handler for id %s", id)
	}
	infra.resources[id] = resourceDef
	return handler, nil
}

// apply2 is a generic function that encapsulates common logic for resource creation.
// It checks for the presence of a provider and resources, ensuring id uniqueness and provider
// ability to create and store the resource.
func apply2[T any](infra *Infra, id string, resourceDef T, resourceManager ResourceManager[T]) (*Resource[T], error) {
	if infra.provider == nil || infra.resources == nil {
		return nil, fmt.Errorf("missing provider or resources; ensure Infra is initialized with New")
	}
	if _, exists := infra.resources[id]; exists {
		return nil, fmt.Errorf("resource with id %s already exists", id)
	}
	if !infra.provider.Exists(id) {
		resource, err := resourceManager.Create(resourceDef)
		if err != nil {
			return nil, err
		}
		infra.provider.Set(id, resource)
	} else {
		resource, err := resourceManager.Update(resourceDef, infra.provider.Get(id).(*Resource[T]))
		if err != nil {
			return nil, err
		}
		infra.provider.Set(id, resource)
	}
	handler, ok := infra.provider.Get(id).(*Resource[T])
	if !ok || handler == nil {
		return nil, fmt.Errorf("provider returned a nil or invalid handler for id %s", id)
	}
	infra.resources[id] = resourceDef
	return handler, nil
}

// mapInterface defines basic operations for mapping keys to values, utilized by HandlerStore
// and DefinitionStore for managing cloud resource identifiers and configurations.
type mapInterface interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Exists(key string) bool
}
