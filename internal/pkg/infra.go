package infra

import "fmt"

// DefinitionStore store the definition
type DefinitionStore interface {
	mapInterface
}

// Provider implements Creator interfaces
// External struct should implement provider to create cloud components
type Provider interface {
	VPCCreator
	DNSCreator
	SubnetCreator
	HandlerStore
}

// HandlerStore store the handlers
// Handlers is a abstraction about the created component at the cloud
type HandlerStore interface {
	mapInterface
}

// SubnetCreator is the interface that does create the VPC
type SubnetCreator interface {
	CreateSubnet(Subnet) (*Resource[Subnet], error)
}

// VPCCreator is the interface that does create the VPC
type VPCCreator interface {
	CreateVPC(VPC) (*Resource[VPC], error)
}

// DNSCreator is the interface that does create the DNSCreator
type DNSCreator interface {
	CreateDNS(DNS) (*Resource[DNS], error)
}

// Resource represents an infra component that was already deployed
type Resource[T any] struct {
	id      string
	handler interface{} // this should be known by the provider
}

// Infra manages the creation, update and deletion of infra components
type Infra struct {
	provider  Provider
	resources map[string]interface{}
}

// New returns a new infra provider
func New(provider Provider) *Infra {
	return &Infra{provider, make(map[string]interface{})}
}

// VPC defines a VPC to be deployed
type VPC struct {
}

// Subnet defines a subnet in the VPC
type Subnet struct {
	vpc *Resource[VPC]
}

// DNS defines a DNS
type DNS struct {
}

// CreateVPC takes a VPC resource and creates it in the cloud
func (i *Infra) CreateVPC(id string, resourceDef VPC) (*Resource[VPC], error) {
	return createResource[VPC](i, id, resourceDef, i.provider.CreateVPC)
}

// CreateSubnet takes a Subnet resource and creates it in the cloud
func (i *Infra) CreateSubnet(id string, resourceDef Subnet) (*Resource[Subnet], error) {
	return createResource[Subnet](i, id, resourceDef, i.provider.CreateSubnet)
}

// CreateDNS is a dns creator
func (i *Infra) CreateDNS(id string, resourceDef DNS) (*Resource[DNS], error) {
	return createResource[DNS](i, id, resourceDef, i.provider.CreateDNS)
}

// createResource provides shared functionality among the create* functions
func createResource[T any](infra *Infra, id string, resourceDef T, create func(T) (*Resource[T], error)) (*Resource[T], error) {
	if infra.provider == nil {
		return nil, fmt.Errorf("provider is nil")
	}
	if infra.resources == nil {
		return nil, fmt.Errorf("you must create Infra with the New function")
	}
	if _, ok := infra.resources[id]; ok {
		return nil, fmt.Errorf("you cannot call createResource for the same id twice")
	}
	if !infra.provider.Exists(id) {
		Resource, err := create(resourceDef)
		if err != nil {
			return nil, err
		}
		infra.provider.Set(id, Resource)
	}
	handler := infra.provider.Get(id)
	if handler == nil {
		return nil, fmt.Errorf("Provider has returned a nil handler")
	}
	infra.resources[id] = resourceDef
	return handler.(*Resource[T]), nil
}

type mapInterface interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Exists(key string) bool
}
