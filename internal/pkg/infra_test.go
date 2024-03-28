package infra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TVPCCreator mocks the creation of VPC resources for testing.
type TVPCCreator struct {
	Response *Resource[VPC]
	Err      error
	creates  uint
	updates  uint
	deletes  uint
}

// CreateVPC simulates VPC creation, returning a predefined response or error.
func (c *TVPCCreator) CreateVPC(r VPC) (*Resource[VPC], error) {
	return c.Response, c.Err
}

// CreateVPC simulates VPC creation, returning a predefined response or error.
func (c *TVPCCreator) Create(r VPC) (*Resource[VPC], error) {
	return c.Response, c.Err
}

// CreateVPC simulates VPC creation, returning a predefined response or error.
func (c *TVPCCreator) Update(r VPC, old *Resource[VPC]) (*Resource[VPC], error) {
	return c.Response, c.Err
}

// CreateVPC simulates VPC creation, returning a predefined response or error.
func (c *TVPCCreator) Delete(old *Resource[VPC]) error {
	return c.Err
}

// TDNSCreator mocks the creation of DNS resources for testing.
type TDNSCreator struct {
	Response *Resource[DNS]
	Err      error
}

// TSubnetCreator mocks the creation of Subnet resources for testing.
type TSubnetCreator struct {
	Response *Resource[Subnet]
	Err      error
}

// CreateSubnet simulates Subnet creation, returning a predefined response or error.
func (c *TSubnetCreator) CreateSubnet(r Subnet) (*Resource[Subnet], error) {
	return c.Response, c.Err
}

// CreateDNS simulates DNS creation, returning a predefined response or error.
func (c *TDNSCreator) CreateDNS(r DNS) (*Resource[DNS], error) {
	return c.Response, c.Err
}

// TResourceStore provides a mock store for resource handlers, supporting basic operations.
type TResourceStore struct {
	store map[string]interface{}
}

// Set stores a resource handler associated with a key.
func (p *TResourceStore) Set(key string, handler interface{}) {
	p.store[key] = handler
}

// Get retrieves a resource handler by key.
func (p *TResourceStore) Get(key string) interface{} {
	return p.store[key]
}

// Exists checks if a resource handler exists for a given key.
func (p *TResourceStore) Exists(key string) bool {
	_, ok := p.store[key]
	return ok
}

// TestProvider aggregates mocks for various resource creators and a resource store.
type TestProvider struct {
	TVPCCreator
	TDNSCreator
	TSubnetCreator
	TResourceStore
}

func (p *TestProvider) VPC() ResourceManager[VPC] {
	return &p.TVPCCreator
}

// TestCreateResource evaluates the behavior of the createResource function under various conditions.
func TestCreateResource(t *testing.T) {
	// Tests the scenario where both the infra provider and resources are nil, expecting an error.
	{
		infra := &Infra{provider: nil, resources: nil}
		_, err := apply(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where the infra resources are nil, expecting an error.
	{
		infra := &Infra{
			provider:  &TestProvider{TResourceStore: TResourceStore{store: make(map[string]interface{})}},
			resources: nil,
		}
		_, err := apply(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}

	// Verifies that resource creation succeeds when both infra resources and provider are correctly provided.
	{
		provider := &TestProvider{TResourceStore: TResourceStore{store: make(map[string]interface{})}}
		infra := New(provider)
		resource, err := apply(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.Equal(t, "id", resource.id)
		assert.Equal(t, "handler", resource.handler)
		assert.Nil(t, err)
	}

	// Tests the case where a resource with the given ID already exists, expecting an error.
	{
		infra := &Infra{
			provider:  &TestProvider{TResourceStore: TResourceStore{store: make(map[string]interface{})}},
			resources: map[string]interface{}{"id": "id"},
		}
		_, err := apply(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where the provider's resource store returns nil for the given ID, expecting an error.
	{
		provider := &TestProvider{TResourceStore: TResourceStore{store: map[string]interface{}{"id": nil}}}
		infra := New(provider)
		_, err := apply(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}

	// Tests the case where the provider encounters an error during resource creation.
	{
		provider := &TestProvider{TResourceStore: TResourceStore{store: make(map[string]interface{})}}
		infra := New(provider)
		_, err := apply(infra, "id", "def", func(string) (*Resource[string], error) {
			return nil, fmt.Errorf("Error on creation")
		})
		assert.Equal(t, "Error on creation", fmt.Sprintf("%s", err))
	}
}

// TestCreate verifies the functionality of creating VPC, Subnet, and DNS resources
// using a mock infrastructure provider. It checks the correctness of the created
// resources' IDs and handlers against expected values.
func TestCreate(t *testing.T) {
	// Setup: Initializes a new infrastructure with a mock provider containing predefined responses.
	infra := New(&TestProvider{
		TResourceStore: TResourceStore{store: make(map[string]interface{})},
		TVPCCreator:    TVPCCreator{Response: &Resource[VPC]{id: "myvpc", handler: "handler"}},
		TSubnetCreator: TSubnetCreator{Response: &Resource[Subnet]{id: "mysubnet", handler: "handler"}},
		TDNSCreator:    TDNSCreator{Response: &Resource[DNS]{id: "mydns", handler: "handler"}},
	})

	// Test VPC Creation: Ensures the VPC is created as expected with correct ID and handler.
	vpc, err := infra.CreateVPC("myvpc", VPC{})
	assert.Equal(t, "myvpc", vpc.id, "VPC ID should match the expected value.")
	assert.Equal(t, "handler", vpc.handler, "VPC handler should match the expected value.")
	assert.Nil(t, err, "No error should be returned during VPC creation.")

	// Test Subnet Creation: Validates the Subnet creation with its ID and handler.
	subnet, err := infra.CreateSubnet("mysubnet", Subnet{})
	assert.Equal(t, "mysubnet", subnet.id, "Subnet ID should match the expected value.")
	assert.Equal(t, "handler", subnet.handler, "Subnet handler should match the expected value.")
	assert.Nil(t, err, "No error should be returned during Subnet creation.")

	// Test DNS Creation: Confirms the DNS resource creation is accurate with correct ID and handler.
	dns, err := infra.CreateDNS("mydns", DNS{})
	assert.Equal(t, "mydns", dns.id, "DNS ID should match the expected value.")
	assert.Equal(t, "handler", dns.handler, "DNS handler should match the expected value.")
	assert.Nil(t, err, "No error should be returned during DNS creation.")
}
