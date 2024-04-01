package infra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TVPCCreator mocks the creation of VPC resources for testing.
type TResourceManager[T any] struct {
	Resource *Resource[T]
	Err      error
	creates  uint
	updates  uint
	deletes  uint
}

// Create simulates a resource creation, returning a predefined response or error.
func (rm *TResourceManager[T]) Create(id string, r T) (*Resource[T], error) {
	rm.creates++
	return rm.Resource, rm.Err
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

// TestInnerCreate evaluates the behavior of the innerCreate function under various conditions.
// The innerCreate function is shared among all create* functions
func TestInnerCreate(t *testing.T) {

	// Tests the scenario where infra is nil
	{
		_, err := innerCreate(nil, "id", "def", &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where id is empty
	{
		_, err := innerCreate(New(&TestProvider{}), "", "def", &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where manager is nil
	{
		_, err := innerCreate(New(&TestProvider{}), "id", "def", nil)
		assert.NotNil(t, err)
	}

	// Tests the scenario where both the infra provider and resources are nil, expecting an error.
	{
		_, err := innerCreate(&Infra{
			provider:  nil,
			resources: nil,
		}, "id", "def", &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Tests the scenario where the infra resources are nil, expecting an error.
	{
		_, err := innerCreate(&Infra{
			provider:  &TestProvider{},
			resources: nil,
		}, "id", "def", &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		})
		assert.NotNil(t, err)
	}

	// Verifies that resource creation succeeds when both infra resources and provider are correctly provided.
	{
		rm := &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		}
		_, err := innerCreate(New(&TestProvider{}), "id", "def", rm)
		assert.Equal(t, "id", rm.Resource.id)
		assert.Equal(t, "handler", rm.Resource.handler)
		assert.Nil(t, err)
	}

	// Testing the scenario where the resource manager returns a nil handler
	{
		rm := &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: nil,
			},
		}
		_, err := innerCreate(New(&TestProvider{}), "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Testing the scenario where the resource manager returns an id different from input
	{
		rm := &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id2",
				handler: "handler",
			},
		}
		_, err := innerCreate(New(&TestProvider{}), "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Tests the case where a resource with the given ID already exists, expecting an error.
	{
		rm := &TResourceManager[string]{
			Resource: &Resource[string]{
				id:      "id",
				handler: "handler",
			},
		}
		provider := &TestProvider{}
		infra := &Infra{
			provider:  provider,
			resources: map[string]interface{}{"id": "id"},
		}
		_, err := innerCreate(infra, "id", "def", rm)
		assert.NotNil(t, err)
	}

	// Tests the scenario where the provider's resource store returns nil for the given ID, expecting an error.
	{
		provider := &TestProvider{}
		rm := &TResourceManager[string]{
			Resource: nil,
		}
		infra := New(provider)
		_, err := innerCreate(infra, "id", "def", rm)
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
		_, err := innerCreate(infra, "id", "def", rm)
		assert.NotNil(t, err)
		assert.Equal(t, err, expectedErr)
	}
}

func testCreate[T any](t *testing.T, id string, handler interface{}, def T, create func(id string, resourceDef T) (*Resource[T], error)) {
	resource, err := create(id, def)
	assert.Equal(t, id, resource.id, "ID should match the expected value.")
	assert.Equal(t, handler, resource.handler, "handler should match the expected value.")
	assert.Nil(t, err, "No error should be returned during resource creation.")
}

// TestCreate verifies the functionality of creating VPC, Subnet, and DNS resources
// using a mock infrastructure provider. It checks the correctness of the created
// resources' IDs and handlers against expected values.
func TestResourceCreation(t *testing.T) {
	// Setup: Initializes a new infrastructure with a mock provider containing predefined responses.
	provider := &TestProvider{
		vpc:            TResourceManager[VPC]{Resource: &Resource[VPC]{id: "myvpc", handler: "vpchandler"}},
		dns:            TResourceManager[DNS]{Resource: &Resource[DNS]{id: "mydns", handler: "dnshandler"}},
		subnet:         TResourceManager[Subnet]{Resource: &Resource[Subnet]{id: "mysubnet", handler: "subnethandler"}},
		loadBalancer:   TResourceManager[LoadBalancer]{Resource: &Resource[LoadBalancer]{id: "mylb", handler: "mylbhandler"}},
		launchTemplate: TResourceManager[LaunchTemplate]{Resource: &Resource[LaunchTemplate]{id: "mylt", handler: "mylthandler"}},
		autoScale:      TResourceManager[AutoScale]{Resource: &Resource[AutoScale]{id: "myas", handler: "myashandler"}},
	}
	infra := New(provider)
	testCreate(t, provider.vpc.Resource.id, provider.vpc.Resource.handler, VPC{}, infra.CreateVPC)
	testCreate(t, provider.subnet.Resource.id, provider.subnet.Resource.handler, Subnet{}, infra.CreateSubnet)
	testCreate(t, provider.dns.Resource.id, provider.dns.Resource.handler, DNS{}, infra.CreateDNS)
	testCreate(t, provider.loadBalancer.Resource.id, provider.loadBalancer.Resource.handler, LoadBalancer{}, infra.CreateLoadBalancer)
	testCreate(t, provider.launchTemplate.Resource.id, provider.launchTemplate.Resource.handler, LaunchTemplate{}, infra.CreateLaunchTemplate)
	testCreate(t, provider.autoScale.Resource.id, provider.autoScale.Resource.handler, AutoScale{}, infra.CreateAutoScale)
}

func testCreateError[T any](t *testing.T, id string, expectedErr error, def T, create func(id string, resourceDef T) (*Resource[T], error)) {
	resource, err := create(id, def)
	assert.Nil(t, resource, "Resource should be nil")
	assert.Equal(t, err, expectedErr, "Errors should match")
}

// TestCreate verifies the functionality of creating VPC, Subnet, and DNS resources
// using a mock infrastructure provider. It checks the correctness of the created
// resources' IDs and handlers against expected values.
func TestResourceCreationError(t *testing.T) {
	// Setup: Initializes a new infrastructure with a mock provider containing predefined responses.
	provider := &TestProvider{
		vpc:            TResourceManager[VPC]{Err: fmt.Errorf("vpc error")},
		dns:            TResourceManager[DNS]{Err: fmt.Errorf("dns error")},
		subnet:         TResourceManager[Subnet]{Err: fmt.Errorf("subnet error")},
		loadBalancer:   TResourceManager[LoadBalancer]{Err: fmt.Errorf("lb error")},
		launchTemplate: TResourceManager[LaunchTemplate]{Err: fmt.Errorf("lt error")},
		autoScale:      TResourceManager[AutoScale]{Err: fmt.Errorf("as error")},
	}
	infra := New(provider)
	testCreateError(t, "id", provider.vpc.Err, VPC{}, infra.CreateVPC)
	testCreateError(t, "id", provider.subnet.Err, Subnet{}, infra.CreateSubnet)
	testCreateError(t, "id", provider.dns.Err, DNS{}, infra.CreateDNS)
	testCreateError(t, "id", provider.loadBalancer.Err, LoadBalancer{}, infra.CreateLoadBalancer)
	testCreateError(t, "id", provider.launchTemplate.Err, LaunchTemplate{}, infra.CreateLaunchTemplate)
	testCreateError(t, "id", provider.autoScale.Err, AutoScale{}, infra.CreateAutoScale)
}
