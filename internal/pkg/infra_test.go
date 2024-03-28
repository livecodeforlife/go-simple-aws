package infra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TVPCCreator mocks the creation of VPC resources for testing.
type TResourceManager[T any] struct {
	Response *Resource[T]
	Err      error
	creates  uint
	updates  uint
	deletes  uint
}

// Create simulates a resource creation, returning a predefined response or error.
func (rm *TResourceManager[T]) Create(id string, r T) (*Resource[T], error) {
	rm.creates++
	return rm.Response, rm.Err
}

// Delete simulates a resource deletion
func (rm *TResourceManager[T]) Delete(id string) error {
	return rm.Err
}

// TestProvider aggregates mocks for various resource creators and a resource store.
type TestProvider struct {
	vpc            TResourceManager[VPC]
	dns            TResourceManager[DNS]
	subnet         TResourceManager[Subnet]
	loadBalancer   TResourceManager[LoadBalancer]
	launchTemplate TResourceManager[LaunchTemplate]
	autoScale      TResourceManager[AutoScale]
}

func (p *TestProvider) VPC() ResourceManager[VPC] {
	return &p.vpc
}
func (p *TestProvider) DNS() ResourceManager[DNS] {
	return &p.dns
}
func (p *TestProvider) Subnet() ResourceManager[Subnet] {
	return &p.subnet
}
func (p *TestProvider) LoadBalancer() ResourceManager[LoadBalancer] {
	return &p.loadBalancer
}
func (p *TestProvider) LaunchTemplate() ResourceManager[LaunchTemplate] {
	return &p.launchTemplate
}
func (p *TestProvider) AutoScale() ResourceManager[AutoScale] {
	return &p.autoScale
}

// TestCreateResource evaluates the behavior of the createResource function under various conditions.
func TestCreateResource(t *testing.T) {

	// Tests the scenario where infra is nil
	{
		_, err := create(nil, "id", "def", &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where id is empty
	{
		_, err := create(New(&TestProvider{}), "", "def", &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where manager is nil
	{
		_, err := create(New(&TestProvider{}), "id", "def", nil)
		assert.NotNil(t, err)
	}

	// Tests the scenario where both the infra provider and resources are nil, expecting an error.
	{
		_, err := create(&Infra{
			provider:  nil,
			resources: nil,
		}, "id", "def", &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where the infra resources are nil, expecting an error.
	{
		_, err := create(&Infra{
			provider:  &TestProvider{},
			resources: nil,
		}, "id", "def", &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Verifies that resource creation succeeds when both infra resources and provider are correctly provided.
	{
		rm := &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		}
		_, err := create(New(&TestProvider{}), "id", "def", rm)
		assert.Equal(t, "id", rm.Response.id)
		assert.Equal(t, "handler", rm.Response.handler)
		assert.Nil(t, err)
	}

	// Testing the scneario where the resource manager returns a nil handler
	{
		rm := &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: nil,
			},
		}
		_, err := create(New(&TestProvider{}), "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Testing the scenario where the resource manager returns an id different from input
	{
		rm := &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id2",
				handler: "handler",
			},
		}
		_, err := create(New(&TestProvider{}), "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Tests the case where a resource with the given ID already exists, expecting an error.
	{
		rm := &TResourceManager[string]{
			Response: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		}
		provider := &TestProvider{}
		infra := &Infra{
			provider:  provider,
			resources: map[string]interface{}{"id": "id"},
		}
		_, err := create(infra, "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Tests the scenario where the provider's resource store returns nil for the given ID, expecting an error.
	{
		provider := &TestProvider{}
		rm := &TResourceManager[string]{
			Response: nil,
		}
		infra := New(provider)
		_, err := create(infra, "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Tests the case where the provider encounters an error during resource creation.
	{
		expectedErr := fmt.Errorf("Error on creation")
		rm := &TResourceManager[string]{
			Err: expectedErr,
		}
		provider := &TestProvider{}
		infra := New(provider)
		_, err := create(infra, "id", "def", rm)
		assert.NotNil(t, err)
		assert.Equal(t, err, expectedErr)
	}
}

func testCreate[T any](t *testing.T, id string, def T, create func(id string, resourceDef T) (*Resource[T], error)) {
	resource, err := create(id, def)
	assert.Equal(t, id, resource.id, "ID should match the expected value.")
	assert.Equal(t, "handler", resource.handler, "handler should match the expected value.")
	assert.Nil(t, err, "No error should be returned during resource creation.")
}

// TestCreate verifies the functionality of creating VPC, Subnet, and DNS resources
// using a mock infrastructure provider. It checks the correctness of the created
// resources' IDs and handlers against expected values.
func TestResourceCreation(t *testing.T) {
	// Setup: Initializes a new infrastructure with a mock provider containing predefined responses.
	provider := &TestProvider{
		vpc:            TResourceManager[VPC]{Response: &Resource[VPC]{id: "myvpc", handler: "handler"}},
		dns:            TResourceManager[DNS]{Response: &Resource[DNS]{id: "mydns", handler: "handler"}},
		subnet:         TResourceManager[Subnet]{Response: &Resource[Subnet]{id: "mysubnet", handler: "handler"}},
		loadBalancer:   TResourceManager[LoadBalancer]{Response: &Resource[LoadBalancer]{id: "mylb", handler: "handler"}},
		launchTemplate: TResourceManager[LaunchTemplate]{Response: &Resource[LaunchTemplate]{id: "mylt", handler: "handler"}},
		autoScale:      TResourceManager[AutoScale]{Response: &Resource[AutoScale]{id: "myas", handler: "handler"}},
	}
	infra := New(provider)
	testCreate(t, provider.vpc.Response.id, VPC{}, infra.CreateVPC)
	testCreate(t, provider.subnet.Response.id, Subnet{}, infra.CreateSubnet)
	testCreate(t, provider.dns.Response.id, DNS{}, infra.CreateDNS)
	testCreate(t, provider.loadBalancer.Response.id, LoadBalancer{}, infra.CreateLoadBalancer)
	testCreate(t, provider.launchTemplate.Response.id, LaunchTemplate{}, infra.CreateLaunchTemplate)
	testCreate(t, provider.autoScale.Response.id, AutoScale{}, infra.CreateAutoScale)
}
