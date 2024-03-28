package infra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TVPCCreator struct {
	Response *Resource[VPC]
	Err      error
}

func (c *TVPCCreator) CreateVPC(r VPC) (*Resource[VPC], error) {
	return c.Response, c.Err
}

type TDNSCreator struct {
	Response *Resource[DNS]
	Err      error
}

func (c *TSubnetCreator) CreateSubnet(r Subnet) (*Resource[Subnet], error) {
	return c.Response, c.Err
}

type TSubnetCreator struct {
	Response *Resource[Subnet]
	Err      error
}

func (c *TDNSCreator) CreateDNS(r DNS) (*Resource[DNS], error) {
	return c.Response, c.Err
}

type TResourceStore struct {
	store map[string]interface{}
}

func (p *TResourceStore) Set(key string, handler interface{}) {
	p.store[key] = handler
}
func (p *TResourceStore) Get(key string) interface{} {
	return p.store[key]
}
func (p *TResourceStore) Exists(key string) bool {
	_, ok := p.store[key]
	return ok
}

type TestProvider struct {
	TVPCCreator
	TDNSCreator
	TSubnetCreator
	TResourceStore
}

func TestCreateResource(t *testing.T) {
	{
		infra := &Infra{
			provider:  nil,
			resources: nil,
		}
		_, err := createResource(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}
	{
		infra := &Infra{
			provider: &TestProvider{
				TResourceStore: TResourceStore{
					store: make(map[string]interface{}),
				},
			},
			resources: nil,
		}
		_, err := createResource(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}
	{
		//Given a right initialized infra
		provider := &TestProvider{
			TResourceStore: TResourceStore{
				store: make(map[string]interface{}),
			},
		}
		infra := New(provider)
		resource, err := createResource(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.Equal(t, resource.id, "id")
		assert.Equal(t, resource.handler, "handler")
		assert.Nil(t, err)
	}
	{
		//Test with same id
		infra := &Infra{
			provider: &TestProvider{
				TResourceStore: TResourceStore{
					store: make(map[string]interface{}),
				},
			},
			resources: map[string]interface{}{
				"id": "id",
			},
		}
		_, err := createResource(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}
	{
		//When provider returns nil
		provider := &TestProvider{
			TResourceStore: TResourceStore{
				store: map[string]interface{}{
					"id": nil,
				},
			},
		}
		infra := New(provider)
		_, err := createResource(infra, "id", "def", func(string) (*Resource[string], error) {
			return &Resource[string]{id: "id", handler: "handler"}, nil
		})
		assert.NotNil(t, err)
	}

	{
		//When provider returns error on creation
		provider := &TestProvider{
			TResourceStore: TResourceStore{
				store: make(map[string]interface{}),
			},
		}
		infra := New(provider)
		_, err := createResource(infra, "id", "def", func(string) (*Resource[string], error) {
			return nil, fmt.Errorf("Error on creation")
		})
		assert.Equal(t, fmt.Sprintf("%s", err), "Error on creation")
	}

}

func TestCreate(t *testing.T) {
	infra := New(&TestProvider{
		TResourceStore: TResourceStore{
			store: make(map[string]interface{}),
		},
		TVPCCreator: TVPCCreator{
			Response: &Resource[VPC]{id: "myvpc", handler: "handler"},
		},
		TSubnetCreator: TSubnetCreator{
			Response: &Resource[Subnet]{id: "mysubnet", handler: "handler"},
		},
		TDNSCreator: TDNSCreator{
			Response: &Resource[DNS]{id: "mydns", handler: "handler"},
		},
	})
	vpc, err := infra.CreateVPC("myvpc", VPC{})
	assert.Equal(t, vpc.id, infra.provider.(*TestProvider).TVPCCreator.Response.id)
	assert.Equal(t, vpc.handler, infra.provider.(*TestProvider).TVPCCreator.Response.handler)
	assert.Equal(t, err, infra.provider.(*TestProvider).TVPCCreator.Err)

	subnet, err := infra.CreateSubnet("mysubnet", Subnet{})
	assert.Equal(t, subnet.id, infra.provider.(*TestProvider).TSubnetCreator.Response.id)
	assert.Equal(t, subnet.handler, infra.provider.(*TestProvider).TSubnetCreator.Response.handler)
	assert.Equal(t, err, infra.provider.(*TestProvider).TSubnetCreator.Err)

	dns, err := infra.CreateDNS("mydns", DNS{})
	assert.Equal(t, dns.id, infra.provider.(*TestProvider).TDNSCreator.Response.id)
	assert.Equal(t, dns.handler, infra.provider.(*TestProvider).TDNSCreator.Response.handler)
	assert.Equal(t, err, infra.provider.(*TestProvider).TDNSCreator.Err)

}
